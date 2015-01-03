import json

from eve import Eve
from settings import DOMAIN


app = Eve()


def post(client, url, data, headers=None, content_type='application/json'):
    if not headers:
        headers = []
    headers.append(('Content-Type', content_type))
    request = client.post(url, data=json.dumps(data), headers=headers)
    try:
        value = json.loads(request.get_data())
    except json.JSONDecodeError:
        value = None
    print("post({}): {}, {}".format(data, value.get('_status'), value.get('_id')))
    return value, request.status_code


if __name__ == '__main__':
    # Initialize data
    with app.app_context():
        # Drop everything
        for collection in DOMAIN.keys():
            app.data.driver.db[collection].drop()

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

        session = post(client, '/sessions', {
            'name': 'new year super challenge',
        })

        print('@' * 80)
        print(users[0]['_items'][0]['_id'])

        post(client, '/organizations', {
            'name': 'pwn-around-the-world',
            'users': [
                users[0]['_items'][0]['_id'],
            ],
        })


    # Run
    app.run(debug=True, host='0.0.0.0')
