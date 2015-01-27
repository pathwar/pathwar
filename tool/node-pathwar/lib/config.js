var config = require('rc')('pathwar', {
  api_endpoint: 'https://api.pathwar.net/',

  token: null
});

// FIXME: data validation

module.exports = config;
