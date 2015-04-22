
achievements = {
    'item_title': 'achievement',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

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

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['moderator', 'admin'],
    'allowed_write_roles': ['moderator', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['moderator', 'admin'],
    'allowed_item_write_roles': ['moderator', 'admin'],

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
        'validations_limit': {
            'type': 'integer',
            'default': 1,
        },
        'validations_left': {
            'type': 'integer',
        },
        'session': {
            'type': 'uuid',
            # 'required': True,
            'data_relation': {
                'resource': 'sessions',
                'field': '_id',
                'embeddable': True,
            },
        },
    },
}


infrastructure_hijacks = {
    'item_title': 'infrastructure hijack',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['moderator', 'admin'],
    'allowed_write_roles': ['moderator', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['moderator', 'admin'],
    'allowed_item_write_roles': ['moderator', 'admin'],

    'schema': {
        'author': {
            'type': 'uuid',
            # 'required': True,
            'data_relation': {
                'resource': 'users',
                'field': '_id',
                'embeddable': True,
            },
        },
        'tags': {
            'type': 'list',
        },
        'description': {
            'type': 'string',
        },
    },
}


items = {
    'item_title': 'item',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

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

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

    # 'url': 'levels/<level>/hints',
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
            'type': 'uuid',
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

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

    # 'url': 'levels/<level>/instances',
    'schema': {
        # FIXME: Add instance data: blob ?
        'level': {
            'type': 'uuid',
            'required': True,
            'data_relation': {
                'resource': 'levels',
                'field': '_id',
                # 'field': 'name',
                'embeddable': True,
            },
        },
        'active': {
            'type': 'boolean',
            'default': True,
        },
        'server': {
            'type': 'uuid',
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
                'type': 'dict',
                # FIXME: repair the list of dict
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
        'urls': {
            'type': 'list',
            'schema': {
                'type': 'dict',
                'schema': {
                    'name': {
                        'type': 'string',
                    },
                    'url': {
                        'type': 'string',
                    },
                },
            },
        },
        'passphrases': {
            'type': 'list',
            'schema': {
                'type': 'dict',
                'schema': {
                    'key': {
                        'type': 'string',
                    },
                    'value': {
                        'type': 'string',
                    },
                },
            },
        },
    },
}


level_instance_users = {
    'item_title': 'level instance user',
    'resource_title': 'level instance users',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'server', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'server', 'admin'],
    'allowed_item_write_roles': ['user', 'admin'],

    'cache_control': 'private, no-cache, no-store, must-revalidate',

    'schema': {
        'level_instance': {
            'type': 'uuid',
            'required': True,
            'data_relation': {
                'resource': 'level-instances',
                'field': '_id',
                'embeddable': True,
            },
        },
        'level': {
            'type': 'uuid',
            'required': False,
            'data_relation': {
                'resource': 'levels',
                'field': '_id',
                'embeddable': True,
            },
        },
        'organization': {
            'type': 'uuid',
            'required': False,
            'data_relation': {
                'resource': 'organizations',
                'field': '_id',
                'embeddable': True,
            },
        },
        'organization_level': {
            'type': 'uuid',
            'required': True,
            'data_relation': {
                'resource': 'organization-levels',
                'field': '_id',
                'embeddable': True,
            },
        },
        'user': {
            'type': 'uuid',
            'required': False,
            'data_relation': {
                'resource': 'users',
                'field': '_id',
                'embeddable': True,
            },
        },
        'hash': {
            'type': 'string',
        },
    },
}


level_statistics = {
    'item_title': 'level statistics',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

    'cache_control': 'private, no-cache, no-store, must-revalidate',

    'schema': {
        'level': {
            'type': 'uuid',
            'required': True,
            'data_relation': {
                'resource': 'levels',
                'field': '_id',
                'embeddable': True,
            },
        },
        'amount_bought': {'type': 'integer', 'default': 0},
        'amount_finished': {'type': 'integer', 'default': 0},
        'fivestar_average': {'type': 'integer', 'default': 0},
        'duration_average': {'type': 'integer', 'default': 0},
        'amount_hints_bought': {'type': 'integer', 'default': 0},
    },
}


levels = {
    'item_title': 'level',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

    'schema': {
        'name': {
            'type': 'string',
            'minlength': 3,
            'maxlength': 16,
            'unique': True,
        },
        'statistics': {
            'type': 'uuid',
            'required': False,
            'data_relation': {
                'resource': 'level-statistics',
                'field': '_id',
                'embeddable': True,
            },
        },
        'availability': {
            'type': 'dict',
            'schema': {
                'sessions': {
                    'type': 'list',
                    'schema': {
                        'type': 'uuid',
                        'data_relation': {
                            'resource': 'sessions',
                            'field': '_id',
                            'embeddable': True,
                        },
                    },
                },
            },
        },
        'description': {
            'type': 'string',
        },
        'price': {
            'type': 'integer',
            'default': 1,
        },
        'tags': {
            'type': 'list',
        },
        'url': {
            'type': 'string',
        },
        'author': {
            'type': 'dict',
            'schema': {
                'string': {
                    'type': 'string',
                },
                'user': {
                    'type': 'uuid',
                    'data_relation': {
                        'resource': 'user',
                        'field': '_id',
                        'embeddable': True,
                    },
                },
                'organization': {
                    'type': 'uuid',
                    'data_relation': {
                        'resource': 'organizations',
                        'field': '_id',
                        'embeddable': True,
                    },
                },
            },
        },
        'passphrases_amount': {
            'type': 'integer',
            'default': 1,
        },
        'sources': {
            'type': 'dict',
            'schema': {
                'git': {
                    'type': 'string',
                },
                'build': {
                    'type': 'string',
                },
            },
        },
        'version': {
            'type': 'string',
            'default': 'dev',
        },
        'reward': {
            'type': 'integer',
            'default': 1
        },
        'difficulty': {
            'type': 'integer',
            'default': 1
        },
        'lang': {
            'type': 'string',
            'default': 'en',
        },
        'version': {
            'type': 'string',
            'default': 'latest',
        },
        'defaults': {
            'type': 'dict',
            'schema': {
                'memory_limit': {
                    'type': 'string',
                    'default': '16M',
                },
                'cpu_shares': {
                    'type': 'string',
                    'default': '1/10',
                },
                'redump': {
                    'type': 'integer',
                    'default': 600,
                },
                'rootable': {
                    'type': 'boolean',
                    'default': True,
                },
            },
        },
    },
}


organization_achievements = {
    'item_title': 'organization earned achievement',
    'resource_title': 'organization earned achievements',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['user', 'admin'],

    # 'url': 'organizations/<organization>/achievements',
    'schema': {
        'organization': {
            'type': 'uuid',
            'required': True,
            'data_relation': {
                'resource': 'organizations',
                'field': '_id',
                'embeddable': True,
            },
        },
        'achievement': {
            'type': 'uuid',
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

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['user', 'admin'],

    # 'url': 'organizations/<organization>/coupons',
    'schema': {
        'organization': {
            'type': 'uuid',
            'required': True,
            'data_relation': {
                'resource': 'organizations',
                'field': '_id',
                'embeddable': True,
            },
        },
        'author': {
            'type': 'uuid',
            'required': False,
            'data_relation': {
                'resource': 'users',
                'field': '_id',
                'embeddable': True,
            },
        },
        'coupon': {
            'type': 'uuid',
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

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['user', 'admin'],

    # 'url': 'organizations/<organization>/levels/<level>/validations',
    'schema': {
        'status': {
            'type': 'string',
            'allowed': ['pending', 'accepted', 'refused'],
            'default': 'pending',
        },
        'organization': {
            'type': 'uuid',
            'data_relation': {
                'resource': 'organizations',
                'field': '_id',
                'embeddable': True,
            },
        },
        'level': {
            'type': 'uuid',
            'data_relation': {
                'resource': 'levels',
                'field': '_id',
                'embeddable': True,
            },
        },
        'organization_level': {
            'type': 'uuid',
            'required': True,
            'data_relation': {
                'resource': 'organization-levels',
                'field': '_id',
                'embeddable': True,
            },
        },
        'author': {
            'type': 'uuid',
            'required': False,
            'data_relation': {
                'resource': 'users',
                'field': '_id',
                'embeddable': True,
            },
        },
        'explanation': {
            'type': 'string',
        },
        'screenshot': {
            'type': 'string',
        },
        'passphrases': {
            'required': True,
            'type': 'list',
            'schema': {
                'type': 'string',
            },
        },
    },
}


organization_level_hints = {
    'item_title': 'organization level hint',
    'resource_title': 'organization level hints',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['user', 'admin'],

    'schema': {
        'organization': {
            'type': 'uuid',
            'data_relation': {
                'resource': 'organizations',
                'field': '_id',
                'embeddable': True,
            },
        },
        'level': {
            'type': 'uuid',
            'data_relation': {
                'resource': 'levels',
                'field': '_id',
                'embeddable': True,
            },
        },
        'organization_level': {
            'type': 'uuid',
            'required': True,
            'data_relation': {
                'resource': 'organization-levels',
                'field': '_id',
                'embeddable': True,
            },
        },
        'level_hint': {
            'type': 'uuid',
            'required': True,
            'data_relation': {
                'resource': 'level-hints',
                'field': '_id',
                'embeddable': True,
            },
        },
    },
}


organization_levels = {
    'item_title': 'organization bought level',
    'resource_title': 'organization bought levels',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['user', 'admin'],

    # 'url': 'organizations/<organization>/levels',
    'schema': {
        'organization': {
            'type': 'uuid',
            'required': True,
            'data_relation': {
                'resource': 'organizations',
                'field': '_id',
                'embeddable': True,
            },
        },
        'rank': {
            'type': 'integer',
            'default': 0,
        },
        'status': {
            'type': 'string',
            'allowed': ['in progress', 'validated'],
            'default': 'in progress'
        },
        'level': {
            'type': 'uuid',
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

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['user', 'admin'],

    # 'url': 'organizations/<organization>/items',
    'schema': {
        'organization': {
            'type': 'uuid',
            'required': True,
            'data_relation': {
                'resource': 'organizations',
                'field': '_id',
                'embeddable': True,
            },
        },
        'item': {
            'type': 'uuid',
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
    'item_title': 'organization user',
    'resource_title': 'organization users',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

    # 'url': 'organizations/<organization>/users',
    'schema': {
        'organization': {
            'type': 'uuid',
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
            'type': 'uuid',
            'required': True,
            'data_relation': {
                'resource': 'users',
                'field': '_id',
                'embeddable': True,
            },
        },
    },
}


organization_statistics = {
    'item_title': 'organization statistics',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

    'cache_control': 'private, no-cache, no-store, must-revalidate',

    'schema': {
        'organization': {
            'type': 'uuid',
            'required': True,
            'data_relation': {
                'resource': 'organizations',
                'field': '_id',
                'embeddable': True,
            },
        },
        'cash': {'type': 'integer', 'default': 0},
        'score': {'type': 'integer', 'default': 0},
        'gold_medals': {'type': 'integer', 'default': 0},
        'silver_medals': {'type': 'integer', 'default': 0},
        'bronze_medals': {'type': 'integer', 'default': 0},
        'achievements': {'type': 'integer', 'default': 0},
    },
}


organizations = {
    'item_title': 'organization',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'moderator', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['user', 'moderator', 'admin'],

    'schema': {
        'name': {
            'type': 'string',
            'unique': True,
        },
        'session': {
            'type': 'uuid',
            'required': True,
            'data_relation': {
                'resource': 'sessions',
                'field': '_id',
                'embeddable': True,
            },
        },
        'visibility': {
            'type': 'string',
            'allowed': ['public', 'private', 'unlisted'],
            'default': 'public',
        },
        'owner': {
            'type': 'uuid',
            'required': False,
            'data_relation': {
                'resource': 'users',
                'field': '_id',
                'embeddable': True,
            },
        },
        'statistics': {
            'type': 'uuid',
            'required': False,
            'data_relation': {
                'resource': 'organization-statistics',
                'field': '_id',
                'embeddable': True,
            },
        },
        'gravatar_email': {
            'type': 'string',
            # 'unique': True,
            # 'required': True,
        },
        'gravatar_hash': {
            'type': 'string',
            'readonly': True,
        },
    },
}


sessions = {
    'item_title': 'session',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

    'schema': {
        'name': {
            'type': 'string',
            'unique': True,
        },
        'public': {
            'type': 'boolean',
            'default': False,
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

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['server', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['server', 'admin'],
    'allowed_item_write_roles': ['admin'],

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


activities = {
    # FIXME: INTERNAL
    'item_title': 'activity',
    'resource_title': 'activities',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

    'schema': {
        'user': {
            'type': 'uuid',
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
        'organization': {
            'type': 'uuid',
            'required': False,
            'data_relation': {
                'resource': 'organizations',
                'field': '_id',
                'embeddable': True,
            },
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
                        'type': 'uuid',
                    },
                },
            },
        },
    },
}


user_hijack_proofs = {
    'item_title': 'user hijack proof',
    'resource_title': 'user hijack proofs',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['user', 'admin'],

    'schema': {
        'status': {
            'type': 'string',
            'allowed': ['success', 'failure'],
            'default': 'pending'
        },
        'from': {
            'type': 'dict',
            'schema': {
                'organization': {
                    'type': 'uuid',
                    'required': True,
                    'data_relation': {
                        'resource': 'organizations',
                        'field': '_id',
                        'embeddable': True,
                    },
                },
                'author': {
                    'type': 'uuid',
                    'data_relation': {
                        'resource': 'users',
                        'field': '_id',
                        'embeddable': True,
                    },
                },
            },
        },
        'to': {
            'type': 'dict',
            'schema': {
                'organization': {
                    'type': 'uuid',
                    'required': True,
                    'data_relation': {
                        'resource': 'organizations',
                        'field': '_id',
                        'embeddable': True,
                    },
                },
                'author': {
                    'type': 'uuid',
                    'required': True,
                    'data_relation': {
                        'resource': 'users',
                        'field': '_id',
                        'embeddable': True,
                    },
                },
            },
        },
    },
}


user_organization_invites = {
    'item_title': 'user organization invite',
    'resource_title': 'user organization invites',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['user', 'admin'],

    'schema': {
        'user': {
            'type': 'uuid',
            'required': True,
            'data_relation': {
                'resource': 'users',
                'field': '_id',
                'embeddable': True,
            },
        },
        'organization': {
            'type': 'uuid',
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

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

    # 'url': 'users/<user>/notifications',
    'schema': {
        'read': {
            'type': 'boolean',
            'default': False,
        },
        'title': {
            'type': 'string',
        },
        'user': {
            'type': 'uuid',
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

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['user', 'admin'],

    # 'url': 'users/<user>/tokens',
    'schema': {
        'token': {
            'type': 'string',
            'default': 'FIXME: generate a random token',
            'unique': True,
            # 'readonly': True,
        },
        'description': {  # For the user
            'type': 'string',
            'default': '',
        },
        'is_session': {  # If true, the token will have an expiry date
            'type': 'boolean',
            'default': False,
        },
        'user': {  # Will be computed using the credentials
            'type': 'uuid',
            'required': False,
            # 'readonly': True,
            'data_relation': {
                'resource': 'users',
                'field': '_id',
                'embeddable': True,
            },
        },
        'scopes': {  # Access scope of the token
            'type': 'string',
            'default': '*',  # All access
            'required': False,
            'empty': False,
            'nullable': False,
        },
        'expiry_date': {  # Null expiry_date means no expiration
            'type': 'datetime',
            'default': None,  # For tokens without expiry date
            'readonly': True,
            'nullable': True,
        },
    },
}


users = {
    'item_title': 'user',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': ['POST'],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'moderator', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['user', 'moderator', 'admin'],

    'schema': {
        'login': {
            'type': 'string',
            'unique': True,
            'required': True,
            # chmod 646
        },
        'email': {
            'type': 'string',
            'unique': True,
            'required': True,
            # chmod 606
        },
        'active': {
            'type': 'boolean',
            'default': False,
            # chmod 644
        },
        'email_verification_token': {
            'type': 'string',
            'readonly': True,
            # chmod 002
        },
        'visibility': {
            'type': 'string',
            'allowed': ['public', 'private', 'unlisted'],
            'default': 'public',
        },
        'password_salt': {  # Generated on subscription
            'type': 'string',
            'readonly': True,
            # chmod 000
        },
        'password': {  # In reality, this will become a hash
            'type': 'string',
            # chmod 202
        },
        'otp_secret': {
            'type': 'string',
            'readonly': True,
            # chmod 202
        },
        'otp_enabled': {
            'type': 'boolean',
            'default': False,
            # chmod 646
        },
        'role': {
            'type': 'string',
            'allowed': ['user', 'level-server', 'moderator', 'admin'],
            'default': 'user',
            # chmod 600
        },
        'location': {
            'type': 'dict',
            'schema': {
                'city': {'type': 'string'},
                'country': {'type': 'string'},
            },
            # chmod 646
        },
        'social_links': {
            'type': 'list',
            'schema': {
                'type': 'dict',
                'schema': {
                    'kind': {
                        'type': 'string',
                    },
                },
            },
            # chmod 446
        },
        'gravatar_hash': {
            'type': 'string',
            'readonly': True,
        },
        'groups': {
            'type': 'list',
            'schema': {
                'type': 'string',
                'allowed': ['security', 'staff', 'developer', 'beta'],
            },
        },
        'available_sessions': {
            'type': 'list',
            'schema': {
                'type': 'uuid',
                'data_relation': {
                    'resource': 'sessions',
                    'field': '_id',
                    'embeddable': True,
                },
            },
            # chmod 644
        },
    },
}


whoswho_attempts = {
    'item_title': 'whoswho attempt',
    'resource_title': 'whoswho attempts',

    # collection
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['user', 'admin'],

    'schema': {
        'status': {
            'type': 'string',
            'allowed': ['success', 'failure'],
            'default': 'pending'
        },
        'from': {
            'type': 'dict',
            'schema': {
                'organization': {
                    'type': 'uuid',
                    'required': True,
                    'data_relation': {
                        'resource': 'organizations',
                        'field': '_id',
                        'embeddable': True,
                    },
                },
                'author': {
                    'type': 'uuid',
                    'data_relation': {
                        'resource': 'users',
                        'field': '_id',
                        'embeddable': True,
                    },
                },
            },
        },
        'to': {
            'type': 'dict',
            'schema': {
                'organization': {
                    'type': 'uuid',
                    'required': True,
                    'data_relation': {
                        'resource': 'organizations',
                        'field': '_id',
                        'embeddable': True,
                    },
                },
                'author': {
                    'type': 'uuid',
                    'required': True,
                    'data_relation': {
                        'resource': 'users',
                        'field': '_id',
                        'embeddable': True,
                    },
                },
            },
        },
    },
}


DOMAIN = {
    # Exposed
    'achievements': achievements,
    'activities': activities,
    'coupons': coupons,
    'infrastructure-hijacks': infrastructure_hijacks,
    'items': items,
    'level-hints': level_hints,
    'level-instance-users': level_instance_users,
    'level-instances': level_instances,
    'level-statistics': level_statistics,
    'levels': levels,
    'organization-achievements': organization_achievements,
    'organization-coupons': organization_coupons,
    'organization-items': organization_items,
    'organization-level-hints': organization_level_hints,
    'organization-level-validations': organization_level_validations,
    'organization-levels': organization_levels,
    'organization-statistics': organization_statistics,
    'organization-users': organization_users,
    'organizations': organizations,
    'servers': servers,
    'sessions': sessions,
    'user-hijack-proofs': user_hijack_proofs,
    'user-notifications': user_notifications,
    'user-organization-invites': user_organization_invites,
    'user-tokens': user_tokens,
    'users': users,
    'whoswho-attempts': whoswho_attempts,
}


# Use defaults
uuid_regex = '[a-f0-9]{8}-?' \
             '[a-f0-9]{4}-?' \
             '4[a-f0-9]{3}-?' \
             '[89ab][a-f0-9]{3}-?' \
             '[a-f0-9]{12}'
defaults = {
    'item_url': 'regex("{}")'.format(uuid_regex),
    'public_methods': [],
    'public_item_methods': [],
}

for resource_name, resource_obj in DOMAIN.items():
    DOMAIN[resource_name]['schema']['_id'] = {'type': 'uuid'}
    DOMAIN[resource_name]['schema']['_schema_version'] = {'type': 'integer'}
    for key, value in defaults.items():
        if key not in resource_obj:
            DOMAIN[resource_name][key] = value
