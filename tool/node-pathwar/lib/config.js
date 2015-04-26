var defaultConfig = {
  api_endpoint: 'https://api.pathwar.net/',

  token: null,
  organization: null,
  username: null,
  password: null
};

var config = require('rc')('pathwar', defaultConfig);

// FIXME: data validation

module.exports = config;
module.exports.defaultConfig = defaultConfig;
