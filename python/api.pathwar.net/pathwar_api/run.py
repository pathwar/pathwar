from uuid import UUID, uuid4
import bcrypt
import json
import random
import sys

from eve import Eve
from eve.auth import BasicAuth, TokenAuth
from eve.io.base import BaseJSONEncoder
from eve.io.mongo import Validator
from eve.methods.post import post, post_internal
from eve_docs import eve_docs
from flask import abort, url_for
from flask.ext.bootstrap import Bootstrap

from settings import DOMAIN
from seeds import db_reset, db_init, db_seed
from tools import bp_tools
from mail import mail, send_mail


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
            # item['otp_secret'] = ...
            on_update_user(item)

            if not app.is_seed:
                # FIXME: put after insert success
                verification_url = url_for(
                    'tools.email_verify',
                    user_id=item['_id'],
                    email_verification_token=item['email_verification_token'],
                    _external=True,
                )
                message = 'Verification link: {}'.format(verification_url)
                send_mail(
                    message=message,
                    subject='Email verification',
                    recipients=[item]
                )

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
        # print(response.get_data())
        app.logger.warn(dir(response))

    elif resource == 'organizations':
        organization = json.loads(response.get_data())
        user = request_get_user(request)
        post_internal('organization-users', {
            'organization': organization['_id'],
            'role': 'owner',
            'user': user['_id'],
        })
        orga_statistics = post_internal('organization-statistics', {
            'organization': organization['_id'],
        })
        app.data.driver.db['organizations'].update(
            { '_id': organization['_id'] },
            { 'statistics': orga_statistics[0]['_id'] },
        )


# Initialize Eve
app = Eve(
    auth=MockBasicAuth,
    # auth=MockTokenAuth,
    json_encoder=UUIDEncoder,
    validator=UUIDValidator,
)


def eve_init():
    # eve-docs
    Bootstrap(app)
    app.register_blueprint(eve_docs, url_prefix='/docs')

    # tools
    app.register_blueprint(bp_tools, url_prefix='/tools')

    # mail
    mail.init_app(app)

    # Attach hooks
    app.on_pre_GET += pre_get_callback
    app.on_insert += insert_callback
    app.on_pre_POST += pre_post_callback
    app.on_post_POST += post_post_callback
    # getattr(app, 'on_pre_POST_user-tokens') += pre_post_user_tokens_callback

    # Initialize db
    db_init(app)


def main(argv):
    eve_init()

    if len(argv) > 1:
        if argv[1] == 'flush-db':
            with app.app_context():
                db_reset(app)
        elif argv[1] == 'seed-db':
            with app.app_context():
                db_reset(app)
                app.is_seed = True
                db_seed(app)
                app.is_seed = False

    else:
        # Run
        app.run(
            debug=True,
            host='0.0.0.0',
        )


if __name__ == '__main__':
    main(sys.argv)
