from uuid import UUID, uuid4
import bcrypt
import json
import random

from eve import Eve
from eve.auth import BasicAuth, TokenAuth
from eve.io.base import BaseJSONEncoder
from eve.io.mongo import Validator
from eve.methods.post import post, post_internal
from eve_docs import eve_docs
from flask.ext.bootstrap import Bootstrap

from settings import DOMAIN
from seeds import load_seeds


def request_get_user(request):
    auth = request.authorization
    if auth.get('username'):
        if auth.get('password'):  # user:pass
            app.logger.warn('FIXME: check password')
            return app.data.driver.db['users'].find({
                'login': auth.get('username'),
                'active': True
            })
        else:  # token
            user_token = app.data.driver.db['user-tokens'] \
                                        .find_one({
                                            'token': auth.get('username')
                                        })
            if user_token:
                return app.data.driver.db['users'] \
                                      .find_one({'_id': user_token['user']})
    return None


class MockBasicAuth(BasicAuth):
    def check_auth(self, username, password, allowed_roles, resource,
                   method):
        if not len(username):
            return False
        if len(password):  # Login+Password based
            app.logger.warn('FIXME: check password')
            return app.data.driver.db['users'].find({
                'login': username, 'active': True
            })
        else:  # Token-based
            user_token = app.data.driver.db['user-tokens'] \
                                        .find_one({'token': username})
            if user_token:
                user = app.data.driver.db['users'] \
                                      .find_one({'_id': user_token['user']})
                if user['active']:
                    return True
            return False


class MockTokenAuth(TokenAuth):
    def check_auth(self, token, allowed_roles, resource, method):
        user_tokens = app.data.driver.db['user-tokens']
        return user_tokens.find_one({'token': token})


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


def pre_get_callback(resource, request, lookup):
    resources_with_me_filter = (
        'user-notifications', 'user-organization-invites', 'user-tokens',
        'organization-users',
    )
    # Handle users/me
    if resource == 'users' and 'login' in lookup:
        del lookup['login']
        lookup['_id'] = request_get_user(request)['_id']
    elif resource in resources_with_me_filter:
        app.logger.warn('FIXME: handle where user==me')


def insert_callback(resource, items):
    for item in items:
        item['_id'] = str(uuid4())

    if resource == 'users':
        for item in items:
            item['password_salt'] = bcrypt.gensalt().encode('utf-8')
            #item['otp_secret'] = ...

            # POST/PATCH/PUT
            if item['password'] and \
               len(item['password']) and \
               not item['password'].startswith('$2a$'):
                # FIXME: better check for bcrypt format
                password = item['password'].encode('utf-8')
                item['password'] = bcrypt.hashpw(
                    password, item['password_salt']
                )

    app.logger.info('### insert_callback({}) {}'.format(resource, items))


def pre_post_callback(resource, request):
    if resource == 'user-tokens':
        # Handle login
        user = request_get_user(request)
        app.logger.warn('@@@ {}'.format(user))
        payload = request.get_json()
        payload['token'] = str(uuid4())
        payload['user'] = user['_id']
        app.logger.warn('### {}'.format(payload))
        return post_internal(resource, payload, skip_validation=True)


# Initialize Eve
app = Eve(
    auth=MockBasicAuth,
    # auth=MockTokenAuth,
    json_encoder=UUIDEncoder,
    validator=UUIDValidator,
)


def main():
    # eve-docs
    Bootstrap(app)
    app.register_blueprint(eve_docs, url_prefix='/docs')

    # Attach hooks
    app.on_pre_GET += pre_get_callback
    app.on_insert += insert_callback
    app.on_pre_POST += pre_post_callback
    #getattr(app, 'on_pre_POST_user-tokens') += pre_post_user_tokens_callback

    # Initialize data
    with app.app_context():
        load_seeds(app, reset=True)

    # Run
    app.run(
        debug=True,
        host='0.0.0.0',
    )


if __name__ == '__main__':
    main()
