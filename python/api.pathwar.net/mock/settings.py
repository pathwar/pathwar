import os


MONGO_DBNAME = 'api-bench'
MONGO_HOST = os.environ['MONGO_PORT_27017_TCP_ADDR']
MONGO_PORT = os.environ['MONGO_PORT_27017_TCP_PORT']


XML = False


achievements = {
    'item_title': 'achievement',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'schema': {
        'name': {
            'type': 'string',
            'unique': True,
        },
    },
}


coupons = {
    'item_title': 'coupon',
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


items = {
    'item_title': 'item',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'schema': {
        'name': {
            'type': 'string',
            'unique': True,
        },
    },
}


level_hints = {
    'item_title': 'level hint',
    'resource_title': 'level hints',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'url': 'levels/<level>/hints',
    'additional_lookup': {
        'url': 'string',
        'field': 'name',
    },
    'schema': {
        'name': {
            'type': 'string',
            # 'unique': True,
        },
        'level': {
            'type': 'objectid',
            'required': True,
            'data_relation': {
                'resource': 'levels',
                'field': '_id',
                # 'field': 'name',
                'embeddable': True,
            },
        },
    },
}


levels = {
    'item_title': 'level',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'additional_lookup': {
        'url': 'regex(".*")',
        'field': 'name',
    },
    'schema': {
        'name': {
            'type': 'string',
            'minlength': 3,
            'maxlength': 16,
            'unique': True,
        },
    },
}


organization_levels = {
    'item_title': 'organization level',
    'resource_title': 'organization levels',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'url': 'organizations/<organization>/levels',
    'schema': {
        'organization': {
            'type': 'objectid',
            'required': True,
            'data_relation': {
                'resource': 'organizations',
                'field': '_id',
                'embeddable': True,
            },
        },
        'level': {
            'type': 'objectid',
            'required': True,
            'data_relation': {
                'resource': 'levels',
                'field': '_id',
                'embeddable': True,
            },
        },
    },
}


organization_achievements = {
    'item_title': 'organization achievement',
    'resource_title': 'organization achievements',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'url': 'organizations/<organization>/achievements',
    'schema': {
        'organization': {
            'type': 'objectid',
            'required': True,
            'data_relation': {
                'resource': 'organizations',
                'field': '_id',
                'embeddable': True,
            },
        },
        'achievement': {
            'type': 'objectid',
            'required': True,
            'data_relation': {
                'resource': 'achievements',
                'field': '_id',
                'embeddable': True,
            },
        },
    },
}


organizations = {
    'item_title': 'organization',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'additional_lookup': {
        'url': 'regex(".*")',
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


sessions = {
    'item_title': 'session',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'schema': {
        'name': {
            'type': 'string',
            'unique': True,
        },
    },
}


user_notifications = {
    'item_title': 'user notification',
    'resource_title': 'user notifications',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'url': 'users/<user>/notifications',
    'schema': {
        'title': {
            'type': 'string',
        },
        'user': {
            'type': 'objectid',
            'required': True,
            'data_relation': {
                'resource': 'users',
                'field': '_id',
                'embeddable': True,
            },
        },
    },
}


user_tokens = {
    'item_title': 'user token',
    'resource_title': 'user tokens',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'url': 'users/<user>/tokens',
    'schema': {
        'token': {
            'type': 'string',
            'default': 'random token',
            'unique': True,
        },
        'user': {
            'type': 'objectid',
            'required': True,
            'data_relation': {
                'resource': 'users',
                'field': '_id',
                'embeddable': True,
            },
        },
    },
}


users = {
    'item_title': 'user',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'additional_lookup': {
        'url': 'regex(".*")',
        'field': 'login',
    },
    'schema': {
        'login': {
            'type': 'string',
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
    'level_hints': level_hints,
    'levels': levels,
    'organization_achievements': organization_achievements,
    'organization_levels': organization_levels,
    'organizations': organizations,
    'sessions': sessions,
    'user_notifications': user_notifications,
    'user_tokens': user_tokens,
    'users': users,
}
