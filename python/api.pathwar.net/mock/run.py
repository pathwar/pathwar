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
        coupon_id = app.data.driver.db['coupons'].insert({
        })
        user_id = app.data.driver.db['users'].insert({
            'login': 'example-user',
            'role': 'participant',
        })
        item_id = app.data.driver.db['items'].insert({
        })
        achievement_id = app.data.driver.db['achievements'].insert({
            'name': 'pwn da world',
        })
        user_token_id = app.data.driver.db['user_tokens'].insert({
            'user': user_id,
        })
        notification_id = app.data.driver.db['user_notifications'].insert({
            'user': user_id,
            'title': 'Welcome',
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
        organization_level = app.data.driver.db['organization_levels'].insert({
            'level': level_id,
            'organization': organization_id,
        })
        organization_achievement = app.data.driver.db['organization_achievements'].insert({
            'achievement': achievement_id,
            'organization': organization_id,
        })
        organization_coupon = app.data.driver.db['organization_coupons'].insert({
            'coupon': coupon_id,
            'organization': organization_id,
        })
        organization_item = app.data.driver.db['organization_items'].insert({
            'item': item_id,
            'organization': organization_id,
        })
        organization_user = app.data.driver.db['organization_users'].insert({
            'user': user_id,
            'organization': organization_id,
        })
        organization_level_validation = app.data.driver.db['organization_level_validations'].insert({
            'level': level_id,
            'organization': organization_id,
        })

    # Run
    app.run(debug=True, host='0.0.0.0')
