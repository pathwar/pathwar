var httpinvoke = require('httpinvoke'),
    debug = require('debug')('pathwar:lib'),
    config = require('./config'),
    _ = require('lodash');


var Client = module.exports = function(config) {
  this.config = config;
};


(function() {
  this.get = function(path, options, cb) {
    var url = config.api_endpoint + path.replace(/^\//, '');
    options = options || {};
    _.defaults(options, {
      converters: {
        'text json': JSON.parse,
        'json text': JSON.stringify
      },
      headers: {},
      outputType: 'json'
    });
    _.defaults(options.headers, {
      Authorization: 'Basic ' + new Buffer(config.token + ':').toString('base64')
    });
    debug('GET ' + url, options);
    return httpinvoke(url, 'GET', options, cb);
  };
}).call(Client.prototype);
