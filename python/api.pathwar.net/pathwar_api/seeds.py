import datetime
from base64 import b64encode
from uuid import uuid4
import json

from settings import DOMAIN


def post(app, client, url, data, headers=None,
         content_type='application/json', auth_token='root-token'):
    if not headers:
        headers = []
    headers.append(('Content-Type', content_type))
    if auth_token:
        headers.append(('Authorization', 'Basic {}'
                        .format(b64encode('{}:'.format(auth_token)))))
    request = client.post(url, data=json.dumps(data), headers=headers)
    try:
        result = json.loads(request.get_data())
        if result['_items']:
            entries = result['_items']
        else:
            entries = [result]
        logger_method = app.logger.info
    except Exception, e:
        result = {}
        entries = []
        logger_method = app.logger.error
        app.logger.error(e)
    logger_method(request.get_data())
    logger_method("post({}, {}): {}, {}".
                  format(url, data, result.get('_status'),
                         result.get('message'),
                         result.get('_id')))
    return entries, result, request.status_code


def db_init(app):
    # Check if db is ready
    pass


def db_reset(app):
    app.logger.warn('DB_RESET')
    # Drop everything
    for collection in DOMAIN.keys():
        app.data.driver.db[collection].drop()

    # Create first user manually, the next requets (`post()`) will use the
    # token
    root_id = app.data.driver.db['users'].insert({
        'login': 'root',
        'role': 'admin',
        '_id': str(uuid4()),
        'active': True,
    })
    app.data.driver.db['user-tokens'].insert({
        'user': root_id,
        'token': 'root-token',
        '_id': str(uuid4()),
    })


def db_seed(app):
    app.logger.warn('DB_SEED')
    client = app.test_client()

    sessions, _, _ = post(
        app, client, '/sessions', [{
            'name': 'new year super challenge',
            'public': True,
        }, {
            'name': 'world battle',
        }])

    if sessions:
        users, _, _ = post(
            app, client, '/users', [{
                'login': 'joe',
                'email': 'joe@pathwar.net',
                'password': 'secure',
            }, {
                'login': 'm1ch3l',
                'email': 'm1ch3l@pathwar.net',
                'role': 'moderator',
                'active': True,
                'password': 'secure',
                'available_sessions': [
                    sessions[0]['_id'],
                    sessions[1]['_id'],
                ],
            }])
    lemming_users, _, _ = post(
        app, client, '/users', [{
            'login': 'lemming-{}'.format(i),
            'email': 'lemming-{}@lemming.net'.format(i),
            # 'password': 'secure',
            # 'active': True,
        } for i in xrange(200)]
    )

    # user_tokens, _, _ = post(app, client, '/user-tokens', [{
    #     'user': users[0]['_id'],
    # }, {
    #     'user': users[1]['_id'],
    #     #'expiry_date': ,
    # }])

    if sessions:
        coupons, _, _ = post(
            app, client, '/coupons', [{
                'hash': '1234567890',
                'value': 42,
                'session': sessions[0]['_id'],
            }, {
                'hash': '0987654321',
                'value': 24,
                'session': sessions[1]['_id'],
            }])

    servers, _, _ = post(
        app, client, '/servers', [{
            'name': 'c1-123',
            'ip_address': '1.2.3.4',
            'active': True,
            'token': '123456789000987654321',
            'tags': ['docker', 'armhf'],
        }, {
            'name': 'dedibox-152',
            'ip_address': '4.3.2.1',
            'active': True,
            'token': 'asodijgasodigjasoidgjoisadgj',
            'tags': ['docker', 'x86_64'],
        }])

    if sessions:
        organizations, _, _ = post(
            app, client, '/organizations', [{
                'name': 'pwn-around-the-world',
                'session': sessions[0]['_id'],
            }, {
                'name': 'staff',
                'session': sessions[1]['_id'],
            }])

    if sessions:
        lemming_organization, _, _ = post(
            app, client, '/organizations', [{
                'name': 'the-lemmings',
                'session': sessions[1]['_id'],
            }])

    if organizations:
        organization_statistics, _, _ = post(
            app, client, '/organization-statistics', [{
                'organization': organizations[0]['_id'],
                'cash': 42,
                'score': 42,
                'gold_medals': 3,
                'silver_medals': 3,
                'bronze_medals': 3,
                'achievements': 23,
            }, {
                'organization': organizations[1]['_id'],
                'cash': 42,
                'score': 42,
                'gold_medals': 3,
                'silver_medals': 3,
                'bronze_medals': 3,
                'achievements': 23,
            }])

    achievements, _, _ = post(
        app, client, '/achievements', [{
            'name': 'flash-gordon',
            'description': 'Validate a level in less than 1 minute',
        }, {
            'name': 'API',
            'description': 'Hack the API',
        }])

    if organizations and users:
        organizations_users, _, _ = post(
            app, client, '/organization-users', [{
                'organization': organizations[0]['_id'],
                'role': 'owner',
                'user': users[0]['_id'],
            }, {
                'organization': organizations[0]['_id'],
                'role': 'pwner',
                'user': users[1]['_id'],
            }, {
                'organization': organizations[1]['_id'],
                'role': 'owner',
                'user': users[1]['_id'],
            }])

    if lemming_users[0] and lemming_organization:
        for lemming_user in lemming_users:
            lemmings_family, _, _ = post(
                app, client, '/organization-users', {
                    'organization': lemming_organization[0]['_id'],
                    'role': 'pwner',
                    'user': lemming_user['_id']
                })

    levels, _, _ = post(
        app, client, '/levels', [{
            'name': 'welcome',
            'description': 'An easy welcome level',
            'price': 42,
            'tags': ['easy', 'welcome', 'official'],
            'author': 'Pathwar Team',
        }, {
            'name': 'pnu',
            'description': 'Possible not upload',
            'price': 420,
            'tags': ['php', 'advanced'],
            'author': 'Pathwar Team',
        }])

    if levels:
        level_hints, _, _ = post(
            app, client, '/level-hints', [{
                'name': 'welcome sources',
                'price': 42,
                'level': levels[0]['_id'],
            }, {
                'name': 'welcome full solution',
                'price': 420,
                'level': levels[0]['_id'],
            }, {
                'name': 'pnu sources',
                'price': 42,
                'level': levels[1]['_id'],
            }])

    if organizations and levels:
        organization_levels, _, _ = post(
            app, client, '/organization-levels', [{
                'organization': organizations[0]['_id'],
                'level': levels[0]['_id'],
            }, {
                'organization': organizations[0]['_id'],
                'level': levels[1]['_id'],
            }, {
                'organization': organizations[1]['_id'],
                'level': levels[1]['_id'],
            }])

    if organizations and achievements:
        organization_achievements, _, _ = post(
            app, client, '/organization-achievements', [{
                'organization': organizations[0]['_id'],
                'achievement': achievements[0]['_id'],
            }, {
                'organization': organizations[0]['_id'],
                'achievement': achievements[1]['_id'],
            }, {
                'organization': organizations[1]['_id'],
                'achievement': achievements[1]['_id'],
            }])

    if organizations and users:
        user_organization_invites, _, _ = post(
            app, client, '/user-organization-invites', [{
                'organization': organizations[0]['_id'],
                'user': users[0]['_id'],
            }, {
                'organization': organizations[0]['_id'],
                'user': users[1]['_id'],
            }, {
                'organization': organizations[1]['_id'],
                'user': users[1]['_id'],
            }])

    if organizations and levels and servers:
        level_instances, _, _ = post(
            app, client, '/level-instances', [{
                'hash': '123456789',
                'level': levels[0]['_id'],
                'server': servers[0]['_id'],
            # 'overrides': [{'key': 'cpu_shares', 'value': 42}],
            }, {
                'hash': '987654321',
                'level': levels[0]['_id'],
                'server': servers[1]['_id'],
                'organizations': [
                    organizations[0]['_id'],
                    organizations[1]['_id'],
                ],
            }, {
                'hash': '585185815',
                'level': levels[1]['_id'],
                'server': servers[1]['_id'],
            }])

    if users:
        user_notifications, _, _ = post(
            app, client, '/user-notifications', [{
                'user': users[0]['_id'],
                'title': 'hello !',
            }, {
                'user': users[0]['_id'],
                'title': 'what\s up?',
            }, {
                'user': users[1]['_id'],
                'title': 'hello !',
            }])

    if organizations and coupons:
        organization_coupons, _, _ = post(
            app, client, '/organization-coupons', [{
                'organization': organizations[0]['_id'],
                'coupon': coupons[0]['_id'],
            }, {
                'organization': organizations[0]['_id'],
                'coupon': coupons[1]['_id'],
            }, {
                'organization': organizations[1]['_id'],
                'coupon': coupons[1]['_id'],
            }])

    items, _, _ = post(
        app, client, '/items', [{
            'name': 'spiderpig-glasses',
            'description': 'Unlock all level hints',
            'price': 4242,
            'quantity': 1000,
        }, {
            'name': 'whoswho shield',
            'description': 'Cannot be attacked on whoswho',
            'price': 200,
            'quantity': 1000,
        }, {
            'name': 'for glory',
            'description': 'Unlock the "for glory" achievement',
            'price': 1000,
            'quantity': 1,
        }])

    if organizations and users and levels:
        activities, _, _ = post(
            app, client, '/activities', [{
                'user': users[0]['_id'],
                'category': 'level',
                'action': 'bought-level',
                'arguments': ['pnu'],
                'linked_resources': [{
                    'kind': 'levels',
                    'id': levels[1]['_id'],
                }],
                'organization': organizations[0]['_id'],
            }, {
                'user': users[0]['_id'],
                'action': 'created-token',
                'category': 'account',
            }, {
                'user': users[1]['_id'],
                'organization': organizations[1]['_id'],
                'category': 'whoswho',
                'action': 'whoswho-pwned',
                'arguments': ['pwn-around-the-world'],
                'linked_resources': [{
                    'kind': 'organizations',
                    'id': organizations[0]['_id'],
                }]
            }])

    if organizations and items:
        organization_items, _, _ = post(
            app, client, '/organization-items', [{
                'organization': organizations[0]['_id'],
                'item': items[0]['_id'],
            }, {
                'organization': organizations[0]['_id'],
                'item': items[1]['_id'],
            }, {
                'organization': organizations[1]['_id'],
                'item': items[1]['_id'],
            }])

    if organizations and levels and organization_levels:
        organization_level_validations, _, _ = post(
            app, client, '/organization-level-validations', [{
                'organization': organizations[0]['_id'],
                'level': levels[0]['_id'],
                'organization_level': organization_levels[0]['_id'],
                'status': 'pending',
                'explanation': 'blah blah',
                'screenshot': '/screenshots/1234567890',
            }, {
                'organization': organizations[0]['_id'],
                'level': levels[0]['_id'],
                'organization_level': organization_levels[1]['_id'],
                'status': 'accepted',
                'explanation': 'blih blih',
                'screenshot': '/screenshots/1234567891',
            }, {
                'organization': organizations[1]['_id'],
                'level': levels[0]['_id'],
                'organization_level': organization_levels[1]['_id'],
                'status': 'refused',
                'explanation': 'bloh bloh',
                'screenshot': '/screenshots/1234567892',
            }])
