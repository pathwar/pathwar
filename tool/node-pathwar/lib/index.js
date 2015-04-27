var rp = require('request-promise'),
    debug = require('debug')('pathwar:lib'),
    config = require('./config'),
    _ = require('lodash');


var Client = module.exports = function(options) {
  this.config = _.defaults(options || {}, config);
  this.user_id = null;
  this.organization_id = null;
  this.scope = '*';
};


(function() {
  var client;

  // FIXME: check what we can keep in this code
  var hook_finished = function(err, output, statusCode, headers) {
    // error handling
    var err_ret;
    if(err) {
      return arguments;
    }

    // login handling
    if (!client.config.token &&
        _.has(output, '_links') &&
        _.has(output._links, 'self') &&
        _.has(output._links.self, 'href') &&
        _.has(output, 'token') &&
        _.startsWith(output._links.self.href, 'user-tokens/')) {
      client.config.token = output.token;
      client.user_id = output.user._id;
      client.scope = output.scope;
      debug("Logged. user_id=" + client.user_id + ", scope=" + client.scope + ".");
    }

    // POST
    if (output._status === 'OK') {
      output.ids = _.chain(output._items)
        .map(function(obj) { return obj._id; })
        .value();
    }

    // Enrich fetched objects
    if (output._links && output._links.self && output._links.self.href) {
      var parts = output._links.self.href.split('?')[0].split('/');
      switch (parts.length) {
      case 1:  // resources collection
        output._ids = _.chain(output._items)
          .map(function(obj) { return obj._id; })
          .value();
        break;
      case 2:  // resource item
        // Nothing to do
        break;
      }
    }

    return arguments;
  };
  //this.httpinvoke = httpinvoke.hook('finished', hook_finished);

  // requests
  this.request = function(path, method, input, options, cb) {
    client = this;

    // options parameter is optional
    if (typeof options === 'function') {
      cb = options;
      options = {};
    }
    options = options || {};

    // build options
    _.defaults(options, {
      method: method,
      headers: {},
      resolveWithFullResponse: true
    });

    if (typeof(path) == 'object' && path._links != undefined) {
      var object = path;
      options.etag = options.etag || object._etag;
      path = object._links.self.href;
    }

    options.url = options.url || this.config.api_endpoint + path.replace(/^\//, '');

    if (options.etag) {
      _.defaults(options.headers, {
        'If-Match': options.etag
      });
    }

    // default headers
    // FIXME: add user-agent
    _.defaults(options.headers, {
      Accept: 'application/json'
    });

    // token-based authentication
    if (client.config.token) {
      _.defaults(options.headers, {
        Authorization: 'Basic ' + new Buffer(this.config.token + ':').toString('base64')
      });
    }

    // input is passed in the options object to rp, if input is empty
    // we need to use json:true to enable automatic response JSON parsong
    _.defaults(options, {
      json: input || true
    });

    debug(method + ' ' + options.url, options);

    // FIXME: debug response from server

    // FIXME: handle dry-run mode

    return rp(options).promise().nodeify(cb);
  };

  this.get = function(path, options, cb) {
    return this.request(path, 'GET', null, options, cb);
  };

  this.post = function(path, input, options, cb) {
    return this.request(path, 'POST', input, options, cb);
  };

  this.delete = function(path, options, cb) {
    return this.request(path, 'DELETE', null, options, cb);
  };

  this.patch = function(path, input, options, cb) {
    return this.request(path, 'PATCH', input, options, cb);
  };

  this.put = function(path, input, options, cb) {
    return this.request(path, 'PUT', input, options, cb);
  };

  // helpers
  this.login = function(username, password, cb) {
    debug("Logging in as " + username);
    var basicAuth = new Buffer(username + ':' + password).toString('base64');

    return this.post('/user-tokens', {
      is_session: true
    }, {
      headers: {
        Authorization: 'Basic ' + basicAuth
      }
    }).then(function(res) {
      var url = '/user-tokens/' + res.body._id + '?embedded={"user":1}';
      return client.get(url, {
        headers: {
          Authorization: 'Basic ' + basicAuth
        }
      }, cb);
    });
  };

  //FIXME: this.logout = function(cb) {};

  this.require_auth = function(cb) {
    if (!client.config.token && client.config.username && client.config.password) {
      return client.login(client.config.username, client.config.password, cb);
    }
    if (!client.user_id) {
      return client.get('/user-tokens/' + client.config.token, {}, cb);
    }
    return client.get('/', {}, cb);
  };

  this.organizations_list = function(cb) {
    return client.require_auth().then(
      function(res) {
        return client.get('/organization-users?where={"user":"'+client.user_id+'"}', {}, cb);
      }
    );
  };

  this.organizations_select = function(organization_id, cb) {
    return client.require_auth().then(
      function(res) {
        return client.get('/organizations/' + organization_id, {}, cb).then(
          function(res) {
            client.organization_id = res.body._id;
            debug('Selected organization ' + client.organization_id);
            return res;
          });
      }
    );
  };

  this.my_organization_users = function(cb) {
    return client.require_auth().then(
      function(res) {
        return client.get('/organization-users?where{"organization":"' + client.organization_id + '"}');
      }
    );
  };

  this.my_organization_levels = function(cb) {
    return client.require_auth().then(
      function(res) {
        return client.get('/organization-levels?where{"organization":"' + client.organization_id + '"}');
      }
    );
  };

}).call(Client.prototype);
