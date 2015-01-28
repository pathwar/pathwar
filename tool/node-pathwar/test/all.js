"use strict";


var chai = require("chai"),
    debug = require("debug")("tests"),
    Client = require(".."),
    util = require("util");


chai.should();


var valid_token = 'root-token',
    api_endpoint = null;

// if we run using Docker
if (process.env['API_PORT_5000_TCP_ADDR']) {
  api_endpoint = 'http://' + process.env['API_PORT_5000_TCP_ADDR'] + ':' + process.env['API_PORT_5000_TCP_PORT'] + '/';
}

var inspect = function(name, obj) {
  debug(name, util.inspect(obj, {showHidden: false, depth: null}));
};


suite("[client]", function() {
  var client;

  suite('#http-requests', function() {
    setup(function() {
      var options = {
        token: valid_token
      };
      if (api_endpoint) {
        options['api_endpoint'] = api_endpoint;
      }
      client = new Client(options);
    });
    teardown(function() {
      client = null;
    });

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

  suite('#authentication', function() {
    setup(function() {
      var options = {
        token: null
      };
      if (api_endpoint) {
        options['api_endpoint'] = api_endpoint;
      }
      client = new Client(options);
    });
    teardown(function() {
      client = null;
    });

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

    test("should receive a 401 when trying to GET /users", function(done) {
      client.get("/users").then(
        function(res) {
          inspect('res', res);
          done(true);
        },
        function(err) {
          inspect('err', err);
          try {
            (err.statusCode).should.equal(401);
            done();
          } catch (e) {
            done(e);
          }
        });
    });

    test("should authenticate with user and password", function(done) {
      (client.config.token == null).should.true();
      client.login('root', 'toor').then(
        function(res) {
          inspect('res', res);
          try {
            (res.statusCode).should.equal(200);
            (res.body.token == null).should.false();
            (client.config.token == null).should.false();
            (client.config.token).should.equal(res.body.token);
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

    test("should authenticate with user and password then GET /users", function(done) {
      (client.config.token == null).should.true();
      client.login('root', 'toor').then(
        function(res) {
          inspect('res', res);

          client.get('/users').then(
            function(res) {
              try {
                (res.statusCode).should.equal(200);
                done();
              } catch (e) {
                done(e);
              }
            },
            function(err) {
              done(err);
            });

        },
        function(err) {
          inspect('err', err);
          done(err);
        });
    });
  });

  suite('#authenticated', function() {
    setup(function(done) {
      var options = {
        token: null
      };
      if (api_endpoint) {
        options['api_endpoint'] = api_endpoint;
      }
      client = new Client(options);
      return client.login('root', 'password').then(
        function() {
          done();
        });
    });
    teardown(function() {
      client = null;
    });

    suite("#organizations", function() {
      test("should list organizations of the user", function(done) {
        client.organizations_list().then(
          function(res) {
            inspect('res', res);
            try {
              (res.body._ids[0]).should.be.a('string');
              (res.body._ids[0]).should.equal(res.body._items[0]._id);
              (res.body._items[0].organization == null).should.false();
              (res.body._items[0].user).should.equal(client.user_id);
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
    });
  });
});
