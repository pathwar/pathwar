#!/usr/bin/env python

import requests
import sys
import os
import json
import pprint
import time
import socket

config = {}
config['api_scheme'] = os.environ.get('PATHWAR_API_SCHEME','http')
config['api_url'] = os.environ.get('PATHWAR_API_HOST','localhost:5000')
config['api_user'] = os.environ.get('PATHWAR_API_USER','default')
config['api_pass'] = os.environ.get('PATHWAR_API_PASS','')
config['ngx_tpl'] = os.environ.get('PATHWAR_NGX_CONF_PATH','/pathwar/conf_generation/nginx.tpl')
config['ngx_available_location'] = os.environ.get('PATHWAR_NGX_AVAILABLE','/etc/nginx/sites-available/')
config['ngx_enabled_location'] = os.environ.get('PATHWAR_NGX_ENABLED','/etc/nginx/sites-enabled/')

def daemonize():
    pid = os.fork()
    if pid > 0:
        exit(0)
    os.chdir('/')
    os.setsid()
    os.umask(0)
    pid = os.fork()
    if pid > 0:
        exit(0)

def api_request(endpoint, **kwargs):
    query = '{0}://{1}{2}'.format(config['api_scheme'],
                                  config['api_url'],
                                  endpoint)
    params = {}
    for arg in kwargs:
        if type(kwargs[arg]) == dict:
            params[arg]=json.dumps(kwargs[arg])
        else:
            params[arg]=kwargs[arg]
    r = requests.get(query, params=params, auth=(config['api_user'],config['api_pass']))
    return r.json()

ngx_tpl = open(config['ngx_tpl']).read()

daemonize()

while True:
    confs = {}
    active_levels = []

    links = api_request('/level-instances')['_links']
    last = int(links['last']['href'].split('?')[1].split('=')[1])
    for x in xrange(1,last+1):
        levels=api_request('/level-instances',embedded={'server':1, 'level':1},page=x)
        for k in levels['_items']:
            if k['active'] is True:
                active_levels.append(k['_id'])
            server_name = k['urls'][0]['url']
            if 'poney' in server_name:
                print k['_id']
            server_name = server_name[7:] if server_name.startswith('http://') else server_name[8:]
            server_name = server_name[:server_name.find(':')]
            confs[k['_id']] = ngx_tpl
            confs[k['_id']] = confs[k['_id']].replace('_LEVEL_ID_', k['level']['_id']);
            confs[k['_id']] = confs[k['_id']].replace('_LEVEL_INSTANCE_ID_', k['_id']);
            confs[k['_id']] = confs[k['_id']].replace('_LISTEN_PORT_', '80');
            confs[k['_id']] = confs[k['_id']].replace('_SERVER_NAME_', '{0}.{1}'.format(k['_id'], 'levels.pathwar.net'));
#            confs[k['_id']] = confs[k['_id']].replace('_LEVEL_URL_', k['server']['ip_address']);
            confs[k['_id']] = confs[k['_id']].replace('_LEVEL_URL_', server_name);
    pprint.pprint(confs)
    for id in confs:
        conf_name = '/' + id + '.conf'
        with open(config['ngx_available_location'] + conf_name, 'w') as fd:
            fd.write(confs[id])
        if id in active_levels:
            try:
                os.unlink(config['ngx_enabled_location'] + conf_name)
            except:
                pass
            os.symlink(config['ngx_available_location'] + conf_name,
                       config['ngx_enabled_location'] + conf_name)
        else:
            try:
                os.unlink(config['ngx_enabled_location'] + conf_name)
            except:
                pass        
    os.system('/usr/sbin/nginx -s reload')
    time.sleep(60 * 5) #refresh conf every 5 minutes
