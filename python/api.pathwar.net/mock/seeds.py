import json


def post(client, url, data, headers=None, content_type='application/json'):
    if not headers:
        headers = []
    headers.append(('Content-Type', content_type))
    request = client.post(url, data=json.dumps(data), headers=headers)
    try:
        value = json.loads(request.get_data())
    except json.JSONDecodeError:
        value = None
    print("post({}): {}, {}".
          format(data, value.get('_status'), value.get('_id')))
    return value, request.status_code


def load_seeds(app):
    client = app.test_client()

    users = post(client, '/users', [{
        'login': 'joe',
        'email': 'joe@pathwar.net',
    }, {
        'login': 'm1ch3l',
        'email': 'm1ch3l@pathwar.net',
        'role': 'superuser',
    }, {
        'login': 'root',
        'email': 'root@pathwar.net',
        'role': 'admin',
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
        'user': users[0]['_items'][2]['_id'],
    }])
