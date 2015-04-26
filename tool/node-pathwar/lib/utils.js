var Api = require('./index'),
    Table = require('cli-table'),
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
  if (options.parent.dryRun) {
    config.dry_run = options.parent.dryRun;
  }
  return new Api(config);
};
