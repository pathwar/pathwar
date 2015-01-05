from base64 import b64encode
import json


def post(client, url, data, headers=None, content_type='application/json',
         auth_token='root-token'):
    if not headers:
        headers = []
    headers.append(('Content-Type', content_type))
    if auth_token:
        headers.append(('Authorization', 'Basic {}'
                        .format(b64encode('{}:'.format(auth_token)))))
    request = client.post(url, data=json.dumps(data), headers=headers)
    try:
        value = json.loads(request.get_data())
    except json.JSONDecodeError:
        value = None
    print("post({}, {}): {}, {}".
          format(url, data, value.get('_status'), value.get('message'),
                 value.get('_id')))
    return value, request.status_code


def load_seeds(app):
    client = app.test_client()

    root_id = app.data.driver.db['users'].insert({
        'login': 'root',
        'role': 'admin',
    })
    app.data.driver.db['user-tokens'].insert({
        'user': root_id,
        'token': 'root-token',
    })

    users = post(client, '/users', [{
        'login': 'joe',
        'email': 'joe@pathwar.net',
    }, {
        'login': 'm1ch3l',
        'email': 'm1ch3l@pathwar.net',
        'role': 'superuser',
    }])

    sessions = post(client, '/sessions', [{
        'name': 'new year super challenge',
    }, {
        'name': 'world battle',
    }])

    organizations = post(client, '/organizations', [{
        'name': 'pwn-around-the-world',
    }, {
        'name': 'staff',
    }])

    organizations_users = post(client, '/organization-users', [{
        'organization': organizations[0]['_items'][0]['_id'],
        'role': 'owner',
        'user': users[0]['_items'][0]['_id'],
    }, {
        'organization': organizations[0]['_items'][0]['_id'],
        'role': 'pwner',
        'user': users[0]['_items'][1]['_id'],
    }, {
        'organization': organizations[0]['_items'][1]['_id'],
        'role': 'owner',
        'user': str(root_id),
    }])

    levels = post(client, '/levels', [{
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

    level_hints = post(client, '/level-hints', [{
        'name': 'welcome sources',
        'price': 42,
        'level': levels[0]['_items'][0]['_id'],
    }, {
        'name': 'welcome full solution',
        'price': 420,
        'level': levels[0]['_items'][0]['_id'],
    }, {
        'name': 'pnu sources',
        'price': 42,
        'level': levels[0]['_items'][1]['_id'],
    }])
