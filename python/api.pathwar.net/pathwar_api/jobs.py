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
    OrganizationLevelValidation,
    OrganizationStatistics,
    Session,
    Task,
    User,
    UserToken,
)
from utils import mongo_list_to_dict


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

            organizations = mongo_list_to_dict(
                Organization.find({
                    'session': session['_id'],
                })
            )

            levels = mongo_list_to_dict(
                Level.find({})
            )
            
            validations = OrganizationLevelValidation.find()

            for validation in validations:
                if validation['status'] == 'refused':
                    continue
                organization = organizations.get(validation['organization'])
                if not organization:
                    continue
                level = levels.get(validation['level'])
                if not level:
                    current_app.logger.warn(level)
                    continue
                level['validations'] = level.get('validations', 0)

                defaults = {
                    'validated_levels': [],
                    'gold_medals': 0,
                    'silver_medals': 0,
                    'bronze_medals': 0,
                }
                for key, value in defaults.items():
                    if key not in organization:
                        organization[key] = value

                if level['validations'] == 0:
                    organization['gold_medals'] += 1
                elif level['validations'] == 1:
                    organization['silver_medals'] += 1
                elif level['validations'] == 2:
                    organization['bronze_medals'] += 1
                organization['validated_levels'].append(level['_id'])
                level['validations'] += 1


            for organization in organizations.values():
                coupons = OrganizationCoupon.count({
                    'organization': organization['_id'],
                })
                achievements = OrganizationAchievement.count({
                    'organization': organization['_id'],
                })
                validated_levels = list(set(
                    organization.get('validated_levels', [])
                ))
                score = 0
                score += organization.get('gold_medals', 0) * 5
                score += organization.get('silver_medals', 0) * 3
                score += organization.get('bronze_medals', 0) * 1
                score += len(validated_levels) * 10
                score += achievements * 2
                # current_app.logger.debug(
                #     'session=%s organization=%s coupons=%d',
                #     session['name'],
                #     organization['name'],
                #     coupons,
                # )
                OrganizationStatistics.update_by_id(organization['statistics'], {
                    '$set': {
                        'coupons': coupons,
                        'score': score,
                        'achievements': achievements,
                        'gold_medals': organization.get('gold_medals', 0),
                        'silver_medals': organization.get('silver_medals', 0),
                        'bronze_medals': organization.get('bronze_medals', 0),
                        'finished_levels': len(validated_levels),
                        'bought_levels': OrganizationLevel.count({
                            'organization': organization['_id'],
                        })
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
