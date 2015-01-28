var httpinvoke = require('httpinvoke'),
    debug = require('debug')('pathwar:lib'),
    config = require('./config'),
    _ = require('lodash');


var Client = module.exports = function(options) {
  this.config = _.defaults(options, config);
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
        _.startsWith(output._links.self.href, 'user-tokens/')) {
      client.config.token = output.token;
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
      });
    });
  };

  //this.logout = function(cb) {};

}).call(Client.prototype);
