import os

#from models import DOMAIN

MONGO_DBNAME = 'api-bench'
MONGO_HOST = os.environ['MONGO_PORT_27017_TCP_ADDR']
MONGO_PORT = int(os.environ['MONGO_PORT_27017_TCP_PORT'])

PUBLIC_METHODS = ['GET']
PUBLIC_ITEM_METHODS = ['GET']

DOMAIN = {'eve-mongoengine': {}} # sadly this is needed for eve

# FIXME: enable oplogs
my_settings = {
    'MONGO_DBNAME': MONGO_DBNAME,
    'MONGO_HOST': MONGO_HOST,
    'MONGO_PORT': MONGO_PORT,
    'DOMAIN': DOMAIN,
}
