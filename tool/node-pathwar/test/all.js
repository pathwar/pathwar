"use strict";

var assert = require("assert"),
    debug = require("debug")("tests"),
    Client = require(".."),
    util = require('util');

var inspect = function(name, obj) {
  debug(name, util.inspect(obj, {showHidden: false, depth: null}));
};

describe("[client]", function() {
  var client;

  beforeEach(function() {
    client = new Client();
  });

  it("should successfully execute GET /", function(next) {
    client.get("/")
      .then(function(res) {
        inspect('res', res);
        assert.equal(200, res.statusCode);
        next();
      }, function(err) {
        inspect('err', err);
        assert();
      });
  });

  it("should successfully execute GET /levels", function(next) {
    client.get("/levels")
      .then(function(res) {
        inspect('res', res);
        assert.equal(200, res.statusCode);
        assert.equal(true, res.body._items[0].price > 0);
        next();
      }, function(err) {
        inspect('err', err);
        assert();
      });
  });

});
