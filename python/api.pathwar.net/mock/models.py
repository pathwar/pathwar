achievements = {
    'item_title': 'achievement',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': [],
    'public_item_methods': [],
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
    'public_methods': [],
    'public_item_methods': [],
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
            'type': 'uuid',
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
    'public_methods': [],
    'public_item_methods': [],
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
    'public_methods': [],
    'public_item_methods': [],
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
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': [],
    'public_item_methods': [],
    # 'url': 'levels/<level>/instances',
    'schema': {
        'hash': {
            'type': 'string',
            'unique': True,
        },
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
        'server': {
            'type': 'uuid',
            'required': True,
            'data_relation': {
                'resource': 'servers',
                'field': '_id',
                'embeddable': True,
            },
        },
        'organizations': {
            'type': 'list',
            'schema': {
                'type': 'uuid',
                'data_relation': {
                    'resource': 'organizations',
                    'field': '_id',
                    'embeddable': True,
                },
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
    },
}


levels = {
    'item_title': 'level',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': [],
    'public_item_methods': [],
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
    'public_methods': [],
    'public_item_methods': [],
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
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': [],
    'public_item_methods': [],
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
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': [],
    'public_item_methods': [],
    # 'url': 'organizations/<organization>/levels/<level>/validations',
    'schema': {
        'status': {
            'type': 'string',
            'allowed': ['pending', 'accepted', 'refused'],
            'default': 'pending',
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
        'level': {
            'type': 'uuid',
            'required': True,
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
        'explanation': {
            'type': 'string',
        },
        'screenshot': {
            'type': 'string',
        },
    },
}


organization_levels = {
    'item_title': 'organization bought level',
    'resource_title': 'organization bought levels',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': [],
    'public_item_methods': [],
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
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': [],
    'public_item_methods': [],
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
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': [],
    'public_item_methods': [],
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


scorings = {
    'item_title': 'scoring',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': [],
    'public_item_methods': [],
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
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': [],
    'public_item_methods': [],
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
        'scoring': {
            'type': 'uuid',
            'required': False,
            'data_relation': {
                'resource': 'scorings',
                'field': '_id',
                'embeddable': True,
            },
        },
    },
}


sessions = {
    'item_title': 'session',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': ['GET'],
    'public_item_methods': ['GET'],
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
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': [],
    'public_item_methods': [],
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
    # FIXME: INTERNAL
    'item_title': 'user activitie',
    'resource_title': 'user activities',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': [],
    'public_item_methods': [],
    # 'url': 'users/<user>/activities',
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
                        'type': 'uuid',
                    },
                },
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
    },
}


user_organization_invites = {
    'item_title': 'user organization invite',
    'resource_title': 'user organization invites',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': [],
    'public_item_methods': [],
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
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': [],
    'public_item_methods': [],
    # 'url': 'users/<user>/notifications',
    'schema': {
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


user_sessions = {
    'item_title': 'user session',
    'resource_title': 'user sessions',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': [],
    'public_item_methods': [],
    # 'url': 'users/<user>/sessions',
    'schema': {
        'token': {
            'type': 'string',
            'default': 'random token',
            'unique': True,
        },
        # FIXME: date/expire
        'user': {
            'type': 'uuid',
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


user_tokens = {
    'item_title': 'user token',
    'resource_title': 'user tokens',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': [],
    'public_item_methods': [],
    # 'url': 'users/<user>/tokens',
    'schema': {
        'token': {
            'type': 'string',
            'default': 'FIXME: generate a random token',
            'unique': True,
            'required': False,
            'readonly': True,
        },
        'description': {
            'type': 'string',
            'required': False,
            'nullable': True,
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
        'scopes': {
            'type': 'string',
            'default': '*',
            'required': False,
            'empty': False,
            'nullable': False,
        },
        'expiry_date': {
            'type': 'datetime',
            'default': None,  # for token without expory date
            'required': False,
            'nullable': True,
        },
    },
}


users = {
    'item_title': 'user',
    'resource_methods': ['GET', 'POST', 'DELETE'],
    'item_methods': ['GET', 'PATCH', 'PUT', 'DELETE'],
    'public_methods': [],
    'public_item_methods': [],
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
        'active': {
            'type': 'boolean',
            'default': False,
        },
        'email_verification_token': {
            # INTERNAL
            'type': 'string',
            'required': False,
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
    'scorings': scorings,
    'servers': servers,
    'sessions': sessions,
    'user-activities': user_activities,
    'user-organization-invites': user_organization_invites,
    'user-notifications': user_notifications,
    'user-sessions': user_sessions,
    'user-tokens': user_tokens,
    'users': users,
}


# Use defaults
defaults = {
    'item_url': 'regex("[a-f0-9]{8}-?[a-f0-9]{4}-?4[a-f0-9]{3}-?' \
       '[89ab][a-f0-9]{3}-?[a-f0-9]{12}")',
    'public_methods': [],
    'public_item_methods': [],
}

for resource_name, resource_obj in DOMAIN.items():
    DOMAIN[resource_name]['schema']['_id'] = {'type': 'uuid'}
    for key, value in defaults.items():
        if key not in resource_obj:
            DOMAIN[resource_name][key] = value
