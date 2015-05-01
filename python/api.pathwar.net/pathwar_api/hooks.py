# -*- coding: utf-8 -*-


from flask import current_app

from app import app
from models import resource_get_model


def pre_get_callback(resource, request, lookup):
    """ Callback called before a GET request, we can alter the lookup. """
    klass = resource_get_model(resource)
    klass().on_pre_get(request, lookup)


def insert_callback(resource, items):
    """ Callback called just before inserting a resource in mongo. """
    # app.logger.debug('### insert_callback({}) {}'.format(resource, items))
    klass = resource_get_model(resource)
    for item in items:
        klass().on_insert(item)


def inserted_callback(resource, items):
    """ Callback called just after inserting a resource in mongo. """
    # app.logger.debug('### inserted_callback({}) {}'.format(resource, items))
    klass = resource_get_model(resource)
    for item in items:
        klass().on_inserted(item)


def update_callback(resource, item, original):
    """ Callback called just before updating a resource in mongo. """
    # app.logger.debug('### update_callback({}) {}'.format(resource, items))
    klass = resource_get_model(resource)
    klass().on_update(item, original)


def updated_callback(resource, item, original):
    """ Callback called just after updating a resource in mongo. """
    # app.logger.debug('### updated_callback({}) {}'.format(resource, items))
    klass = resource_get_model(resource)
    klass().on_updated(item, original)


def pre_post_callback(resource, request):
    """ Callback called just before the normal processing behavior of a POST
    request.
    """
    klass = resource_get_model(resource)
    klass().on_pre_post(request)


def post_post_callback(resource, request, response):
    """ Callback called just after a POST request ended. """
    # app.logger.info('### post_post({}) request: {}, response: {}'
    #                 .format(resource, request, response))
    klass = resource_get_model(resource)
    klass().on_post_post(request, response)


def pre_patch_callback(resource, request, query):
    """ Callback called just before the normal processing behavior of a PATCH
    request.
    """
    klass = resource_get_model(resource)
    klass().on_pre_patch(request, query)


def setup_hooks(_app):
    # Attach hooks
    app.on_pre_GET += pre_get_callback
    app.on_insert += insert_callback
    app.on_inserted += inserted_callback
    app.on_update += update_callback
    app.on_updated += updated_callback
    app.on_pre_POST += pre_post_callback
    app.on_post_POST += post_post_callback
    app.on_pre_PATCH += pre_patch_callback
    # getattr(app, 'on_pre_POST_user-tokens') += pre_post_user_tokens_callback
