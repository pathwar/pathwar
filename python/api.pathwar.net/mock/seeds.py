from base64 import b64encode
from uuid import uuid4
import json

from settings import DOMAIN


def post(client, url, data, headers=None, content_type='application/json',
         auth_token='root-token'):
    if not headers:
        headers = []
    headers.append(('Content-Type', content_type))
    if auth_token:
        headers.append(('Authorization', 'Basic {}'
                        .format(b64encode('{}:'.format(auth_token)))))
    request = client.post(url, data=json.dumps(data), headers=headers)
    print(request.get_data())
    try:
        value = json.loads(request.get_data())
    except ValueError:
        value = {}
    print("post({}, {}): {}, {}".
          format(url, data, value.get('_status'), value.get('message'),
                 value.get('_id')))
    return value, request.status_code


def load_seeds(app, reset=True):
    if reset:
        # Drop everything
        for collection in DOMAIN.keys():
            app.data.driver.db[collection].drop()

    client = app.test_client()

    root_id = app.data.driver.db['users'].insert({
        'login': 'root',
        'role': 'admin',
        '_id': str(uuid4()),
    })
    app.data.driver.db['user-tokens'].insert({
        'user': root_id,
        'token': 'root-token',
        '_id': str(uuid4()),
    })

    users = post(client, '/users', [{
        'login': 'joe',
        'email': 'joe@pathwar.net',
        '_id': str(uuid4()),
    }, {
        'login': 'm1ch3l',
        'email': 'm1ch3l@pathwar.net',
        'role': 'superuser',
        '_id': str(uuid4()),
    }])

    sessions = post(client, '/sessions', [{
        'name': 'new year super challenge',
        '_id': str(uuid4()),
    }, {
        'name': 'world battle',
        '_id': str(uuid4()),
    }])

    organizations = post(client, '/organizations', [{
        'name': 'pwn-around-the-world',
        '_id': str(uuid4()),
    }, {
        'name': 'staff',
        '_id': str(uuid4()),
    }])

    organizations_users = post(client, '/organization-users', [{
        'organization': organizations[0]['_items'][0]['_id'],
        'role': 'owner',
        'user': users[0]['_items'][0]['_id'],
        '_id': str(uuid4()),
    }, {
        'organization': organizations[0]['_items'][0]['_id'],
        'role': 'pwner',
        'user': users[0]['_items'][1]['_id'],
        '_id': str(uuid4()),
    }, {
        'organization': organizations[0]['_items'][1]['_id'],
        'role': 'owner',
        'user': str(root_id),
        '_id': str(uuid4()),
    }])

    levels = post(client, '/levels', [{
        'name': 'welcome',
        'description': 'An easy welcome level',
        'price': 42,
        'tags': ['easy', 'welcome', 'official'],
        'author': 'Pathwar Team',
        '_id': str(uuid4()),
    }, {
        'name': 'pnu',
        'description': 'Possible not upload',
        'price': 420,
        'tags': ['php', 'advanced'],
        'author': 'Pathwar Team',
        '_id': str(uuid4()),
    }])

    level_hints = post(client, '/level-hints', [{
        'name': 'welcome sources',
        'price': 42,
        'level': levels[0]['_items'][0]['_id'],
        '_id': str(uuid4()),
    }, {
        'name': 'welcome full solution',
        'price': 420,
        'level': levels[0]['_items'][0]['_id'],
        '_id': str(uuid4()),
    }, {
        'name': 'pnu sources',
        'price': 42,
        'level': levels[0]['_items'][1]['_id'],
        '_id': str(uuid4()),
    }])
