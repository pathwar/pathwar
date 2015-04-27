var Api = require('./index'),
    Table = require('cli-table'),
    debug = require('debug')('pathwar:utils'),
    rc = require('./config');


module.exports.getVersion = function(module) {
  return require(module + '/package.json').version;
};


module.exports.newTable = function(options) {
  options = options || {};
  options.chars = options.chars || {
    'top': '', 'top-mid': '', 'top-left': '', 'top-right': '',
    'bottom': '', 'bottom-mid': '', 'bottom-left': '', 'bottom-right': '',
    'left': '', 'left-mid': '', 'mid': '', 'mid-mid': '',
    'right': '', 'right-mid': '', 'middle': ' '
  };
  options.style = options.style || {
    // 'padding-left': 0, 'padding-right': 0
  };
  return new Table(options);
};


module.exports.newApi = function(options) {
  var config = rc;

  options = options || {};
  options.parent = options.parent || {};
  if (options.parent.apiEndpoint) {
    config.api_endpoint = options.parent.apiEndpoint;
  }
  if (options.parent.token) {
    config.token = options.parent.token;
  }
  if (options.parent.dryRun) {
    config.dry_run = options.parent.dryRun;
  }
  return new Api(config);
};


module.exports.truncateUUID = function(input, truncStatus) {
  if (truncStatus || truncStatus == undefined) {
    return input.toString().substring(0, 8);
  } else {
    return input;
  }
};


var error = module.exports.error = function(msg) {
  if (msg && msg.options && msg.options.method && msg.options.url &&
      msg.statusCode && msg.error && msg.error._error) {
    debug('panic', msg);
    console.error('> ' + msg.options.method + ' ' + msg.options.url);
    console.error('< ' + msg.error._error + ' (' + msg.statusCode + ')');
    if (msg.error.fields) {
      _.forEach(msg.error.fields, function(value, key) {
        console.log(' - ' + key + ': ' + value.join('. '));
      });
    }
  } else {
    console.error(msg);
  }
};


var panic = module.exports.panic = function(msg) {
  error(msg);
  console.error('');
  console.error('   Hey ! this is probably a bug !');
  console.error('   Fresh beers will be waiting for you on our next meetup');
  console.error('                          if you report a new issue :) 🍻');
  console.error('');
  console.error('          https://github.com/pathwar/node-pathwar/issues');
  console.error('');
  process.exit(-1);
};
