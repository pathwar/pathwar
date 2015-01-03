import os


MONGO_DBNAME = 'api-bench'
MONGO_HOST = os.environ['MONGO_PORT_27017_TCP_ADDR']
MONGO_PORT = os.environ['MONGO_PORT_27017_TCP_PORT']


achievements = {}


coupons = {
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'schema': {
        'hash': {
            'type': 'string',
            'minlength': 3,
            'maxlength': 32,
            'unique': True,
        },
    },
}


items = {}


levels = {}


organizations = {
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'additional_lookup': {
        'url': 'regex("[\w]+")',
        'field': 'name',
    },
    'schema': {
        'name': {
            'type': 'string',
            'minlength': 3,
            'maxlength': 16,
            'unique': True,
        },
        'points': {
            'type': 'integer',
        },
    },
}


sessions = {}


users = {
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'additional_lookup': {
        'url': 'regex("[\w]+")',
        'field': 'login',
    },
    'schema': {
        'login': {
            'type': 'string',
            'minlength': 3,
            'maxlength': 10,
            'unique': True,
        },
        'role': {
            'type': 'list',
            'allowed': ['participant', 'guest', 'admin']
        },
    },
}


DOMAIN = {
    'achievements': achievements,
    'coupons': coupons,
    'items': items,
    'levels': levels,
    'organizations': organizations,
    'sessions': sessions,
    'users': users,
}
