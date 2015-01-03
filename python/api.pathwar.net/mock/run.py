from eve import Eve
from settings import DOMAIN
app = Eve()

if __name__ == '__main__':
    with app.app_context():
        # Drop everything
        for collection in DOMAIN.keys():
            app.data.driver.db[collection].drop()

        organization_id = app.data.driver.db['organizations'].insert({
            'name': 'example-organization',
        })
        user_id = app.data.driver.db['users'].insert({
            'login': 'example-user',
            'role': 'participant',
        })
        level_id = app.data.driver.db['levels'].insert({
            'name': 'example-level',
        })
        hint_id = app.data.driver.db['hints'].insert({
            'level_id': level_id,
        })

    app.run(debug=True, host='0.0.0.0')
