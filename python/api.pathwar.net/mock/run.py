from eve import Eve
from eve.auth import BasicAuth, TokenAuth

from settings import DOMAIN
from seeds import load_seeds


class MockBasicAuth(BasicAuth):
    def check_auth(self, username, password, allowed_roles, resource,
                   method):
        return username == 'admin' and password == 'secret'


class MockTokenAuth(TokenAuth):
    def check_auth(self, token, allowed_roles, resource, method):
        print('token', token)
        user_tokens = app.data.driver.db['user-tokens']
        return user_tokens.find_one({'token': token})


# Initialize Eve
app = Eve(
    # auth=MockBasicAuth,
    auth=MockTokenAuth,
)


def main():

    # Initialize data
    with app.app_context():
        # Drop everything
        for collection in DOMAIN.keys():
            app.data.driver.db[collection].drop()

        # Create basic content
        load_seeds(app)

    # Run
    app.run(
        debug=True,
        host='0.0.0.0',
    )


if __name__ == '__main__':
    main()
