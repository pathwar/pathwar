from uuid import UUID
import json
import os
import bcrypt
import logging.config
import logging

from eve import Eve
from eve.auth import BasicAuth, TokenAuth
from eve.io.base import BaseJSONEncoder
from eve.io.mongo import Validator
from flask import abort, url_for
from raven.handlers.logging import SentryHandler
from raven.conf import setup_logging


SENTRY_URL = os.environ.get('SENTRY_URL', '')


class MockBasicAuth(BasicAuth):
    def check_auth(self, username, password, allowed_roles, resource,
                   method):
        if not len(username):
            return False

        user = None

        if len(password):
            # FIXME: restrict login+password access for a minimal amount of
            #        resources
            user = app.data.driver.db['raw-users'].find_one({
                'login': username,
                # 'active': True,  # FIXME: re-enable
            })
            if user and 'password_salt' in user:
                hash_ = bcrypt.hashpw(password, str(user['password_salt']))
                if hash_ != user['password']:
                    user = None
            else:
                user = None
        else:  # Token-based
            user_token = app.data.driver.db['raw-user-tokens'] \
                                        .find_one({'token': username})
            if user_token:
                user = app.data.driver.db['raw-users'] \
                                      .find_one({'_id': user_token['user']})
                if user['role'] != 'admin':
                    self.set_request_auth_value(user['_id'])
                # FIXME: Re-enable later
                if not user['active']:
                    user = None
                    abort(401, 'You need to validate your email address first')

        if user:
            # app.logger.debug(
            #     'username: {}, password: {}, allowed_roles: {}, '
            #     'resource: {}, method: {}, user: {}'.format(
            #         username, password, allowed_roles, resource,
            #         method, user,
            #     )
            # )
            if user['role'] in allowed_roles:
                return True
            else:
                return False
        else:
            app.logger.warn('No such active user: {}'.format(username))
            return False


class UUIDValidator(Validator):
    """
    Extends the base mongo validator adding support for the uuid data-type
    """
    def _validate_type_uuid(self, field, value):
        try:
            UUID(value)
        except ValueError:
            self._error(field, "value '%s' cannot be converted to a UUID" %
                        value)


class UUIDEncoder(BaseJSONEncoder):
    """ JSONEconder subclass used by the json render function.
    This is different from BaseJSONEoncoder since it also addresses
    encoding of UUID
    """

    def default(self, obj):
        if isinstance(obj, UUID):
            return str(obj)
        else:
            # delegate rendering to base class method (the base class
            # will properly render ObjectIds, datetimes, etc.)
            return super(UUIDEncoder, self).default(obj)


# Initialize Eve
app = Eve(
    auth=MockBasicAuth,
    json_encoder=UUIDEncoder,
    validator=UUIDValidator,
)

if len(SENTRY_URL):
    LOGGING = {
        'version': 1,
        'disable_existing_loggers': False,
        'formatters': {
            'console': {
                'format': (
                    '[%(asctime)s][%(levelname)s] %(name)s ' +
                    '%(filename)s:%(funcName)s:%(lineno)d | %(message)s'
                ),
                'datefmt': '%H:%M:%S',
            },
        },

        'handlers': {
            'console': {
                'level': 'DEBUG',
                'class': 'logging.StreamHandler',
                'formatter': 'console'
            },
            'sentry': {
                'level': 'INFO',
                'class': 'raven.handlers.logging.SentryHandler',
                'dsn': SENTRY_URL,
                'site': 'api',
            },
        },

        'loggers': {
            '': {
                'handlers': ['console', 'sentry'],
                'level': 'DEBUG',
                'propagate': True,
            },
            'api': {
                'handlers': ['console', 'sentry'],
                'level': 'DEBUG',
                'propagate': True,
            },
        }
    }
    logging.config.dictConfig(LOGGING)
