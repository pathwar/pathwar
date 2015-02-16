from uuid import uuid4
import json
import random
import md5

import bcrypt
from eve.methods.post import post, post_internal
from flask import abort, current_app

from utils import request_get_user, default_session


class BaseItem(object):
    def __init__(self):
        pass

    @classmethod
    def get_by_id(cls, uuid):
        return current_app.data.driver.db[cls.resource].find_one(
            {'_id': uuid}
        )

    @classmethod
    def update_by_id(cls, uuid, data):
        return current_app.data.driver.db[cls.resource].update(
            {'_id': uuid},
            data
        )

    def on_update(self, item):
        pass

    def on_insert(self, item):
        item['_id'] = str(uuid4())

    def on_inserted(self, item):
        pass

    def on_pre_get(self, request, lookup):
        pass

    def on_post_post_item(self, request, response, item):
        pass

    def on_pre_post_item(self, request, item):
        pass

    def on_pre_post(self, request):
        data = request.get_json()
        if isinstance(data, list):
            items = data
        else:
            items = [data]

        for item in items:
            self.on_pre_post_item(request, item)

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

    def _on_update(self, item):
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
        self._on_update(item)

    def on_inserted(self, item):
        post_internal('activities', {
            'user': item['_id'],
            'action': 'users-create',
            'category': 'account',
            'linked_resources': [
                {'kind': 'users', 'id': item['_id']}
            ],
        })
        post_internal('user-notifications', {
            'title': 'Welcome to your account !',
            'user': item['_id'],
        })

        # Create an organization in the default session
        default_organization = post_internal('organizations', {
            'name': 'default organization for {}'.format(item['login']),
            'session': default_session()['_id'],
            'owner': item['_id'],
        })

        # Send verification email
        if not current_app.is_seed and not item['active']:
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

    def on_pre_post_item(self, request, item):
        # FIXME: check for a password, users without password are built
        #        internally
        pass

    def on_pre_get(self, request, lookup):
        # Handle users/me
        if 'login' in lookup:
            del lookup['login']
            lookup['_id'] = request_get_user(request)['_id']


class UserTokenItem(BaseItem):
    resource = 'user-tokens'

    def on_pre_post_item(self, request, item):
        # Handle login
        user = request_get_user(request)

        if not user:
            abort(401)

        # FIXME: do not accept passing token/user (read-only)

        item['token'] = str(uuid4())
        item['user'] = user['_id']

        # FIXME: add expiry_date

    def on_inserted(self, item):
        post_internal('activities', {
            'user': item['user'],
            'action': 'user-tokens-create',
            'category': 'account',
            'linked_resources': [
                {'kind': 'users', 'id': item['user']},
                {'kind': 'user-tokens', 'id': item['_id']}
            ],
        })


class OrganizationItem(BaseItem):
    resource = 'organizations'

    def on_pre_post_item(self, request, item):
        # FIXME: add a security check to ensure owner is preset by
        #        an internal commands, else drop it

        if 'owner' not in item:
            item['owner'] = request_get_user(request)['_id']

    def on_inserted(self, item):
        post_internal('organization-users', {
            'organization': item['_id'],
            'role': 'owner',
            'user': item['owner'],
        })
        post_internal('organization-statistics', {
            'organization': item['_id'],
        })
        post_internal('activities', {
            'user': item['owner'],
            'organization': item['_id'],
            'action': 'organizations-create',
            'category': 'organizations',
            'linked_resources': [
                {'kind': 'organizations', 'id': item['_id']}
            ],
        })
        post_internal('user-notifications', {
            'title': 'You just created a new organization !',
            'user': item['owner'],
        })


class OrganizationLevelItem(BaseItem):
    resource = 'organization-levels'

    def on_inserted(self, item):
        organization = OrganizationItem.get_by_id(item['organization'])
        level = LevelItem.get_by_id(item['level'])

        # Removing cash
        if level['price']:
            OrganizationStatisticItem.update_by_id(
                organization['statistics'],
                {
                    '$inc': {
                        'cash': -level['price'],
                    }
                }
            )

        # FIXME: create notification for each organization members
        # FIXME: add transaction history for statistics computing

        # Add an activity
        post_internal('activities', {
            # 'user': item['owner'],
            'organization': item['organization'],
            'action': 'organization-levels-create',
            'category': 'levels',
            'linked_resources': [
                {'kind': 'organizations', 'id': item['organization']},
                {'kind': 'levels', 'id': item['level']},
                {'kind': 'organization-levels', 'id': item['_id']},
            ],
        })


class OrganizationStatisticItem(BaseItem):
    resource = 'organization-statistics'

    def on_inserted(self, item):
        OrganizationItem.update_by_id(item['organization'], {
            'statistics': item['_id'],
        })


class LevelItem(BaseItem):
    resource = 'levels'


# Resource name / class mapping
models = {
    'levels': LevelItem,
    'organization-levels': OrganizationLevelItem,
    'organization-statistics': OrganizationStatisticItem,
    'organizations': OrganizationItem,
    'user-tokens': UserTokenItem,
    'users': UserItem,
}


def resource_get_model(resource):
    """ Returns class matching resource name string. """
    return models.get(resource, BaseItem)
