from eve import Eve
from settings import DOMAIN
app = Eve()

if __name__ == '__main__':
    # Initialize data
    with app.app_context():
        # Drop everything
        for collection in DOMAIN.keys():
            app.data.driver.db[collection].drop()

        # Initial content
        organization_id = app.data.driver.db['organizations'].insert({
            'name': 'example-organization',
        })
        user_id = app.data.driver.db['users'].insert({
            'login': 'example-user',
            'role': 'participant',
        })
        token_id = app.data.driver.db['user_tokens'].insert({
            'user': user_id,
        })
        level_id = app.data.driver.db['levels'].insert({
            'name': 'example-level',
        })
        #hint_id = app.data.driver.db['level_hints'].insert({
        #    'level': 'example-level',
        #    'name': 'example-level-hint',
        #})
        hint_id = app.data.driver.db['level_hints'].insert({
            'level': level_id,
            'name': 'example-level-hint-2',
        })

    # Run
    app.run(debug=True, host='0.0.0.0')
