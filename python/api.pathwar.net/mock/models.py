achievements = {
    'item_title': 'achievement',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'schema': {
        'name': {
            'type': 'string',
            'unique': True,
        },
        'description': {
            'type': 'string',
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
        'value': {
            'type': 'integer',
            'default': 1,
        },
        'session': {
            'type': 'objectid',
            'required': True,
            'data_relation': {
                'resource': 'sessions',
                'field': '_id',
                'embeddable': True,
            },
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
        'description': {
            'type': 'string',
        },
        'price': {
            'type': 'integer',
        },
        'quantity': {
            'type': 'integer',
        },
    },
}


level_hints = {
    'item_title': 'level hint',
    'resource_title': 'level hints',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    # 'url': 'levels/<level>/hints',
    'additional_lookup': {
        'url': 'string',
        'field': 'name',
    },
    'schema': {
        'name': {
            'type': 'string',
            # 'unique': True,
        },
        'price': {
            'type': 'integer',
        },
        # FIXME: Add hint data: blob ?
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


level_instances = {
    'item_title': 'level instance',
    'resource_title': 'level instances',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    # 'url': 'levels/<level>/instances',
    'additional_lookup': {
        'url': 'string',
        'field': 'hash',
    },
    'schema': {
        'hash': {
            'type': 'string',
            'unique': True,
        },
        # FIXME: Add instance data: blob ?
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
        'server': {
            'type': 'objectid',
            'required': True,
            'data_relation': {
                'resource': 'servers',
                'field': '_id',
                'embeddable': True,
            },
        },
        'overrides': {
            'type': 'list',
            'schema': {
                'key': {
                    'type': 'string',
                    'allowed': [
                        'cpu_shares', 'memory_limit', 'redump', 'rootable'
                    ],
                },
                'value': {
                    'type': 'string',
                },
            },
        },
    },
}


levels = {
    'item_title': 'level',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'additional_lookup': {
        'url': 'string',
        'field': 'name',
    },
    'schema': {
        'name': {
            'type': 'string',
            'minlength': 3,
            'maxlength': 16,
            'unique': True,
        },
        'description': {
            'type': 'string',
        },
        'price': {
            'type': 'integer',
            'default': 4242,
        },
        'tags': {
            'type': 'list',
        },
        'registry_url': {
            'type': 'string',
        },
        'author': {
            'type': 'string',
            'default': 'Anonymous',
        },
        'passphrases_amount': {
            'type': 'integer',
            'default': 1,
        },
        'version': {
            'type': 'string',
            'default': 'dev',
        },
        'lang': {
            'type': 'string',
            'default': 'en',
        },
        'default_memory_limit': {
            'type': 'string',
            'default': '16M',
        },
        'default_cpu_shares': {
            'type': 'string',
            'default': '1/10',
        },
        'default_redump': {
            'type': 'integer',
            'default': 600,
        },
        'default_rootable': {
            'type': 'boolean',
            'default': True,
        },
    },
}


organization_achievements = {
    'item_title': 'organization earned achievement',
    'resource_title': 'organization earned achievements',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    # 'url': 'organizations/<organization>/achievements',
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


organization_coupons = {
    'item_title': 'organization validated coupon',
    'resource_title': 'organization validated coupons',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    # 'url': 'organizations/<organization>/coupons',
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
        'coupon': {
            'type': 'objectid',
            'required': True,
            'data_relation': {
                'resource': 'coupons',
                'field': '_id',
                'embeddable': True,
            },
        },
    },
}


organization_level_validations = {
    'item_title': 'organization level validation submission',
    'resource_title': 'organization level validation submissions',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    # 'url': 'organizations/<organization>/levels/<level>/validations',
    'schema': {
        'status': {
            'type': 'string',
            'allowed': ['pending', 'accepted', 'refused'],
            'default': 'pending',
        },
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
        'organization_level': {
            'type': 'objectid',
            'required': True,
            'data_relation': {
                'resource': 'organization-levels',
                'field': '_id',
                'embeddable': True,
            },
        },
    },
}


organization_levels = {
    'item_title': 'organization bought level',
    'resource_title': 'organization bought levels',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    # 'url': 'organizations/<organization>/levels',
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


organization_items = {
    'item_title': 'organization item',
    'resource_title': 'organization items',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    # 'url': 'organizations/<organization>/items',
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
        'item': {
            'type': 'objectid',
            'required': True,
            'data_relation': {
                'resource': 'items',
                'field': '_id',
                'embeddable': True,
            },
        },
    },
}


organization_users = {
    'item_title': 'organization item',
    'resource_title': 'organization items',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    # 'url': 'organizations/<organization>/users',
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
        'role': {
            'type': 'string',
            'allowed': ['pwner', 'owner'],
            'default': 'owner',
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


organizations = {
    'item_title': 'organization',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'additional_lookup': {
        'url': 'string',
        'field': 'name',
    },
    'schema': {
        'name': {
            'type': 'string',
            'unique': True,
        },
        'points': {
            'type': 'integer',
            'default': 0,
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
        'active': {
            'type': 'boolean',
            'default': True,
        },
        'allow_new_organizations': {
            'type': 'boolean',
            'default': True,
        },
        'allow_update_organizations': {
            'type': 'boolean',
            'default': True,
        },
    },
}


servers = {
    'item_title': 'server',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'schema': {
        'name': {
            'type': 'string',
            'unique': True,
        },
        'ip_address': {
            'type': 'string',
        },
        'active': {
            'type': 'boolean',
            'default': True,
        },
        'token': {
            'type': 'string',
        },
        'tags': {
            'type': 'list',
            'default': ['linux', 'public'],
        },
    },
}


user_activities = {
    'item_title': 'user activitie',
    'resource_title': 'user activities',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    # 'url': 'users/<user>/activities',
    'schema': {
        'user': {
            'type': 'objectid',
            'required': True,
            'data_relation': {
                'resource': 'users',
                'field': '_id',
                'embeddable': True,
            },
        },
        'action': {
            'type': 'string',
        },
        'arguments': {
            'type': 'list',
        },
        'category': {
            'type': 'string',
            'default': 'common',
        },
        'linked_resources': {
            'type': 'list',
            'schema': {
                'type': 'dict',
                'schema': {
                    'kind': {
                        'type': 'string',
                    },
                    'id': {
                        'type': 'objectid',
                    },
                },
            },
        },
        'organization': {
            'type': 'objectid',
            'required': False,
            'data_relation': {
                'resource': 'organizations',
                'field': '_id',
                'embeddable': True,
            },
        },
    },
}


user_organization_invites = {
    'item_title': 'user organization invite',
    'resource_title': 'user organization invites',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'schema': {
        'user': {
            'type': 'objectid',
            'required': True,
            'data_relation': {
                'resource': 'users',
                'field': '_id',
                'embeddable': True,
            },
        },
        'organization': {
            'type': 'objectid',
            'required': True,
            'data_relation': {
                'resource': 'organizations',
                'field': '_id',
                'embeddable': True,
            },
        },
    },
}


user_notifications = {
    'item_title': 'user notification',
    'resource_title': 'user notifications',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    # 'url': 'users/<user>/notifications',
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
    # 'url': 'users/<user>/tokens',
    'schema': {
        'token': {
            'type': 'string',
            'default': 'random token',
            'unique': True,
        },
        'description': {
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
        # FIXME: Add permissions, range, etc
    },
}


users = {
    'item_title': 'user',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'additional_lookup': {
        'url': 'string',
        'field': 'login',
    },
    'schema': {
        'login': {
            'type': 'string',
            'unique': True,
            'required': True,
        },
        'email': {
            'type': 'string',
            'unique': True,
            'required': True,
        },
        'password_blowfish': {
            'type': 'string',
            # 'required': True,
        },
        'otp_secret': {
            'type': 'string',
            # 'required': False,
        },
        'role': {
            'type': 'string',
            'allowed': ['user', 'superuser', 'admin'],
            'default': 'user',
            # 'required': True,
        },
        'location': {
            'type': 'dict',
            'schema': {
                'city': {'type': 'string'},
                'country': {'type': 'string'},
            },
        },
        'social_links': {
            'type': 'list',
            'schema': {
                'kind': {
                    'type': 'string',
                },
                'kind': {
                    'path': 'string',
                },
            },
        },
    },
}


DOMAIN = {
    # Exposed
    'achievements': achievements,
    'coupons': coupons,
    'items': items,
    'level-hints': level_hints,
    'level-instances': level_instances,
    'levels': levels,
    'organization-achievements': organization_achievements,
    'organization-coupons': organization_coupons,
    'organization-items': organization_items,
    'organization-level-validations': organization_level_validations,
    'organization-levels': organization_levels,
    'organization-users': organization_users,
    'organizations': organizations,
    'servers': servers,
    'sessions': sessions,
    'user-activities': user_activities,
    'user-organization-invites': user_organization_invites,
    'user-notifications': user_notifications,
    'user-tokens': user_tokens,
    'users': users,
}
