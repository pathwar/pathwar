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
    client.get("/").then(
      function(res) {
        inspect('res', res);
        assert.equal(200, res.statusCode);
        next();
      },
      function(err) {
        inspect('err', err);
        assert();
      });
  });

  it("should successfully execute GET /levels", function(next) {
    client.get("/levels").then(
      function(res) {
        inspect('res', res);
        assert.equal(200, res.statusCode);
        assert.equal(true, res.body._items[0].price > 0);
        next();
      },
      function(err) {
        inspect('err', err);
        assert();
      });
  });

  it("should successfully execute POST /user-tokens", function(next) {
    client.post("/user-tokens", {
      is_session: true
    }).then(
      function(res) {
        inspect('res', res);
        next();
      },
      function(err) {
        inspect('err', err);
        assert();
      }
    );
  });

  it("should chain a POST /user-tokens and a GET of the created resource using promise chaining", function(next) {
    client.post("/user-tokens", {
      is_session: true
    }).then(
      function(res) {
        inspect('res', res);
        return client.get(res.body._links.self.href);
      },
      function(err) {
        inspect('err', err);
        assert();
      }
    ).then(
      function(res) {
        inspect('res', res);
        next();
      },
      function(err) {
        inspect('err', err);
        assert();
      }
    );
  });

  it("should chain a POST /user-tokens and a GET of the created resource using promise callback", function(next) {
    client.post("/user-tokens", {
      is_session: true
    }).then(
      function(res) {
        inspect('res', res);
        client.get(res.body._links.self.href).then(
          function(res) {
            inspect('res', res);
            next();
          },
          function(err) {
            inspect('err', err);
            assert();
          });
      },
      function(err) {
        inspect('err', err);
        assert();
      });
  });
});
