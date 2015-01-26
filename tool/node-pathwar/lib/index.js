var axios = require('axios');


var Client = module.exports = function(config) {
  this.config = config;
};


(function() {
  this.get = function(url, config) {
    return axios.get(url, config);
  };
}).call(Client.prototype);
