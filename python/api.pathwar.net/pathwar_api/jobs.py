# -*- coding: utf-8 -*-

import datetime

from flask import current_app

from models import (
    Level,
    LevelInstanceUser,
    LevelStatistics,
    Organization,
    OrganizationLevel,
    OrganizationStatistics,
    Session,
    Task,
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
        sessions = Session.all()
        # levels = Level.all()
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
