from uuid import uuid4
import itertools
import json
import random
import md5

import bcrypt
from eve.methods.post import post, post_internal
from eve.methods.patch import patch_internal
from flask import abort, current_app, url_for

from utils import request_get_user, default_session
from mail import mail, send_mail


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

    @classmethod
    def find(cls, lookup):
        return list(current_app.data.driver.db[cls.resource].find(lookup))

    @classmethod
    def find_one(cls, lookup):
        return current_app.data.driver.db[cls.resource].find_one(lookup)

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


class OrganizationUserItem(BaseItem):
    resource = 'organization-users'


class UserItem(BaseItem):
    resource = 'users'

    @classmethod
    def get_by_organization_id(cls, organization_id):
        users_uuid = [
            str(member['user']) for member in OrganizationUserItem.find({
                'organization': organization_id,
            })
        ]
        return cls.find({
            '_id': {
                '$in': users_uuid,
            },
        })

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
            'name': '{}'.format(item['login']),
            'session': default_session()['_id'],
            'owner': item['_id'],
            'gravatar_email': item['email'],
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

    @classmethod
    def statistics_increment(cls, organization_id, payload):
        organization = cls.get_by_id(organization_id)
        OrganizationStatisticItem.update_by_id(
            organization['statistics'], {
                '$inc': payload,
            }
        )

    @classmethod
    def has_user(cls, organization_id, user_id):
        return OrganizationUserItem.find_one({
            'organization': organization_id,
            'user': user_id,
        })

    def on_pre_post_item(self, request, item):
        # FIXME: add a security check to ensure owner is preset by
        #        an internal commands, else drop it

        if 'owner' not in item:
            item['owner'] = request_get_user(request)['_id']

    def on_insert(self, item):
        super(OrganizationItem, self).on_insert(item)
        if 'gravatar_email' in item and item['gravatar_email']:
            item['gravatar_hash'] = md5.new(
                item['gravatar_email'].lower().strip()
            ).hexdigest()

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
        # Removing cash
        level = LevelItem.get_by_id(item['level'])
        if level['price']:
            OrganizationItem.statistics_increment(
                item['organization'], {
                    'cash': -level['price'],
                })

        # Create a notification for each members of the team
        members = UserItem.get_by_organization_id(item['organization'])
        for user in members:
            post_internal('user-notifications', {
                'title': 'Your team bought a new level',
                'user': user['_id'],
            })

        # FIXME: add transaction history for statistics recomputing

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


class OrganizationLevelValidationItem(BaseItem):
    resource = 'organization-level-validations'

    def on_pre_post_item(self, request, item):
        # Checking for mandatory fields
        if 'organization_level' not in item:
            abort(422, "Missing organization_level")
        if 'passphrases' not in item:
            abort(422, "Missing passphrases")
        if not isinstance(item['passphrases'], list) or \
           not len(item['passphrases']):
            abort(422, "Invalid type for passphrases")
        passphrases = map(str, item['passphrases'])
        current_app.logger.warn(passphrases)
        current_app.logger.warn(list(set(passphrases)))
        if sorted(list(set(passphrases))) != sorted(passphrases):
            abort(422, "Passphrases may be validated once")

        # FIXME: race condition, need an atomic update + fetch

        # Get OrganizationLevel from database
        organization_level = OrganizationLevelItem.get_by_id(
            item['organization_level'],
        )
        if not organization_level:
            abort(422, "No such organization_level")
        current_app.logger.info(
            'organization_level: {}'.format(organization_level)
        )

        # Check if the user validate a level for one if its organizations
        user = request_get_user(request)
        if not OrganizationItem.has_user(
                organization_level['organization'], user['_id']
        ):
            abort(422, "You cannot validate a coupon for another organization")

        # Add author
        item['author'] = user['_id']

        # Add computed organization and level uuids
        item['organization'] = organization_level['organization']
        item['level'] = organization_level['level']

        # FIXME: check if passphrase was already validated in another validation

        # Checking if passphrases are valid
        # FIXME: make the mongodb query filter more restrictive
        level_instances = LevelInstanceItem.find({
            'level': organization_level['level'],
        })
        available_passphrases = [
            passphrase['value']
            for passphrase in list(itertools.chain(*[
                    level_instance['passphrases']
                    for level_instance in level_instances
            ]))
        ]
        current_app.logger.warn('{}'.format(available_passphrases))
        for passphrase in passphrases:
            if passphrase not in available_passphrases:
                abort(422, "Bad passphrase")
        current_app.logger.info('level: {}'.format(available_passphrases))

    def on_inserted(self, item):
        # FIXME: compute all the validations and update the OrganizationLevel
        # FIXME: flag level instance as pwned -> redump if needed
        pass


class OrganizationStatisticItem(BaseItem):
    resource = 'organization-statistics'

    def on_inserted(self, item):
        OrganizationItem.update_by_id(item['organization'], {
            '$set': {
                'statistics': item['_id'],
            },
        })


class LevelItem(BaseItem):
    resource = 'levels'

    # FIXME: on_insert -> create LevelStatisticItem (as for Organization)


class LevelStatisticsItem(BaseItem):
    resource = 'level-statistics'

    # FIXME: mimic organizationstatistics


class LevelInstanceItem(BaseItem):
    resource = 'level-instances'


class CouponItem(BaseItem):
    resource = 'coupons'

    def on_insert(self, item):
        super(CouponItem, self).on_insert(item)
        item['validations_left'] = item['validations_limit']


class OrganizationCouponItem(BaseItem):
    resource = 'organization-coupons'

    def on_pre_post_item(self, request, item):
        if 'coupon' not in item:
            abort(422, "Missing coupon")

        coupon = CouponItem.find_one({
            'hash': item['coupon'],
        })

        if not coupon:
            abort(422, "No such coupon")

        # FIXME: race condition, need an atomic update + fetch

        if coupon['validations_left'] < 1:
            abort(422, "Expired coupon")

        # Check if the user add a coupon to one of its organizations
        user = request_get_user(request)
        if not OrganizationItem.has_user(item['organization'], user['_id']):
            abort(422, "You cannot validate a coupon for another organization")

        # Check if organization has already validated this coupon
        existing_coupon = OrganizationCouponItem.find_one({
            'coupon': coupon['_id'],
            'organization': item['organization'],
        })
        if existing_coupon:
            abort(422, 'You already validated this coupon')

        # Translate coupon name with its uuid
        item['coupon'] = coupon['_id']

        # Add author
        item['author'] = user['_id']

        # Decrease the validations_left
        CouponItem.update_by_id(
            coupon['_id'], {
                '$inc': {
                    'validations_left': -1,
                }
            }
        )

        # Removing cash
        OrganizationItem.statistics_increment(
            item['organization'], {
                'cash': coupon['value']
            })


# Resource name / class mapping
models = {
    'coupons': CouponItem,
    'levels': LevelItem,
    'level-instances': LevelInstanceItem,
    'level-statistics': LevelStatisticsItem,
    'organization-coupons': OrganizationCouponItem,
    'organization-levels': OrganizationLevelItem,
    'organization-level-validations': OrganizationLevelValidationItem,
    'organization-statistics': OrganizationStatisticItem,
    'organization-users': OrganizationUserItem,
    'organizations': OrganizationItem,
    'user-tokens': UserTokenItem,
    'users': UserItem,
}


def resource_get_model(resource):
    """ Returns class matching resource name string. """
    return models.get(resource, BaseItem)
