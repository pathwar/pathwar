from flask import current_app


def request_get_user(request):
    auth = request.authorization
    if auth.get('username'):
        if auth.get('password'):  # user:pass
            return current_app.data.driver.db['users'].find_one({
                'login': auth.get('username'),
                # 'active': True,  # FIXME: Reenable later
            })
        else:  # token
            user_token = current_app \
                .data \
                .driver \
                .db['user-tokens'] \
                .find_one({
                    'token': auth.get('username')
                })
            if user_token:
                return current_app \
                    .data \
                    .driver \
                    .db['users'] \
                    .find_one({'_id': user_token['user']})
    return None
