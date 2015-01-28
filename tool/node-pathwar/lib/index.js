var httpinvoke = require('httpinvoke'),
    debug = require('debug')('pathwar:lib'),
    config = require('./config'),
    _ = require('lodash');


var Client = module.exports = function(options) {
  this.config = _.defaults(options, config);
};


httpinvoke = httpinvoke.hook('finished', function(err, output, statusCode, headers) {
  var ret;
  if(err) {
    return arguments;
  }
  if(typeof statusCode === 'undefined') {
    ret = new Error('Server or client error - undefined HTTP status');
    ret.output = output;
    ret.statusCode = statusCode;
    ret.headers = headers;
    return [ret , output, statusCode, headers];
  }
  if(statusCode >= 400 && statusCode <= 599) {
    ret = new Error(output._error.message + ' - undefined HTTP status');
    ret.output = output;
    ret.statusCode = statusCode;
    ret.headers = headers;
    return [ret, output, statusCode, headers];
  }
  return arguments;
});



(function() {
  this.request = function(path, method, options, cb) {
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
    return httpinvoke(url, method, options, cb);
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
}).call(Client.prototype);
