import os

from models import DOMAIN


# MAIL
if 'SMTP_PORT_25_TCP_ADDR' in os.environ:
    MAIL_SERVER = os.environ['SMTP_PORT_25_TCP_ADDR']
    MAIL_PORT = os.environ['SMTP_PORT_25_TCP_PORT']
MAIL_USE_TLS = False
MAIL_USE_SSL = False
MAIL_DEBUG = True
MAIL_USERNAME = None
MAIL_PASSWORD = None
# FIXME: try to put a name too
DEFAULT_MAIL_SENDER = 'notifications@pathwar.net'


# MONGO
MONGO_DBNAME = 'api-bench'
MONGO_HOST = os.environ['MONGO_PORT_27017_TCP_ADDR']
MONGO_PORT = os.environ['MONGO_PORT_27017_TCP_PORT']


# EVE DEFAULTS
PUBLIC_METHODS = ['GET']
PUBLIC_ITEM_METHODS = ['GET']


# CORS HEADERS
X_DOMAINS = '*'
X_HEADERS = ['Content-Type', 'If-Match', 'Authorization']
X_EXPOSE_HEADERS = ['Content-Length', 'Content-Type']


# FIXME: enable oplogs
