import os

from resources import DOMAIN


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
MONGO_HOST = os.environ.get(
    'MONGO_HOST', os.environ.get('MONGO_PORT_27017_TCP_ADDR')
)
MONGO_DBNAME = os.environ.get('MONGO_DBNAME', 'pathwar')
MONGO_USERNAME = os.environ.get('MONGO_USERNAME')
MONGO_PASSWORD = os.environ.get('MONGO_PASSWORD')
MONGO_PORT = os.environ.get(
    'MONGO_PORT_27017_TCP_PORT', os.environ.get('MONGO_PORT', "27017")
)
MONGO_REPLICA_SET = os.environ.get('MONGO_REPLICA_SET')
# print(MONGO_HOST, MONGO_PORT, MONGO_DBNAME, MONGO_USERNAME, MONGO_PASSWORD)
# sys.exit(1)


# EVE DEFAULTS
PUBLIC_METHODS = []
PUBLIC_ITEM_METHODS = []


# CORS HEADERS
# X_DOMAINS = '*'
X_DOMAINS = os.environ.get('X_DOMAINS', '*').strip().split(',')
if len(X_DOMAINS) == 1:
    X_DOMAINS = X_DOMAINS[0]
X_HEADERS = ['Content-Type', 'If-Match', 'Authorization']
X_EXPOSE_HEADERS = ['Content-Length', 'Content-Type']


# FIXME: enable oplogs
