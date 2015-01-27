"use strict";


var chai = require("chai"),
    debug = require("debug")("tests"),
    Client = require(".."),
    util = require("util");


chai.should();


var inspect = function(name, obj) {
  debug(name, util.inspect(obj, {showHidden: false, depth: null}));
};


suite("[client]", function() {
  var client;

  setup(function() {
    client = new Client();
  });

  suite('#http-requests', function() {
    test("should successfully execute GET /", function(done) {
      client.get("/").then(
        function(res) {
          inspect('res', res);
          try {
            (res.statusCode).should.equal(200);
            done();
          } catch (e) {
            done(e);
          }
        },
        function(err) {
          inspect('err', err);
          done(err);
        });
    });

    test("should successfully execute GET /levels", function(done) {
      client.get("/levels").then(
        function(res) {
          inspect('res', res);
          try {
            (res.statusCode).should.equal(200);
            res.body._items[0].should.have.property('price');
            done();
          } catch (e) {
            done(e);
          }
        },
        function(err) {
          inspect('err', err);
          done(err);
        });
    });

    test("should successfully execute POST /user-tokens", function(done) {
      client.post("/user-tokens", {
        is_session: true
      }).then(
        function(res) {
          inspect('res', res);
          done();
        },
        function(err) {
          inspect('err', err);
          done(err);
        }
      );
    });

    test("should chain a POST /user-tokens and a GET of the created resource using promise chaining", function(done) {
      client.post("/user-tokens", {
        is_session: true
      }).then(
        function(res) {
          inspect('res', res);
          return client.get(res.body._links.self.href);
        },
        function(err) {
          inspect('err', err);
          done(err);
        }
      ).then(
        function(res) {
          inspect('res', res);
          done();
        },
        function(err) {
          inspect('err', err);
          done(err);
        }
      );
    });

    test("should chain a POST /user-tokens and a GET of the created resource using promise callback", function(done) {
      client.post("/user-tokens", {
        is_session: true
      }).then(
        function(res) {
          inspect('res', res);
          client.get(res.body._links.self.href).then(
            function(res) {
              inspect('res', res);
              done();
            },
            function(err) {
              inspect('err', err);
              done(err);
            });
        },
        function(err) {
          inspect('err', err);
          done(err);
        });
    });

    test("should trigger the error callback on 404", function(done) {
      client.get("/do-not-exists").then(
        function(res) {
          inspect('res', res);
          done(true);
        },
        function(err) {
          inspect('err', err);
          done();
        });
      });

  });
});
