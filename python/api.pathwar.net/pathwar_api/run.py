import sys

from flask.ext.bootstrap import Bootstrap
from eve_docs import eve_docs

from app import app
from settings import DOMAIN
from seeds import db_reset, db_init, db_seed
from tools import bp_tools
from mail import mail, send_mail
from hooks import setup_hooks


def eve_init():
    # eve-docs
    Bootstrap(app)
    app.register_blueprint(eve_docs, url_prefix='/docs')

    # tools
    app.register_blueprint(bp_tools, url_prefix='/tools')

    # mail
    mail.init_app(app)

    # hooks
    setup_hooks(app)

    # Initialize db
    db_init(app)


def main(argv):
    eve_init()

    if len(argv) > 1:
        if argv[1] == 'flush-db':
            with app.app_context():
                db_reset(app)
        elif argv[1] == 'seed-db':
            with app.app_context():
                db_reset(app)
                app.is_seed = True
                db_seed(app)
                app.is_seed = False

    else:
        # Run
        app.run(
            debug=True,
            host='0.0.0.0',
        )


if __name__ == '__main__':
    main(sys.argv)
