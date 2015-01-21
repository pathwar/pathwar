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
from flask import abort
from flask.ext.bootstrap import Bootstrap

from settings import DOMAIN
from seeds import load_seeds
from tools import bp_tools
from mail import send_mail


def request_get_user(request):
    auth = request.authorization
    if auth.get('username'):
        if auth.get('password'):  # user:pass
            app.logger.warn('FIXME: check password')
            return app.data.driver.db['users'].find_one({
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
    """ Callback called before a GET request, we can alter the lookup. """
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


def on_update_user(item):
    """ Must be called when saving a user POST/PATCH/PUT on /users. """
    if 'password' in item and \
       len(item['password']) and \
       not item['password'].startswith('$2a$'):
        # FIXME: better check for bcrypt format
        password = item['password'].encode('utf-8')
        item['password'] = bcrypt.hashpw(
            password, item['password_salt']
        )


def insert_callback(resource, items):
    """ Callback called just before inserting a resource in mongo. """
    app.logger.info('### insert_callback({}) {}'.format(resource, items))
    for item in items:
        item['_id'] = str(uuid4())

    if resource == 'users':
        for item in items:
            item['password_salt'] = bcrypt.gensalt().encode('utf-8')
            item['email_verification_token'] = str(uuid4())
            #item['otp_secret'] = ...
            on_update_user(item)

            # FIXME: put after insert success
            continue
            # FIXME: do not send_mail for seeds
            send_mail(item, url_for(
                '/tools/email-verify/{}/{}'.format(
                    item['_id'],
                    'Verification link: {}'.format(
                        item['email_verification_token'],
                    ),
                )
            ))

    app.logger.info('### insert_callback({}) {}'.format(resource, items))


def pre_post_callback(resource, request):
    """ Callback called just before the normal processing behavior of a POST
    request.
    """
    if resource == 'user-tokens':
        # Handle login
        user = request_get_user(request)
        app.logger.warn('@@@ pre_post_callback: user={}'.format(user))
        if not user:
            abort(401)
        # FIXME: try to not accept passing token/user (read-only)
        payload = request.get_json()
        payload['token'] = str(uuid4())
        payload['user'] = user['_id']
        # FIXME: add expiry_date
    elif resource == 'users':
        # FIXME: check for a password, users without password are built
        #        internally
        pass


def post_post_callback(resource, request, response):
    """ Callback called just after a POST request ended. """
    app.logger.info('### post_post({}) request: {}, response: {}'
                .format(resource, request, response))
    if resource == 'users':
        print(response.get_data())
        app.logger.warn(dir(response))


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
    app.register_blueprint(bp_tools, url_prefix='/tools')

    # Attach hooks
    app.on_pre_GET += pre_get_callback
    app.on_insert += insert_callback
    app.on_pre_POST += pre_post_callback
    app.on_post_POST += post_post_callback
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
