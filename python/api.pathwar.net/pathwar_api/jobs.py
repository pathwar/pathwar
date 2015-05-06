# -*- coding: utf-8 -*-

from flask import current_app

from models import (
    Level,
    LevelStatistics,
    Organization,
    OrganizationLevel,
    OrganizationStatistics,
    Session,
    Task,
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
    TestJob,
    UpdateStatsJob,
)


def setup_jobs(app):
    jobs = {}

    for job in JOBS_CLASSES:
        jobs[job.name] = job

    setattr(app, 'jobs', jobs)
