# -*- coding: utf-8 -*-

import datetime

from flask import current_app

from models import (
    GlobalStatistics,
    Level,
    LevelInstanceUser,
    LevelStatistics,
    Organization,
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
        global_statistics = {
            'level_bought': 0,
            'level_finished': 0,
        }
        global_statistics['users'] = User.count({'active': True})
        # FIXME: remove beta levels
        global_statistics['levels'] = Level.count()
        global_statistics['organizations'] = Organization.count()

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
                    '$in': ['pending validation'],
                },
            })
            global_statistics['level_bought'] += amount_bought
            global_statistics['level_finished'] += amount_finished

            LevelStatistics.update_by_id(level['statistics'], {
                '$set': {
                    'amount_bought': amount_bought,
                    'amount_finished': amount_finished,
                    # FIXME: fivestar average
                    # FIXME: duration average
                    # FIXME: amount hints bought
                },
            })

        last_record = GlobalStatistics.last_record()
        if not all(item in last_record.items()
                   for item in global_statistics.items()):
            GlobalStatistics.post_internal(global_statistics)

        return

        sessions = Session.all()
        for session in sessions:
            organization_levels = OrganizationLevel.find({
                'session': session['_id'],
            })
            organizations = Organization.find({
                'session': session['_id'],
            })
            for organization in organizations:
                current_app.logger.info(
                    'session=%s organization=%s',
                    session['name'],
                    organization['name'],
                )


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
