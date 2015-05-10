# -*- coding: utf-8 -*-

import datetime

from flask import current_app

from models import (
    Achievement,
    Coupon,
    GlobalStatistics,
    Level,
    LevelInstanceUser,
    LevelStatistics,
    Organization,
    OrganizationAchievement,
    OrganizationCoupon,
    OrganizationLevel,
    OrganizationStatistics,
    Session,
    Task,
    User,
    UserToken,
)


class Job(object):
    def __init__(self, task_id):
        self.task_id = task_id

    def update(self, data):
        return Task.update_by_id(self.task_id, data)

    def update_status(self, status):
        self.update({
            '$set': {
                'status': status,
            },
        })

    def call(self):
        self.update_status('in progress')

        # FIXME: decorate with try/except
        self.run()

        self.update_status('succeed')


class TestJob(Job):
    name = 'test'

    def run(self):
        current_app.logger.warn('test job called')


class CleanJob(Job):
    name = 'clean'

    def run(self):
        now = datetime.datetime.utcnow()

        # Flush tasks collection
        # current_app.data.driver.db['raw-tasks'].drop()

        # Clean expired api tokens
        UserToken.remove({
            'expiry_date': {
                '$lt': now,
            },
        })

        # Clean expired level access tokens
        LevelInstanceUser.remove({
            'expiry_date': {
                '$lt': now,
            },
        })


class UpdateStatsJob(Job):
    name = 'update-stats'

    def run(self):
        # global statistics
        gs = {
            'level_bought': 0,
            'level_finished': 0,
        }
        gs['users'] = User.count({'active': True})
        # FIXME: remove beta levels
        gs['achievements'] = Achievement.count()
        gs['expired_coupons'] = Coupon.count({
            'validations_left': 0,
        })
        gs['coupons'] = Coupon.count()
        gs['level_bought'] = OrganizationLevel.count()
        gs['level_finished'] = OrganizationLevel.count({
            'status': {
                '$in': ['pending validation', 'validated'],
            },
        })
        gs['levels'] = Level.count()
        gs['organization_achievements'] = OrganizationAchievement.count()
        gs['organization_coupons'] = OrganizationCoupon.count()
        gs['organizations'] = Organization.count()

        # level statistics
        levels = Level.all()
        for level in levels:
            amount_bought = OrganizationLevel.count({
                'level': level['_id'],
            })
            # FIXME: filter, only 1 by organization
            amount_finished = OrganizationLevel.count({
                'level': level['_id'],
                'status': {
                    '$in': ['pending validation', 'validated'],
                },
            })

            LevelStatistics.update_by_id(level['statistics'], {
                '$set': {
                    'amount_bought': amount_bought,
                    'amount_finished': amount_finished,
                    # FIXME: fivestar average
                    # FIXME: duration average
                    # FIXME: amount hints bought
                },
            })

        sessions = Session.all()
        for session in sessions:
            organization_levels = OrganizationLevel.find({
                'session': session['_id'],
            })
            organizations = Organization.find({
                'session': session['_id'],
            })
            for organization in organizations:
                coupons = OrganizationCoupon.count({
                    'organization': organization['_id'],
                })
                achievements = OrganizationAchievement.count({
                    'organization': organization['_id'],
                })
                # current_app.logger.debug(
                #     'session=%s organization=%s coupons=%d',
                #     session['name'],
                #     organization['name'],
                #     coupons,
                # )
                OrganizationStatistics.update_by_id(organization['statistics'], {
                    '$set': {
                        'coupons': coupons,
                        'achievements': achievements,
                    },
                })

        last_record = GlobalStatistics.last_record()
        if not all(item in last_record.items()
                   for item in gs.items()):
            GlobalStatistics.post_internal(gs)


JOBS_CLASSES = (
    CleanJob,
    TestJob,
    UpdateStatsJob,
)


def setup_jobs(app):
    jobs = {}

    for job in JOBS_CLASSES:
        jobs[job.name] = job

    setattr(app, 'jobs', jobs)
