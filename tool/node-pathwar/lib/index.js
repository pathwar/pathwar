var httpinvoke = require('httpinvoke'),
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

  // hooks
  var hook_finished = function(err, output, statusCode, headers) {
    // error handling
    var err_ret;
    if(err) {
      return arguments;
    }
    if(typeof statusCode === 'undefined') {
      err_ret = new Error('Server or client error - undefined HTTP status');
      err_ret.output = output;
      err_ret.statusCode = statusCode;
      err_ret.headers = headers;
      return [err_ret , output, statusCode, headers];
    }
    if(statusCode >= 400 && statusCode <= 599) {
      err_ret = new Error(output._error.message + ' - undefined HTTP status');
      err_ret.output = output;
      err_ret.statusCode = statusCode;
      err_ret.headers = headers;
      return [err_ret, output, statusCode, headers];
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
  this.httpinvoke = httpinvoke.hook('finished', hook_finished);

  // requests
  this.request = function(path, method, options, cb) {
    client = this;
    var url = this.config.api_endpoint + path.replace(/^\//, '');
    options = options || {};
    _.defaults(options, {
      partialOutputMode: 'joined',
      converters: {
        'text json': JSON.parse,
        'json text': JSON.stringify
      },
      headers: {},
      outputType: 'json'
    });
    _.defaults(options.headers, {
      Authorization: 'Basic ' + new Buffer(this.config.token + ':').toString('base64')
    });
    debug(method + ' ' + url, options);
    return this.httpinvoke(url, method, options, cb);
  };

  this.get = function(path, options, cb) {
    return this.request(path, 'GET', options, cb);
  };

  this.post = function(path, input, options, cb) {
    options = options || {};
    _.defaults(options, {
      inputType: 'json',
      headers: {},
      input: input
    });
    _.defaults(options.headers, {
      'Content-Type': 'application/json'
    });
    return this.request(path, 'POST', options, cb);
  };

  // helpers
  this.login = function(username, password, cb) {
    debug("Logging in as " + username);
    return this.post('/user-tokens', {
      is_session: true
    }, {
      headers: {
        Authorization: 'Basic ' + new Buffer(username + ':' + password).toString('base64')
      }
    }).then(function(res) {
      return client.get('/user-tokens/' + res.body._id + '?embedded={"user":1}', {
        headers: {
          Authorization: 'Basic ' + new Buffer(username + ':' + password).toString('base64')
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
