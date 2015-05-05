# -*- coding: utf-8 -*-


achievements = {
    'item_title': 'achievement',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': ['GET'],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
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
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['moderator', 'admin', 'm1ch3l'],
    'allowed_write_roles': ['moderator', 'admin', 'm1ch3l'],
    # item
    'item_methods': ['GET', 'PATCH'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['moderator', 'admin', 'm1ch3l'],
    'allowed_item_write_roles': ['moderator', 'admin'],

    'schema': {
        'hash': {
            'type': 'string',
            'minlength': 3,
            'maxlength': 128,
            'unique': True,
            'required': True,
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
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['moderator', 'admin'],
    'allowed_write_roles': ['moderator', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['moderator', 'admin'],
    'allowed_item_write_roles': ['moderator', 'admin'],

    'schema': {
        'author': {
            'type': 'uuid',
            # 'required': True,
            'data_relation': {
                'resource': 'raw-users',
                'field': '_id',
                'embeddable': False,
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
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
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
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
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
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
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
        'pwn_status': {
            'type': 'string',
            'allowed': ['unknown', 'pwned'],
            'default': 'unknown',
        },
        'active': {
            'type': 'boolean',
            'default': True,
        },
        'server': {
            'type': 'uuid',
            'required': False,
            'nullable': True,
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
        'private_urls': {
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
        'name': {
            'type': 'string',
            'unique': True,
        }
    },
    'views': {
        'level-instances': {
            'datasource': {
                'source': 'raw-level-instances',
                'projection': {
                    '_schema_version': 0,
                    'passphrases': 0,
                    'server': 0,
                    'active': 0,
                    'private_urls': 0,
                },
                'filter': {
                    'active': True,
                },
            },
            'public_methods': [],
            'allowed_write_roles': [],
            'allowed_item_write_roles': [],
        },
        'hypervisor-level-instances': {
            'datasource': {
                'source': 'raw-level-instances',
                'projection': {
                    '_schema_version': 0,
                },
            },
            'schema': {
                'level': {
                    'type': 'uuid',
                    'required': True,
                    'data_relation': {
                        'resource': 'raw-levels',
                        'field': '_id',
                        # 'field': 'name',
                        'embeddable': True,
                    },
                },
            },
            # 'embedded_fields': ['level'],
            'public_methods': [],
            'allowed_read_roles': ['hypervisor', 'admin'],
            'allowed_item_read_roles': ['hypervisor', 'admin'],
            'allowed_write_roles': ['hypervisor', 'admin'],
            'allowed_item_write_roles': ['hypervisor', 'admin'],
        },
    },
}


level_instance_users = {
    'item_title': 'level instance user',
    'resource_title': 'level instance users',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'server', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH', 'DELETE'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'server', 'admin'],
    'allowed_item_write_roles': ['admin'],

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
        'expiry_date': {
            'type': 'datetime',
            'default': None,  # For tokens without expiry date
            # 'readonly': True,
            'nullable': True,
            'required': False,
        },
        'hash': {
            'type': 'string',
        },
    },
    'views': {
        'level-instance-users': {
            'datasource': {
                'source': 'raw-level-instance-users',
                'projection': {
                    '_schema_version': 0,
                },
            },
            'auth_field': 'user',
            'public_methods': [],
            'allowed_write_roles': ['user'],
            'allowed_item_write_roles': [],
        },
    },
}


level_statistics = {
    'item_title': 'level statistics',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
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
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

    'schema': {
        'name': {
            'type': 'string',
            'minlength': 3,
            'maxlength': 32,
            'unique': True,
            'required': True,
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
        'accessible': {
            'type': 'string',
            'allowed': ['until validated', 'forever'],
            'default': 'until validated',
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
            'default': 1000,
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
                        'resource': 'raw-users',
                        'field': '_id',
                        'embeddable': False,
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
    'views': {
        'levels': {
            'datasource': {
                'source': 'raw-levels',
                'projection': {
                    '_schema_version': 0,
                    'defaults': 0,
                    'author': 0,
                    'url': 0,
                },
            },
        },
    },
}


organization_achievements = {
    'item_title': 'organization earned achievement',
    'resource_title': 'organization earned achievements',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

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
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

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
                'resource': 'raw-users',
                'field': '_id',
                'embeddable': False,
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
    'views': {
        'organization-coupons': {
            'datasource': {
                'source': 'raw-organization-coupons',
                'projection': {
                    '_schema_version': 0,
                },
            },
            'public_methods': [],
            'allowed_write_roles': ['user'],
            'allowed_item_write_roles': [],
        },
    },
}


organization_level_validations = {
    'item_title': 'organization level validation submission',
    'resource_title': 'organization level validation submissions',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

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
                'resource': 'raw-users',
                'field': '_id',
                'embeddable': False,
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
    'views': {
        'organization-level-validations': {
            'datasource': {
                'source': 'raw-organization-level-validations',
                'projection': {
                    '_schema_version': 0,
                    'passphrases': 0,
                },
            },
            'public_methods': [],
            'allowed_write_roles': ['user', 'admin'],
            'allowed_item_write_roles': ['admin'],
        },
    },
}


organization_level_hints = {
    'item_title': 'organization level hint',
    'resource_title': 'organization level hints',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

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
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'admin'],
    'allowed_write_roles': ['user', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'admin'],
    'allowed_item_write_roles': ['admin'],

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
            'allowed': [
                'in progress', 'pending validation', 'validated', 'refused'
            ],
            'default': 'in progress'
        },
        'has_access': {
            'type': 'boolean',
            'default': True
        },
        'author': {
            'type': 'uuid',
            'required': False,
            'data_relation': {
                'resource': 'raw-users',
                'field': '_id',
                'embeddable': False,
            },
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
    'views': {
        'organization-levels': {
            'datasource': {
                'source': 'raw-organization-levels',
                'projection': {
                    '_schema_version': 0,
                    'author': 0,
                },
            },
            'public_methods': [],
            'allowed_write_roles': ['user'],
            'allowed_item_write_roles': [],
        },
    },
}


organization_items = {
    'item_title': 'organization item',
    'resource_title': 'organization items',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['admin'],
    'allowed_item_write_roles': ['admin'],

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
    'views': {
        'organization-items': {
            'datasource': {
                'source': 'raw-organization-items',
                'projection': {
                    '_schema_version': 0,
                },
            },
        },
    },
}


organization_users = {
    'item_title': 'organization user',
    'resource_title': 'organization users',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
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
                'resource': 'raw-users',
                'field': '_id',
                'embeddable': False,
            },
        },
    },
    'views': {
        'organization-users': {
            'datasource': {
                'source': 'raw-organization-users',
                'projection': {
                    '_schema_version': 0,
                },
            },
            'public_methods': [],
            'allowed_write_roles': ['user'],
            'allowed_item_write_roles': [],
        },
    },
}


organization_statistics = {
    'item_title': 'organization statistics',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
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
        'session': {
            'type': 'uuid',
            'required': False,
            'data_relation': {
                'resource': 'sessions',
                'field': '_id',
                'embeddable': False,
            },
        },
        'cash': {'type': 'integer', 'default': 0},
        'score': {'type': 'integer', 'default': 0},
        'gold_medals': {'type': 'integer', 'default': 0},
        'silver_medals': {'type': 'integer', 'default': 0},
        'bronze_medals': {'type': 'integer', 'default': 0},
        'achievements': {'type': 'integer', 'default': 0},
        'coupons': {'type': 'integer', 'default': 0},
    },
}


organizations = {
    'item_title': 'organization',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'moderator', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

    'schema': {
        'name': {
            'type': 'string',
            'required': True,
            'unique': True,
            'minlength': 3,
            'maxlength': 32,
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
                'resource': 'raw-users',
                'field': '_id',
                'embeddable': False,
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
    'views': {
        'organizations': {
            'datasource': {
                'source': 'raw-organizations',
                'projection': {
                    'gravatar_email': 0,
                    '_schema_version': 0,
                    'visibility': 0,
                    'owner': 0,
                },
                'filter': {
                    'visibility': 'public',
                },
            },
            'public_methods': [],
            'allowed_write_roles': ['admin'],
            'allowed_item_write_roles': ['admin'],
        },
        'teams': {
            # FIXME: enable user PATCH
            'datasource': {
                'source': 'raw-organizations',
                'projection': {
                    '_schema_version': 0,
                },
            },
        },
    },
}


sessions = {
    'item_title': 'session',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
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
        'anonymous': {
            'type': 'boolean',
            'default': False,
        },
        'email_domain': {
            'type': 'string',
        },
        'description': {
            'type': 'string',
        },
        'tags': {
            'type': 'list',
        },
        'avatar': {
            'type': 'string',
        },
    },
}


servers = {
    'item_title': 'server',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['server', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
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
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

    'schema': {
        'user': {
            'type': 'uuid',
            'data_relation': {
                'resource': 'raw-users',
                'field': '_id',
                'embeddable': False,
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
            'default': 'general',
        },
        'public': {
            'type': 'boolean',
            'default': False,
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
    'views': {
        'activities': {
            'datasource': {
                'source': 'raw-activities',
                'projection': {
                    '_schema_version': 0,
                    'user': 0,
                    'organization': 0,
                    'public': 0,
                },
                'filter': {
                    'public': True,
                },
            },
            'public_methods': [],
            'allowed_read_roles': ['user', 'moderator', 'admin', 'm1ch3l'],
            'allowed_item_read_roles': ['user', 'moderator', 'admin'],
            'allowed_write_roles': [],
            'allowed_item_write_roles': [],
        },
        'user-activities': {
            'datasource': {
                'source': 'raw-activities',
                'projection': {
                    '_schema_version': 0,
                },
            },
            'public_methods': [],
            'allowed_write_roles': [],
            'allowed_item_write_roles': [],
        },
    },
}


user_hijack_proofs = {
    'item_title': 'user hijack proof',
    'resource_title': 'user hijack proofs',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
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
                        'resource': 'raw-users',
                        'field': '_id',
                        'embeddable': False,
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
                        'resource': 'raw-users',
                        'field': '_id',
                        'embeddable': False,
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
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['user', 'admin'],

    'schema': {
        'user': {
            'type': 'uuid',
            'required': True,
            'data_relation': {
                # FIXME: switch back to embeddable users
                'resource': 'raw-users',
                'field': '_id',
                'embeddable': False,
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
        'author': {
            'type': 'uuid',
            'required': False,
            'data_relation': {
                'resource': 'users',
                'field': '_id',
                'embeddable': True,
            },
        },
        "status": {
            'type': 'string',
            'default': 'pending',
            'allowed': ['pending', 'accepted', 'refused'],
        },
    },
    'views': {
        'user-organization-invites': {
            'datasource': {
                'source': 'raw-user-organization-invites',
                'projection': {
                    '_schema_version': 0,
                },
            },
            'public_methods': [],
            'allowed_write_roles': ['user'],
            'allowed_item_write_roles': ['user'],
        },
        'admin-delete-user-organization-invites': {
            'datasource': {
                'source': 'raw-user-organization-invites',
                'projection': {
                    '_schema_version': 0,
                },
            },
            'public_methods': [],
            'resource_methods': ['GET'],
            'item_methods': ['GET', 'DELETE'],
            'allowed_write_roles': [],
            'allowed_item_write_roles': ['admin'],
            'allowed_read_roles': ['admin'],
            'allowed_item_read_roles': ['admin'],
        },
    },
}


user_notifications = {
    'item_title': 'user notification',
    'resource_title': 'user notifications',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['user', 'admin'],

    'cache_control': 'private, no-cache, no-store, must-revalidate',

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
                'resource': 'raw-users',
                'field': '_id',
                'embeddable': False,
            },
        },
        'action': {
            'type': 'string',
        },
        'category': {
            'type': 'string',
            'default': 'general',
        },
        'arguments': {
            'type': 'list',
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
    'views': {
        'user-notifications': {
            'datasource': {
                'source': 'raw-user-notifications',
                'projection': {
                    '_schema_version': 0,
                },
            },
            'auth_field': 'user',

            'public_methods': [],
            'allowed_write_roles': [],
            'allowed_item_write_roles': ['user'],
        },
    },
}


user_tokens = {
    'item_title': 'user token',
    'resource_title': 'user tokens',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'admin'],
    # item
    'item_methods': ['GET', 'DELETE', 'PATCH'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

    # ACL
    'auth_field': 'user',

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
        'is_admin': {
            'type': 'boolean',
            'default': False,
        },
        'user': {  # Will be computed using the credentials
            'type': 'uuid',
            'required': False,
            # 'readonly': True,
            'data_relation': {
                'resource': 'raw-users',
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
            # 'readonly': True,
            'nullable': True,
        },
    },
}


users = {
    'item_title': 'user',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': ['POST'],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'moderator', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

    'schema': {
        'myself': {
            'type': 'uuid',
            'readonly': True,
        },
        'last_login': {
            'type': 'string',
            'nullable': True,
            'readonly': True,
        },
        'login': {
            'type': 'string',
            'unique': True,
            'required': True,
            # chmod 646
            'minlength': 3,
            'maxlength': 32,
        },
        'email': {
            'type': 'string',
            'unique': True,
            'required': True,
            # chmod 606
        },
        'name': {
            'type': 'string',
            'maxlength': 64,
        },
        'website': {
            'type': 'string',
            'maxlength': 128,
        },
        'twitter_handle': {
            'type': 'string',
            'maxlength': 64,
        },
        'github_handle': {
            'type': 'string',
            'maxlength': 64,
        },
        'company': {
            'type': 'string',
            'maxlength': 128,
        },
        'location': {
            'type': 'string',
            'maxlength': 64,
        },
        'student_id': {
            'type': 'string',
            'maxlength': 64,
        },
        'active': {
            'type': 'boolean',
            'default': False,
            'readonly': True,
        },
        'blocked': {
            'type': 'boolean',
            'default': False,
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
            'allowed': [
                'user', 'level-server', 'moderator', 'admin', 'm1ch3l',
                'authproxy', 'hypervisor',
            ],
            'default': 'user',
            # chmod 600
        },
        'pnj': {
            'type': 'boolean',
            'default': False,
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
    },
    'views': {
        'users': {
            'datasource': {
                'source': 'raw-users',
                'projection': {
                    '_schema_version': 0,
                    'active': 0,
                    'email': 0,
                    'email_verification_token': 0,
                    'groups': 0,
                    'last_login': 0,
                    'myself': 0,
                    'otp_enabled': 0,
                    'otp_secret': 0,
                    'password': 0,
                    'password_salt': 0,
                    'student_id': 0,
                    'visibility': 0,
                    'pnj': 0,
                },
                'filter': {
                    'role': 'user',
                    'active': True,
                    # 'blocked': False,
                    # 'visibility': 'public',
                    # 'pnj': False,
                },
            },
            'public_methods': [],
            'allowed_read_roles': ['user', 'admin'],
            'allowed_item_read_roles': ['user', 'admin'],
            'allowed_item_write_roles': [],
        },
        'accounts': {
            'datasource': {
                'source': 'raw-users',
                'projection': {
                    'email_verification_token': 0,
                    'password_salt': 0,
                    'myself': 0,
                    'active': 0,
                    'pnj': 0,
                    '_schema_version': 0,
                }
            },
            'auth_field': 'myself',
            'allowed_read_roles': ['user', 'admin'],
            'allowed_item_read_roles': ['user', 'admin'],
            'allowed_write_roles': ['user', 'admin'],
            'allowed_item_write_roles': ['user', 'admin'],
        },
        'pnjs': {
            'datasource': {
                'source': 'raw-users',
                'projection': {
                    'email_verification_token': 0,
                    'password_salt': 0,
                    'password': 0,
                    '_schema_version': 0,
                },
                'filter': {
                    # 'pnj': True,
                },
            },
            'auth_field': 'myself',
            'allowed_read_roles': ['admin'],
            'allowed_item_read_roles': ['admin'],
            'allowed_write_roles': ['admin'],
            'allowed_item_write_roles': ['admin'],
        },
    },
}


password_recover_requests = {
    'item_title': 'password recover request',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': ['POST'],
    'allowed_read_roles': ['admin'],
    'allowed_write_roles': ['user'],
    # item
    'item_methods': ['GET'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['admin'],
    'allowed_item_write_roles': ['internal'],

    'schema': {
        'user': {
            'type': 'uuid',
            'required': True,
            'data_relation': {
                'resource': 'raw-users',
                'field': '_id',
                'embeddable': True,
            },
        },
        'verification_token': {
            'type': 'string',
            'readonly': True,
        },
        'status': {
            'type': 'string',
            'allowed': ['used', 'cancelled', 'pending'],
            'default': 'pending',
        },
        'password_salt': {
            'type': 'string',
            'readonly': True,
        },
        'password': {
            'required': True,
            'type': 'string',
        },
    },
}


whoswho_attempts = {
    'item_title': 'whoswho attempt',
    'resource_title': 'whoswho attempts',

    # collection
    'resource_methods': ['GET', 'POST'],
    'public_methods': [],
    'allowed_read_roles': ['user', 'moderator', 'admin'],
    'allowed_write_roles': ['user', 'admin'],
    # item
    'item_methods': ['GET', 'PATCH'],
    'public_item_methods': [],
    'allowed_item_read_roles': ['user', 'moderator', 'admin'],
    'allowed_item_write_roles': ['admin'],

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
                        'resource': 'raw-users',
                        'field': '_id',
                        'embeddable': False,
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
                        'resource': 'raw-users',
                        'field': '_id',
                        'embeddable': False,
                    },
                },
            },
        },
    },
}


DOMAIN = {}


BASE_RESOURCES = {
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
    'password-recover-requests': password_recover_requests,
    'servers': servers,
    'sessions': sessions,
    'user-hijack-proofs': user_hijack_proofs,
    'user-notifications': user_notifications,
    'user-organization-invites': user_organization_invites,
    'user-tokens': user_tokens,
    'users': users,
    'whoswho-attempts': whoswho_attempts,
}


for resource_name, resource in BASE_RESOURCES.items():
    if 'views' not in resource:
        resource['views'] = {}

    # Default raw view
    raw_name = 'raw-{}'.format(resource_name)
    if raw_name not in resource['views']:
        resource['views'][raw_name] = {
            'allowed_read_roles': ['admin'],
            'allowed_write_roles': ['admin'],
            'allowed_item_read_roles': ['admin'],
            'allowed_item_write_roles': ['admin'],
            'public_methods': [],
            'public_item_methods': [],
        }

    # Default base view
    if resource_name not in resource['views']:
        resource['views'][resource_name] = {
            'datasource': {
                'source': raw_name,
                'projection': {
                    '_schema_version': 0,
                },
            },
        }

    # Update resource
    resource['schema']['_id'] = {'type': 'uuid'}
    resource['schema']['_schema_version'] = {'type': 'integer'}
    uuid_regex = '[a-f0-9]{8}-?' \
                 '[a-f0-9]{4}-?' \
                 '4[a-f0-9]{3}-?' \
                 '[89ab][a-f0-9]{3}-?' \
                 '[a-f0-9]{12}'
    defaults = {
        'item_url': 'regex("{}")'.format(uuid_regex),
        'public_methods': [],
        'public_item_methods': [],
        'allowed_item_read_roles': ['admin'],
        'allowed_item_write_roles': ['admin'],
        'allowed_read_roles': ['admin'],
        'allowed_write_roles': ['admin'],
    }
    for key, value in defaults.items():
        if key not in resource:
            resource[key] = value

    # Register each views
    for view_name in resource['views'].keys():

        # if 'schema' in resource['views'][view_name]:
        #     view['schema'].update(resource['views'][view_name]['schema'])
        #     del resource['views'][view_name]['schema']

        view = resource.copy()
        view.update(resource['views'][view_name])
        # del view['views']
        DOMAIN[view_name] = view
