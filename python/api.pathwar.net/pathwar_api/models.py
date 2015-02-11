from uuid import uuid4
import json
import random
import md5

import bcrypt
from eve.methods.post import post, post_internal

from app import app
from utils import request_get_user


class BaseItem(object):
    def __init__(self):
        pass

    def on_update(self, item):
        pass

    def on_insert(self, item):
        item['_id'] = str(uuid4())

    def on_pre_post(self, request):
        pass

    def on_pre_get(self, request, lookup):
        pass

    def on_post_post_item(self, request, response, item):
        pass

    def on_post_post(self, request, response):
        dct = json.loads(response.get_data())
        if '_items' in dct:
            items = dct['_items']
        else:
            items = [dct]

        for item in items:
            self.on_post_post_item(request, response, item)


class UserItem(BaseItem):
    resource = 'users'

    def on_update(self, item):
        if 'password' in item and \
           len(item['password']) and \
           not item['password'].startswith('$2a$'):
            # FIXME: better check for bcrypt format
            password = item['password'].encode('utf-8')
            item['password'] = bcrypt.hashpw(
                password, item['password_salt']
            )
        item['gravatar_hash'] = md5.new(
            item['email'].lower().strip()
        ).hexdigest()

    def on_insert(self, item):
        super(UserItem, self).on_insert(item)
        item['password_salt'] = bcrypt.gensalt().encode('utf-8')
        item['email_verification_token'] = str(uuid4())
        # item['otp_secret'] = ...
        self.on_update(item)

        # Send verification email
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

    def on_pre_post(self, request):
        # FIXME: check for a password, users without password are built
        #        internally
        pass

    def on_pre_get(self, request, lookup):
        # Handle users/me
        if 'login' in lookup:
            del lookup['login']
            lookup['_id'] = request_get_user(request)['_id']

    def on_post_post_item(self, request, response, item):
        # FIXME: create a user notification
        app.logger.warn('%' * 800)

        worldwide_session = app.data.driver.db['sessions'].find_one({
            'name': 'Worldwide'
        })
        default_organization = post_internal('organizations', {
            'session': worldwide_session['_id'],
        })
        app.logger.warn(dir(response))


class UserTokenItem(BaseItem):
    resource = 'user-tokens'

    def on_pre_post(self, request):
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


class OrganizationItem(BaseItem):
    def on_post_post_item(self, request, response, item):
        user = request_get_user(request)
        if not app.is_seed:
            app.logger.error(item)
            post_internal('organization-users', {
                'organization': item['_id'],
                'role': 'owner',
                'user': user['_id'],
            })
            orga_statistics = post_internal('organization-statistics', {
                'organization': item['_id'],
            })

        # app.data.driver.db['organizations'].update(
        #     { '_id': organization['_id'] },
        #     { 'statistics': orga_statistics[0]['_id'] },
        # )



# Resource name / class mapping
models = {
    'users': UserItem,
    'user-tokens': UserTokenItem,
}


def resource_get_model(resource):
    """ Returns class matching resource name string. """
    return models.get(resource, BaseItem)
