from flask import Blueprint, current_app, jsonify, abort

bp_tools = Blueprint('tools', __name__)


@bp_tools.route('/email-verify/<string:user_id>/'
                '<string:email_verification_token>')
def email_verify(user_id, email_verification_token):
    user = current_app.data.driver.db['users'].find_one({
        '_id': str(user_id),
        'active': False,
        'email_verification_token': str(email_verification_token),
    })
    if user:
        current_app.data.driver.db['users'].update(
            {
                '_id': str(user_id)
            }, {
                '$set': {
                    'active': True,
                    'email_verification_token': None,
                },
            },
        )
        # FIXME: 302 redirect to the portal
        return 'Email validated (FIXME: 302)'
    else:
        abort(404)
    print(user, user_id, email_verification_token)
