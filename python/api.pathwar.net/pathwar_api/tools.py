from flask import Blueprint, current_app, jsonify, abort

bp_tools = Blueprint('tools', __name__)


@bp_tools.route('/email-verify/<string:user_id>/'
                '<string:email_verification_token>')
def email_verify(user_id, email_verification_token):
    user = current_app.data.driver.db['raw-users'].find_one({
        '_id': str(user_id),
        'active': False,
        'email_verification_token': str(email_verification_token),
    })
    if user:
        current_app.data.driver.db['raw-users'].update(
            {
                '_id': str(user_id)
            }, {
                '$set': {
                    'active': True,
                    'email_verification_token': None,
                },
            },
        )
        return """<html>
  <head>
    <meta http-equiv="refresh" content="10; url=http://portal.pathwar.net/" />
  </head>
  <body>
    <p>Email validated, <a href="http://portal.pathwar.net/">go to portal</a>.</p>
  </body>
</html>"""

    else:
        abort(404)
    print(user, user_id, email_verification_token)
