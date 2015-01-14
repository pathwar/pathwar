import os


from models import DOMAIN


MONGO_DBNAME = 'api-bench'
MONGO_HOST = os.environ['MONGO_PORT_27017_TCP_ADDR']
MONGO_PORT = os.environ['MONGO_PORT_27017_TCP_PORT']

PUBLIC_METHODS = ['GET']
PUBLIC_ITEM_METHODS = ['GET']

X_DOMAINS = '*'

# FIXME: enable oplogs
