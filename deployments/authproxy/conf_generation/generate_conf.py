#!/usr/bin/env python

import ConfigParser
import requests
import sys
import os
import json
import pprint
import time


CONF_FILE='conf.ini'

config = ConfigParser.ConfigParser()

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
    query = '{0}://{1}{2}'.format(config.get('conf', 'api_scheme'),
                                  config.get('conf', 'api_url'),
                                  endpoint)
    params = {}
    for arg in kwargs:
        if type(kwargs[arg]) == dict:
            params[arg]=json.dumps(kwargs[arg])
        else:
            params[arg]=kwargs[arg]
    r = requests.get(query, params=params, auth=(config.get('conf', 'api_token'),''))
    return r.json()


if len(sys.argv) == 2:
    CONF_FILE=sys.argv[1]
config.readfp(open(CONF_FILE))
ngx_tpl = open(config.get('conf','ngx_tpl')).read()

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
            server_name = server_name[7:] if server_name.startswith('http://') else server_name[8:]
            server_name = server_name[:server_name.find(':')]
            confs[k['_id']] = ngx_tpl
            confs[k['_id']] = confs[k['_id']].replace('_LEVEL_ID_', k['level']['_id']);
            confs[k['_id']] = confs[k['_id']].replace('_LISTEN_PORT_', k['urls'][0]['name']);
            confs[k['_id']] = confs[k['_id']].replace('_SERVER_NAME_', server_name);
            confs[k['_id']] = confs[k['_id']].replace('_LEVEL_URL_', k['server']['ip_address']);
        
    for id in confs:
        conf_name = '/' + id + '.conf'
        with open(config.get('conf','ngx_available_location') + conf_name, 'w') as fd:
            fd.write(confs[id])
        if id in active_levels:
            os.symlink(config.get('conf','ngx_available_location') + conf_name,
                       config.get('conf','ngx_enabled_location') + conf_name)
        else:
            try:
                os.unlink(config.get('conf','ngx_enabled_location') + conf_name)
            except:
                pass        
    os.system('/etc/init.d/nginx reload')
    time.sleep(60 * 5) #refresh conf every 5 minutes
