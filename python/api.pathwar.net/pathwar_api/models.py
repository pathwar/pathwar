from uuid import uuid4
import bcrypt
import datetime
import fnmatch
import itertools
import json
import md5
import os
import random
import re

import requests
from eve.methods.post import post, post_internal
from eve.methods.patch import patch_internal
from flask import abort, current_app, url_for, request as flask_request

from utils import request_get_user, generate_name, is_restricted_word
from mail import mail, send_mail
from resources import BASE_RESOURCES


class BaseModel(object):
    SCHEMA_VERSION = 1

    search_fields = ['_id']

    def __init__(self):
        pass

    @classmethod
    def base_resource(cls):
        return BASE_RESOURCES.get(cls.resource)

    @classmethod
    def base_views(cls):
        return cls.base_resource()['views']

    @classmethod
    def search(cls, search, full_object=False):
        query = {
            '$or': [
                {field: search} for field in cls.search_fields
            ]
        }

        if full_object:
            return cls.find(query)
        else:
            return [
                item['_id'] for item in
                cls.find(query, {'_id': 1})
            ]

    @classmethod
    def resolve_input(cls, input_, field, optional=False):
        search = input_.get(field)

        if not search:
            if optional:
                return False
            else:
                abort(422, "Missing {}".format(cls.resource))

        items = cls.search(search, full_object=True)
        if len(items) == 0:  # 0 matching item
            abort(422, "No such item '{}'".format(search))

        elif len(items) > 1:  # multiple matching items
            abort(422, "Too much candidates for item '{}'".format(search))

        else:  # len == 1
            input_[field] = items[0]['_id']
            return items[0]

    @classmethod
    def resolve_list(cls, list_):
        new_list = []
        for search in list_:
            items = cls.search(search)

            if len(items) == 0:  # 0 matching item
                abort(422, "No such item '{}'".format(search))

            elif len(items) > 1:  # multiple matching items
                abort(422, "Too much candidates for item '{}'".format(search))

            else:  # len == 1
                new_list.append(items[0])
        return new_list

    @classmethod
    def mongo_resource(cls):
        return 'raw-{}'.format(cls.resource)

    @classmethod
    def post_internal(cls, payload):
        return post_internal(cls.mongo_resource(), payload)

    @classmethod
    def get_by_id(cls, uuid):
        return current_app.data.driver.db[cls.mongo_resource()].find_one(
            {'_id': uuid}
        )

    @classmethod
    def update_by_id(cls, uuid, data):
        return current_app.data.driver.db[cls.mongo_resource()].update(
            {'_id': uuid},
            data
        )

    @classmethod
    def find(cls, lookup, projection=None):
        return list(current_app.data.driver.db[cls.mongo_resource()].find(
            lookup, projection
        ))

    @classmethod
    def find_one(cls, lookup, projection=None):
        return current_app.data.driver.db[cls.mongo_resource()].find_one(
            lookup, projection
        )

    def on_update(self, item, original):
        pass

    def on_updated(self, item, original):
        pass

    def on_insert(self, item):
        item['_id'] = str(uuid4())
        item['_schema_version'] = self.SCHEMA_VERSION

    def on_inserted(self, item):
        pass

    def on_pre_get(self, request, lookup):
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

    def on_pre_patch_item(self, request, item):
        pass

    def on_pre_patch(self, request, query):
        items = self.find(query)

        for item in items:
            self.on_pre_patch_item(request, item)


class Achievement(BaseModel):
    resource = 'achievements'

    @classmethod
    def unlock(cls, organization, achievements):
        for achievement_name in achievements:
            achievement = Achievement.find_one({
                'name': achievement_name
            })
            if not achievement:
                current_app.logger.error(
                    'Unknown achievement %s', achievement_name,
                )
                continue
            # FIXME: optimize
            if not OrganizationAchievement.find({
                'organization': organization,
                'achievement': achievement['_id'],
            }):
                OrganizationAchievement.post_internal({
                    'organization': organization,
                    'achievement': achievement['_id'],
                })


class Activity(BaseModel):
    resource = 'activities'

    def on_pre_get(self, request, lookup):
        if request.path.split('/')[1] == 'user-activities':
            user = request_get_user(request)
            organization_users = [
                organization_user['organization'] for organization_user in
                OrganizationUser.get_by_user(user['_id'])
            ]

            lookup['$or'] = [{
                'organization': {
                    '$in': organization_users,
                },
            }, {
                'user': user['_id'],
            }]


class OrganizationUser(BaseModel):
    resource = 'organization-users'

    # FIXME: rebuild the cross-references field User.organizations

    @classmethod
    def get_by_user(cls, user_id):
        return cls.find({
            'user': user_id,
        })

    def on_pre_get(self, request, lookup):
        if request.path.split('/')[1] == 'organization-users':
            user = request_get_user(request)
            organization_users = [
                organization_user['organization'] for organization_user in
                OrganizationUser.get_by_user(user['_id'])
            ]

            lookup['organization'] = {
                '$in': organization_users,
            }


class Session(BaseModel):
    resource = 'sessions'
    search_fields = ['_id', 'name']

    @classmethod
    def get_by_name(cls, name):
        return cls.find_one({'name': name})

    @classmethod
    def world_session(cls):
        return cls.get_by_name('World')

    @classmethod
    def beta_session(cls):
        return cls.get_by_name('Beta')


class PasswordRecoverRequest(BaseModel):
    resource = 'password-recover-requests'

    def on_pre_post_item(self, request, item):
        if 'password' not in item:
            abort(422, "Missing password")
        if 'email' not in item:
            abort(422, "Missing email")

        # match user
        user = User.resolve_input(item, 'email')
        item['user'] = user['_id']
        del item['email']

        # FIXME: check if user is blocked

        # Finally let inactive user send password recover requests
        # if not user['active']:
        #     abort(422, 'Account not yet validated')

        # check for captcha
        recaptcha_secret = os.environ.get('RECAPTCHA_SECRET')
        if recaptcha_secret:
            if 'captcha_response' not in item:
                abort(422, "Missing captcha_response")
            req = requests.post(
                'https://www.google.com/recaptcha/api/siteverify', data={
                    'secret': recaptcha_secret,
                    'response': item['captcha_response'],
                    # FIXME: 'remoteip': user_ip
                })
            if not req.json().get('success'):
                abort(422, "Invalid captcha")
        del item['captcha_response']

    def on_insert(self, item):
        super(PasswordRecoverRequest, self).on_insert(item)
        item['password_salt'] = bcrypt.gensalt().encode('utf-8')
        item['verification_token'] = str(uuid4())
        # crypt new password
        if 'password' in item and \
           len(item['password']):
            password = item['password'].encode('utf-8')
            item['password'] = bcrypt.hashpw(
                password, item['password_salt']
            )

    def on_inserted(self, item):
        Activity.post_internal({
            'user': item['_id'],
            'action': 'password-recover-request-create',
            'category': 'accounts',
            'public': False,
            'linked_resources': [
                {'kind': 'users', 'id': item['_id']},
            ],
        })
        UserNotification.post_internal({
            'title': 'You created a password recover request',
            'user': item['_id'],
            'action': 'password-recover-request-create',
            'category': 'accounts',
            'linked_resources': [
                {'kind': 'users', 'id': item['_id']},
            ],
        })

        # Send email
        verification_url = url_for(
            'tools.password_recover_verify',
            user_id=item['user'],
            verification_token=item['verification_token'],
            _external=True,
        )
        message = """
We received a request to reset the the password of this account. If you did not send it, please ignore this email.

If you want to get a new password for your account, please click or copy the link below.

    {}

Thanks,
Pathwar Team

If you received this email by mistake, simply delete it. You won't be subscribed if you don't click the confirmation link above.
""".format(verification_url)

        send_mail(
            message=message,
            subject='Password recover verification',
            recipients=[User.get_by_id(item['user'])]
        )


class User(BaseModel):
    resource = 'users'
    search_fields = ['_id', 'login', 'email']

    @classmethod
    def has_an_organization_for_session(cls, user, session):
        user_organizations = OrganizationUser.find({
            'user': user,
        })
        for user_organization in user_organizations:
            organization = Organization.get_by_id(
                user_organization['organization']
            )
            if organization['session'] == session:
                return True
        return False

    @classmethod
    def get_by_organization_id(cls, organization_id):
        users_uuid = [
            str(member['user']) for member in OrganizationUser.find({
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
        super(User, self).on_insert(item)
        item['password_salt'] = bcrypt.gensalt().encode('utf-8')
        item['email_verification_token'] = str(uuid4())
        item['myself'] = item['_id']
        # item['otp_secret'] = ...
        self._on_update(item)

    def on_inserted(self, item):
        Activity.post_internal({
            'user': item['_id'],
            'action': 'users-create',
            'category': 'accounts',
            'public': False,  # Pass to True when public
            'linked_resources': [
                {'kind': 'users', 'id': item['_id']},
            ],
        })
        UserNotification.post_internal({
            'title': 'Welcome to your account !',
            'user': item['_id'],
            'action': 'users-create',
            'category': 'accounts',
            'linked_resources': [
                {'kind': 'users', 'id': item['_id']},
            ],
        })

        # Create an organization in the default session
        default_organization = Organization.post_internal({
            'name': '{}'.format(item['login']),
            'session': Session.world_session()['_id'],
            'owner': item['_id'],
            'gravatar_email': item['email'],
        })

        # FIXME: automatically open subscriptions based on email pattern
        # matching

        # Send verification email
        if not current_app.is_seed and not item['active']:
            verification_url = url_for(
                'tools.email_verify',
                user_id=item['_id'],
                email_verification_token=item['email_verification_token'],
                _external=True,
            )
            message = """
You're almost done. Please, confirm your registration to Pathwar by clicking the link below.

    {}

Thanks,
Pathwar Team

If you received this email by mistake, simply delete it. You won't be subscribed if you don't click the confirmation link above.
""".format(verification_url)

            send_mail(
                message=message,
                subject='Email verification',
                recipients=[item]
            )

    def on_update(self, item, original):
        if 'email' in item and item['email'] != original['email']:
            abort(422, 'You cannot update your email')
        if 'blocked' in item and item['blocked'] != original.get('blocked'):
            user = request_get_user(flask_request)
            if user.get('role') != 'admin':
                abort(422, 'blocked field is read-only')
        if ('login' in item
            and item['login'] != original['login']):
            if original.get('last_login'):
                abort(422, "You can change your login only 1 time")

            if Organization.find_one({'name': item['login']}):
                abort(422, "Name already taken")

            if User.find_one({'login': item['login']}):
                abort(422, "Name already taken")

            world_organization = Organization.find_one({
                'owner': original['_id'],
                'session': Session.world_session()['_id'],
            })
            Organization.update_by_id(world_organization['_id'], {
                '$set': {
                    'name': item['login'],
                },
            })
            item['last_login'] = original['login']

    def on_pre_post_item(self, request, item):
        if 'password' not in item:
            abort(422, "Missing password")
        if 'login' not in item:
            abort(422, "Missing login")
        if 'email' not in item:
            abort(422, "Missing email")
        if is_restricted_word(item['login']):
            abort(422, "Invalid login")
        existing_user = User.find_one({
            '$or': [
                {'login': item['login']},
                {'email': item['email']},
            ]
        })
        if existing_user:
            abort(422, "Login or email already taken")

    def on_pre_get(self, request, lookup):
        # Handle users/me
        if 'login' in lookup:
            del lookup['login']
            lookup['_id'] = request_get_user(request)['_id']


class UserHijackProof(BaseModel):
    resource = 'user-hijack-proofs'


class UserNotification(BaseModel):
    resource = 'user-notifications'

    def on_inserted(self, item):
        pass


class UserOrganizationInvite(BaseModel):
    resource = 'user-organization-invites'

    def on_pre_get(self, request, lookup):
        if request.path.split('/')[1] == 'user-organization-invites':
            user = request_get_user(request)

            # lookup['status'] = 'pending'
            lookup['$or'] = [
                {'author': user['_id']},
                {'user': user['_id']},
            ]

    def on_pre_post_item(self, request, item):
        # Prepare input
        invited_user = User.resolve_input(item, 'user')
        organization = Organization.resolve_input(item, 'organization')

        # Author field
        user = request_get_user(request)
        item['author'] = user['_id']

        # Inactive users
        if not invited_user['active']:
            abort(422, 'The user is not yet validated')

        # Email restriction
        session = Session.get_by_id(organization['session'])
        if 'email_domain' in session and session['email_domain']:
            regex = fnmatch.translate(session['email_domain'])
            reobj = re.compile(regex)
            if not reobj.match(invited_user['email']):
                abort(422, 'This is a private session')

        # Forbid invite if non-owner
        if organization['owner'] != user['_id']:
            abort(422, "You are not owner of the organization")

        # Forbid self-invitation
        if user['_id'] == item['user']:
            abort(422, "You cannot invite yourself")

        # Forbid duplicates invitations
        existing_items = UserOrganizationInvite.find({
            'user': item['user'],
            'organization': item['organization'],
        })
        if len(existing_items):
            abort(422, 'You cannot invite someone twice')

        # Forbid invitation of someone already in a team in this session
        if User.has_an_organization_for_session(
            item['user'], organization['session']
        ):
            abort(
                422,
                'You cannot invite someone already in another organization'
            )

        # Forbid invitation in world session
        if organization['session'] == Session.world_session()['_id']:
            abort(422, 'You cannot invite someone else in your account')

    def on_inserted(self, item):
        UserNotification.post_internal({
            'title': 'New team invitation',
            'user': item['user'],
            'action': 'user-organization-invite-create',
            'category': 'organizations',
            'linked_resources': [
                {'kind': 'organizations', 'id': item['organization']},
                {'kind': 'author', 'id': item['author']},
            ],
        })

    def on_pre_patch_item(self, request, item):
        # FIXME: check if user is still solvable for accepting invite
        current = UserOrganizationInvite.get_by_id(
            request.path.split('/')[2]
        )
        user = request_get_user(request)
        if not current:
            abort(422, "No such user-organization-invite")

        if current['status'] != 'pending':
            request_user = request_get_user(flask_request)
            if user.get('role') != 'admin':
                abort(422, "Invitation expired")

        if current['user'] != user['_id']:
            abort(422, "This invitation was not for you")

        if item['status'] == 'accepted':
            organization = Organization.get_by_id([
                current['organization']
            ])
            if User.has_an_organization_for_session(
                user['_id'], organization['session']
            ):
                abort(422, "You already have an organization")

    def on_updated(self, item, original):
        if item['status'] != original['status']:
            if item['status'] == 'accepted':
                # Notify teamates
                members = User.get_by_organization_id(original['organization'])
                for user in members:
                    UserNotification.post_internal({
                        'title': 'New team member',
                        'user': user['_id'],
                        'action': 'user-organization-invite-accepted',
                        'category': 'organizations',
                        'linked_resources': [
                            {
                                'kind': 'organizations',
                                'id': original['organization'],
                            },
                            {'kind': 'users', 'id': original['user']},
                        ],
                    })

                # Create OrganizationUser
                OrganizationUser.post_internal({
                    'organization': original['organization'],
                    'role': 'pwner',
                    'user': original['user'],
                })

            # FIXME: post activity
            # FIXME: this kind of activity is private to the team
            elif item['status'] == 'refused':
                owner = Organization.get_by_id(original['organization'])['owner']
                UserNotification.post_internal({
                    'title': 'Your invitation was refused',
                    'user': owner,
                    'action': 'user-organization-invite-accepted',
                    'category': 'organizations',
                    'linked_resources': [
                        {'kind': 'organizations', 'id': original['organization']},
                        {'kind': 'users', 'id': original['user']},
                    ],
                })

    # FIXME: check if user is solvable (no existing organization,
    #        validated user, etc...)
    # FIXME: on PATCH by the user, add him to the new organization
    # FIXME: on POST, send user notification


class UserToken(BaseModel):
    resource = 'user-tokens'

    def on_pre_post_item(self, request, item):
        # Handle login
        user = request_get_user(request)
        if not user:
            abort(401)

        # FIXME: do not accept passing token/user (read-only)

        item['token'] = str(uuid4())
        item['user'] = user['_id']

        item['expiry_date'] = (
            datetime.datetime.utcnow() +
            datetime.timedelta(hours=12)
        )

    def on_inserted(self, item):
        Activity.post_internal({
            'user': item['user'],
            'action': 'user-tokens-create',
            'category': 'accounts',
            'public': False,
            'linked_resources': [
                {'kind': 'users', 'id': item['user']},
                {'kind': 'user-tokens', 'id': item['_id']}
            ],
        })


class Organization(BaseModel):
    resource = 'organizations'
    search_fields = ['_id', 'name']

    @classmethod
    def statistics_increment(cls, organization_id, payload):
        organization = cls.get_by_id(organization_id)
        OrganizationStatistics.update_by_id(
            organization['statistics'], {
                '$inc': payload,
            }
        )

    @classmethod
    def has_user(cls, organization_id, user_id):
        return OrganizationUser.find_one({
            'organization': organization_id,
            'user': user_id,
        })

    def on_pre_post_item(self, request, item):
        myself = request_get_user(request)
        session = Session.resolve_input(item, 'session')

        if 'name' not in item:
            abort(422, "Missing name")
        if is_restricted_word(item['name']):
            abort(422, "Invalid name")

        if 'email_domain' in session and session['email_domain']:
            regex = fnmatch.translate(session['email_domain'])
            reobj = re.compile(regex)
            if not reobj.match(myself['email']):
                abort(422, 'This is a private session')

        # Forbid invitation of someone already in a team in this session
        if User.has_an_organization_for_session(
                myself['_id'], session['_id']
        ):
            abort(
                422,
                'You cannot invite someone already in another organization'
            )

        # FIXME: add a security check to ensure owner is preset by
        #        an internal commands, else drop it

        if 'owner' not in item:
            item['owner'] = myself['_id']

    def on_insert(self, item):
        super(Organization, self).on_insert(item)
        if 'gravatar_email' in item and item['gravatar_email']:
            item['gravatar_hash'] = md5.new(
                item['gravatar_email'].lower().strip()
            ).hexdigest()

    def on_inserted(self, item):
        OrganizationUser.post_internal({
            'organization': item['_id'],
            'role': 'owner',
            'user': item['owner'],
        })
        OrganizationStatistics.post_internal({
            'organization': item['_id'],
        })
        Activity.post_internal({
            'user': item['owner'],
            'organization': item['_id'],
            'action': 'organizations-create',
            'category': 'organizations',
            'public': True,
            'linked_resources': [
                {'kind': 'organizations', 'id': item['_id']}
            ],
        })
        # UserNotification.post_internal({
        #     'title': 'You just created a new organization !',
        #     'user': item['owner'],
        #     'action': 'organizations-create',
        #     'category': 'organizations',
        #     'linked_resources': [
        #         {'kind': 'organizations', 'id': item['_id']}
        #     ],
        # })

    def on_pre_get(self, request, lookup):
        if request.path.split('/')[1] == 'teams':
            user = request_get_user(request)
            organization_users = [
                organization_user['organization'] for organization_user in
                OrganizationUser.get_by_user(user['_id'])
            ]
            if '_id' in lookup:
                if lookup['_id'] not in organization_users:
                    abort(401)
            else:
                lookup['_id'] = {'$in': organization_users}


class OrganizationLevel(BaseModel):
    resource = 'organization-levels'

    def on_pre_get(self, request, lookup):
        if request.path.split('/')[1] == 'organization-levels':
            user = request_get_user(request)
            organization_users = [
                organization_user['organization'] for organization_user in
                OrganizationUser.get_by_user(user['_id'])
            ]

            lookup['organization'] = {
                '$in': organization_users,
            }

    def on_pre_post_item(self, request, item):
        if 'organization' not in item:
            abort(422, "Missing organization")
        if 'level' not in item:
            abort(422, "Missing level")

        organization = Organization.get_by_id(item['organization'])
        if not organization:
            abort(422, "No such organization")

        level = Level.get_by_id(item['level'])
        if not level:
            abort(422, "No such level")

        sessions_availability = level.get('availability', {}).get('sessions')
        if (sessions_availability and
            organization['session'] not in sessions_availability):
            abort(422, "You cannot buy this level in this session")

        organization_statistics = OrganizationStatistics.get_by_id(
            organization['statistics'],
        )
        if not organization_statistics:
            abort(422, "No such organization_statistics")

        if organization_statistics['cash'] < level['price']:
            abort(422, "Not enough money")

        existing_items = OrganizationLevel.find({
            'level': item['level'],
            'organization': item['organization'],
        })
        if existing_items:
            abort(422, "You cannot buy two times a level")

        user = request_get_user(request)
        item['author'] = user['_id']

    def on_inserted(self, item):
        # Removing cash
        level = Level.get_by_id(item['level'])
        if level['price']:
            Organization.statistics_increment(
                item['organization'], {
                    'cash': -level['price'],
                })

        # Create a notification for each members of the team
        members = User.get_by_organization_id(item['organization'])
        for user in members:
            if user['_id'] == item['author']:
                continue
            UserNotification.post_internal({
                'title': 'New level bought',
                'user': user['_id'],
                'action': 'organization-level-create',
                'category': 'levels',
                'linked_resources': [
                    {'kind': 'levels', 'id': item['level']},
                    {'kind': 'organizations', 'id': item['organization']},
                ],
            })

        # FIXME: send notification to teamates)

        # Add an activity
        Activity.post_internal({
            # 'user': item['owner'],
            'organization': item['organization'],
            'action': 'organization-levels-create',
            'category': 'levels',
            'public': True,
            'linked_resources': [
                {'kind': 'organizations', 'id': item['organization']},
                {'kind': 'levels', 'id': item['level']},
            ],
        })

        # FIXME: move achievements computing into a dedicated function so we
        # can call it in a cronjob
        bought_levels = len(
            OrganizationLevel.find({'organization': item['organization']})
        )
        achievements = ['buy-1-level']
        if bought_levels >= 5:
            achievements.append('buy-5-levels')
        if bought_levels >= 10:
            achievements.append('buy-10-levels')
        if bought_levels >= 50:
            achievements.append('buy-50-levels')
        if bought_levels >= 100:
            achievements.append('buy-100-levels')
        Achievement.unlock(item['organization'], achievements)

    # def on_updated(self, item):
        # FIXME: add transaction history for statistics recomputing
        # FIXME: add ranking (for medals)
        # FIXME: check for achievements
        # FIXME: compute rewards


class OrganizationLevelValidation(BaseModel):
    resource = 'organization-level-validations'

    def on_pre_get(self, request, lookup):
        if request.path.split('/')[1] == 'organization-level-validations':
            user = request_get_user(request)
            organization_users = [
                organization_user['organization'] for organization_user in
                OrganizationUser.get_by_user(user['_id'])
            ]

            lookup['organization'] = {
                '$in': organization_users,
            }

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
        if sorted(list(set(passphrases))) != sorted(passphrases):
            abort(422, "Passphrases may be validated once")

        # FIXME: race condition, need an atomic update + fetch

        # Get OrganizationLevel from database
        organization_level = OrganizationLevel.get_by_id(
            item['organization_level'],
        )
        if not organization_level:
            abort(422, "No such organization_level")

        # Check if the user validate a level for one if its organizations
        user = request_get_user(request)
        if not Organization.has_user(
                organization_level['organization'], user['_id']
        ):
            abort(422, "You cannot validate a level for another organization")

        # Add author
        item['author'] = user['_id']

        # Add computed organization and level uuids
        item['organization'] = organization_level['organization']
        item['level'] = organization_level['level']

        # FIXME: check if passphrase was already validated in another
        #        validation

        # Checking if passphrases are valid
        # FIXME: make the mongodb query filter more restrictive
        level_instances = LevelInstance.find({
            'level': organization_level['level'],
        })
        available_passphrases = [
            passphrase['value']
            for passphrase in list(itertools.chain(*[
                level_instance['passphrases']
                for level_instance in level_instances
            ]))
        ]
        for passphrase in passphrases:
            if passphrase not in available_passphrases:
                abort(422, "Bad passphrase")

    def on_inserted(self, item):
        # FIXME: flag level instance as pwned -> redump if needed

        # First validation, give cash
        # FIXME: compute all the validations and check if _all_ passphrases
        #        are valid
        organization_level = OrganizationLevel.get_by_id(
            item['organization_level']
        )
        level = Level.get_by_id(organization_level['level'])
        if organization_level['status'] == 'in progress':
            Organization.statistics_increment(
                item['organization'], {
                    'cash': int(level['reward']),
                    'score': 10,
                })

        OrganizationLevel.update_by_id(item['organization_level'], {
            '$set': {
                'status': 'pending validation',
                'has_access': False,
            },
        })

        # FIXME: remove access to organization for this level

        members = User.get_by_organization_id(item['organization'])
        for user in members:
            if user['_id'] == item['author']:
                continue
            UserNotification.post_internal({
                'title': 'New level validation',
                'user': user['_id'],
                'action': 'organization-level-validation-create',
                'category': 'levels',
                'linked_resources': [
                    {
                        'kind': 'organizations',
                        'id': item['organization'],
                    },
                    {'kind': 'users', 'id': item['author']},
                    {'kind': 'levels', 'id': item['level']},
                ],
            })


class OrganizationLevelHint(BaseModel):
    resource = 'organization-level-hints'


class OrganizationStatistics(BaseModel):
    resource = 'organization-statistics'

    def _on_update(self, item):
        if 'session' not in item:
            organization = Organization.get_by_id(item['organization'])
            item['session'] = organization['session']

    def on_insert(self, item):
        super(OrganizationStatistics, self).on_insert(item)
        self._on_update(item)

    def on_inserted(self, item):
        Organization.update_by_id(item['organization'], {
            '$set': {
                'statistics': item['_id'],
            },
        })


class InfrastructureHijack(BaseModel):
    resource = 'infrastructure-hijacks'


class Item(BaseModel):
    resource = 'items'


class Level(BaseModel):
    resource = 'levels'

    def on_pre_post_item(self, request, item):
        if 'availability' in item:
            if 'sessions' in item['availability']:
                item['availability']['sessions'] = Session.resolve_list(
                    item['availability']['sessions']
                )

    def on_update(self, item, original):
        # FIXME: references update seems to be broken
        if 'availability' in item:
            if 'sessions' in item['availability']:
                item['availability']['sessions'] = Session.resolve_list(
                    item['availability']['sessions']
                )

    def on_inserted(self, item):
        LevelStatistics.post_internal({
            'level': item['_id'],
        })


class LevelStatistics(BaseModel):
    resource = 'level-statistics'

    def on_inserted(self, item):
        Level.update_by_id(item['level'], {
            '$set': {
                'statistics': item['_id'],
            },
        })


class LevelHint(BaseModel):
    resource = 'level-hints'


class LevelInstance(BaseModel):
    resource = 'level-instances'

    def on_pre_post_item(self, request, item):
        if 'name' not in item:
            item['name'] = generate_name()


class LevelInstanceUser(BaseModel):
    resource = 'level-instance-users'

    def on_pre_post_item(self, request, item):
        item['expiry_date'] = datetime.datetime.utcnow() + \
                              datetime.timedelta(hours=12)

        if 'level_instance' not in item:
            abort(422, "Missing level_instance")
        level_instance = LevelInstance.get_by_id(item['level_instance'])
        if not level_instance:
            abort(422, "No such level_instance")

        if 'organization' not in item:
            abort(422, "Missing organization")
        organization_level = OrganizationLevel.find({
            'organization': item['organization'],
            'level': level_instance['level'],
        })
        if not len(organization_level):
            abort(422, "No such organization_level")
        organization_level = organization_level[0]

        # FIXME: race condition, need an atomic update + fetch

        # Check if the user add a coupon to one of its organizations
        user = request_get_user(request)
        if not Organization.has_user(item['organization'], user['_id']):
            abort(422, "You cannot create object for another organization")

        # FIXME: Check if entry already exists, if yes, update the existing one

        # Add nested fields
        item['level'] = level_instance['level']
        item['organization_level'] = organization_level['_id']
        item['user'] = user['_id']

    def on_insert(self, item):
        super(LevelInstanceUser, self).on_insert(item)
        item['hash'] = str(uuid4())


class Coupon(BaseModel):
    resource = 'coupons'

    def on_insert(self, item):
        super(Coupon, self).on_insert(item)
        item['validations_left'] = item['validations_limit']


class OrganizationItem(BaseModel):
    resource = 'organization-items'


class OrganizationAchievement(BaseModel):
    resource = 'organization-achievements'

    def on_inserted(self, item):
        Activity.post_internal({
            'organization': item['organization'],
            'action': 'organization-achievement-create',
            'category': 'achievements',
            'public': True,
            'linked_resources': [
                {'kind': 'organizations', 'id': item['organization']},
                {'kind': 'achievements', 'id': item['achievement']},
            ],
        })

        Organization.statistics_increment(
            item['organization'], {
                'achievements': 1,
                'score': 5,
            })

        members = User.get_by_organization_id(item['organization'])
        for user in members:
            UserNotification.post_internal({
                'title': 'Achievement unlocked',
                'user': user['_id'],
                'action': 'organization-achievement-create',
                'category': 'achievements',
                'linked_resources': [
                    {
                        'kind': 'organizations',
                        'id': item['organization'],
                    },
                    {'kind': 'achievements', 'id': item['achievement']},
                ],
            })


class OrganizationCoupon(BaseModel):
    resource = 'organization-coupons'

    def on_pre_get(self, request, lookup):
        if request.path.split('/')[1] == 'organization-coupons':
            user = request_get_user(request)
            organization_users = [
                organization_user['organization'] for organization_user in
                OrganizationUser.get_by_user(user['_id'])
            ]

            lookup['organization'] = {
                '$in': organization_users,
            }

    def on_pre_post_item(self, request, item):
        if 'coupon' not in item:
            abort(422, "Missing coupon")

        coupon = Coupon.find_one({
            'hash': item['coupon'],
        })

        if not coupon:
            abort(422, "No such coupon")

        # FIXME: race condition, need an atomic update + fetch

        if coupon['validations_left'] < 1:
            abort(422, "Expired coupon")

        # Check if the user add a coupon to one of its organizations
        user = request_get_user(request)
        if not Organization.has_user(item['organization'], user['_id']):
            abort(422, "You cannot validate a coupon for another organization")

        # Check if organization has already validated this coupon
        existing_coupon = OrganizationCoupon.find_one({
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
        Coupon.update_by_id(
            coupon['_id'], {
                '$inc': {
                    'validations_left': -1,
                }
            }
        )

    def on_inserted(self, item):
        members = User.get_by_organization_id(item['organization'])
        for user in members:
            UserNotification.post_internal({
                'title': 'Coupon validated',
                'user': user['_id'],
                'action': 'organization-coupon-create',
                'category': 'coupons',
                'linked_resources': [
                    {
                        'kind': 'organizations',
                        'id': item['organization'],
                    },
                    {'kind': 'users', 'id': item['author']},
                    {'kind': 'organization-coupons', 'id': item['_id']},
                ],
            })

        coupon = Coupon.get_by_id(item['coupon'])

        # Update team cash
        Organization.statistics_increment(
            item['organization'], {
                'cash': coupon['value'],
                'coupons': 1,
            })

        # FIXME: move achievements computing into a dedicated function so we
        # can call it in a cronjob
        validated_coupons = len(
            OrganizationCoupon.find({
                'organization': item['organization'],
            })
        )
        achievements = ['validated-1-coupon']
        if validated_coupons >= 5:
            achievements.append('validated-5-coupons')
        if validated_coupons >= 10:
            achievements.append('validated-10-coupons')
        if validated_coupons >= 50:
            achievements.append('validated-50-coupons')
        if validated_coupons >= 100:
            achievements.append('validated-100-coupons')
        if validated_coupons >= 500:
            achievements.append('validated-500-coupons')
        if validated_coupons >= 1000:
            achievements.append('validated-1000-coupons')
        if validated_coupons >= 5000:
            achievements.append('validated-5000-coupons')
        Achievement.unlock(item['organization'], achievements)


class WhoswhoAttempt(BaseModel):
    resource = 'whoswho-attempts'


class Server(BaseModel):
    resource = 'servers'


# Resource name / class mapping
base_models = [
    Achievement,
    Activity,
    Coupon,
    InfrastructureHijack,
    Item,
    Level,
    Level,
    LevelHint,
    LevelInstance,
    LevelInstance,
    LevelInstanceUser,
    LevelStatistics,
    Organization,
    OrganizationAchievement,
    OrganizationCoupon,
    OrganizationItem,
    OrganizationLevel,
    OrganizationLevelHint,
    OrganizationLevelValidation,
    OrganizationStatistics,
    OrganizationUser,
    PasswordRecoverRequest,
    Server,
    Session,
    User,
    User,
    User,
    UserHijackProof,
    UserNotification,
    UserOrganizationInvite,
    UserToken,
    WhoswhoAttempt,
]


models = {}
for entry in base_models:
    for view_name, view in entry.base_views().items():
        models[view_name] = entry


def resource_get_model(resource):
    """ Returns class matching resource name string. """
    return models.get(resource, BaseModel)
