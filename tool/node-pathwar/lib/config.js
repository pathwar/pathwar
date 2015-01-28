var config = require('rc')('pathwar', {
  api_endpoint: 'https://api.pathwar.net/',

  token: null,
  organization: null,
  username: null,
  password: null
});

// FIXME: data validation

module.exports = config;
