from uuid import uuid4
import json
import bcrypt
import random
import md5

from eve.methods.post import post, post_internal

from app import app


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
    item['gravatar_hash'] = md5.new(item['email'].lower().strip()).hexdigest()


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

            if not app.is_seed and not item['active']:
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
    dct = json.loads(response.get_data())
    if '_items' in dct:
        items = dct['_items']
    else:
        items = [dct]

    if resource == 'users':
        # FIXME: create a user notification
        app.logger.warn('%' * 800)

        worldwide_session = app.data.driver.db['sessions'].find_one({
            'name': 'Worldwide'
        })
        default_organization = post_internal('organizations', {
            'session': worldwide_session['_id'],
        })
        app.logger.warn(dir(response))

    elif resource == 'organizations':
        for organization in items:
            user = request_get_user(request)
            if not app.is_seed:
                app.logger.error(organization)
                post_internal('organization-users', {
                    'organization': organization['_id'],
                    'role': 'owner',
                    'user': user['_id'],
                })
                orga_statistics = post_internal('organization-statistics', {
                    'organization': organization['_id'],
                })

            # app.data.driver.db['organizations'].update(
            #     { '_id': organization['_id'] },
            #     { 'statistics': orga_statistics[0]['_id'] },
            # )


def setup_hooks(_app):
    # Attach hooks
    app.on_pre_GET += pre_get_callback
    app.on_insert += insert_callback
    app.on_pre_POST += pre_post_callback
    app.on_post_POST += post_post_callback
    # getattr(app, 'on_pre_POST_user-tokens') += pre_post_user_tokens_callback
