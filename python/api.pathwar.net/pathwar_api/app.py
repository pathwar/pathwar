from uuid import UUID
import json

from eve import Eve
from eve.auth import BasicAuth, TokenAuth
from eve.io.base import BaseJSONEncoder
from eve.io.mongo import Validator
from flask import abort, url_for


class MockBasicAuth(BasicAuth):
    def check_auth(self, username, password, allowed_roles, resource,
                   method):
        if not len(username):
            return False

        user = None

        if len(password):  # Login+Password based
            app.logger.warn('FIXME: check password')
            user = app.data.driver.db['users'].find_one({
                'login': username, 'active': True
            })
        else:  # Token-based
            user_token = app.data.driver.db['user-tokens'] \
                                        .find_one({'token': username})
            if user_token:
                user = app.data.driver.db['users'] \
                                      .find_one({'_id': user_token['user']})
                if not user['active']:
                    user = None


        if user:
            app.logger.warn('username: {}, password: {}, allowed_roles: {}, '
                            'resource: {}, method: {}, user: {}'.format(
                                username, password, allowed_roles, resource,
                                method, user,
                            )
            )
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
