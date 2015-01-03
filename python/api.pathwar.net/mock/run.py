import json

from eve import Eve
from settings import DOMAIN


app = Eve()


def post(client, url, data, headers=None, content_type='application/json'):
    if not headers:
        headers = []
    headers.append(('Content-Type', content_type))
    request = client.post(url, data=json.dumps(data), headers=headers)
    #return json.loads(request.get_data()), request.status_code


if __name__ == '__main__':
    # Initialize data
    with app.app_context():
        # Drop everything
        for collection in DOMAIN.keys():
            app.data.driver.db[collection].drop()

        client = app.test_client()

        post(client, '/users', [{
            'login': 'joe',
            'email': 'joe@pathwar.net',
            'role': 'user',
        }, {
            'login': 'm1ch3l',
            'email': 'm1ch3l@pathwar.net',
            'role': 'superuser',
        }, {
            'login': 'root',
            'email': 'root@pathwar.net',
            'role': 'admin',
        }])


    # Run
    app.run(debug=True, host='0.0.0.0')
