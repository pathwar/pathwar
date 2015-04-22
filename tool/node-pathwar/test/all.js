"use strict";


var _ = require('lodash'),
    chai = require("chai"),
    debug = require("debug")("tests"),
    Client = require(".."),
    util = require("util"),
    should = chai.should();


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
      client.get("/").then(function(res) {
        inspect('res', res);
        try {
          (res.statusCode).should.equal(200);
          done();
        } catch (e) {
          done(e);
        }
      }).catch(function(err) {
        inspect('err', err);
        done(err);
      });
    });

    return;

    test("should successfully execute GET /levels", function(done) {
      client.get("/levels").then(function(res) {
        inspect('res', res);
        try {
          (res.statusCode).should.equal(200);
          res.body._items[0].should.have.property('price');
          done();
        } catch (e) {
          done(e);
        }
      }).catch(function(err) {
        inspect('err', err);
        done(err);
      });
    });

    test("should successfully execute POST /user-tokens", function(done) {
      client.post("/user-tokens", {
        is_session: true
      }).then(function(res) {
        inspect('res', res);
        done();
      }).catch(function(err) {
        inspect('err', err);
        done(err);
      });
    });

    test("should chain a POST /user-tokens and a GET of the created resource using promise chaining", function(done) {
      client.post("/user-tokens", {
        is_session: true
      }).then(function(res) {
        inspect('res', res);
        return client.get(res.body._links.self.href);
      }).then(function(res) {
        inspect('res', res);
        done();
      }).catch(function(err) {
        inspect('err', err);
        done(err);
      });
    });

    test("should chain a POST /user-tokens and a GET of the created resource using promise callback", function(done) {
      client.post("/user-tokens", {
        is_session: true
      }).then(function(res) {
        inspect('res', res);
        client.get(res.body._links.self.href).then(function(res) {
          inspect('res', res);
          done();
        }).catch(function(err) {
          inspect('err', err);
          done(err);
        });
      }).catch(function(err) {
        inspect('err', err);
        done(err);
      });
    });

    test("should trigger the error callback on 404", function(done) {
      client.get("/do-not-exists").then(function(res) {
        inspect('res', res);
        done(true);
      }).catch(function(err) {
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
      client.get("/").then(function(res) {
        inspect('res', res);
        try {
          (res.statusCode).should.equal(200);
          done();
        } catch (e) {
          done(e);
        }
      }).catch(function(err) {
        inspect('err', err);
        done(err);
      });
    });

    test("should receive a 401 when trying to GET /users", function(done) {
      client.get("/users").then(function(res) {
        inspect('res', res);
        done(true);
      }).catch(function(err) {
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
      should.not.exist(client.config.token);

      client.login('root', 'toor').then(function(res) {
        inspect('res', res);
        try {
          (res.statusCode).should.equal(200);
          (res.body.token == null).should.false();
          should.exist(client.config.token);
          (client.config.token).should.equal(res.body.token);
          done();
        } catch (e) {
          done(e);
        }
      }).catch(function(err) {
        inspect('err', err);
        done(err);
      });
    });

    test("should authenticate with user and password then GET /users", function(done) {
      should.not.exist(client.config.token);

      client.login('root', 'toor').then(function(res) {
        inspect('res', res);

        client.get('/users').then(function(res) {
          try {
            (res.statusCode).should.equal(200);
            done();
          } catch (e) {
            done(e);
          }
        }).catch(function(err) {
          done(err);
        });
      }).catch(function(err) {
        inspect('err', err);
        done(err);
      });
    });
  });

  suite('#authenticated', function() {
    var first_organization;

    setup(function(done) {
      var options = {
        token: null
      };
      if (api_endpoint) {
        options['api_endpoint'] = api_endpoint;
      }
      client = new Client(options);
      return client.login('root', 'password').then(function() {
        done();
      });
    });
    teardown(function() {
      client = null;
    });

    suite("#organizations-manipulation", function() {
      test("should list organizations of the user", function(done) {
        client.organizations_list().then(function(res) {
          inspect('res', res);
          try {
            (res.body._ids[0]).should.be.a('string');
            (res.body._ids[0]).should.equal(res.body._items[0]._id);
            should.exist(res.body._items[0].organization);
            (res.body._items[0].user).should.equal(client.user_id);
            (res.statusCode).should.equal(200);
            first_organization = res.body._items[0]['organization'];
            done();
          } catch (e) {
            done(e);
          }
        }).catch(function(err) {
          inspect('err', err);
          done(err);
        });
      });

      test("should switch to the first organization of the user", function(done) {
        should.not.exist(client.organization_id);
        client.organizations_select(first_organization).then(function(res) {
          inspect('res', res);
          try {
            should.exist(client.organization_id);
            (client.organization_id).should.equal(res.body._id);
            (res.statusCode).should.equal(200);
            done();
          } catch (e) {
            done(e);
          }
        }).catch(function(err) {
          inspect('err', err);
          done(err);
        });
      });
    });

    suite("#organizations-context", function() {
      setup(function(done) {
        client.organizations_select(first_organization).then(function(res) {
          done();
        });
      });

      test("should fetch members of the current organization", function(done) {
        client.my_organization_users().then(function(res) {
          try {
            (res.statusCode).should.equal(200);
            _.findIndex(res.body._items, {'user': client.user_id}).should.not.equal(-1);
            done();
          } catch (e) {
            done(e);
          }
        }).catch(function(err) {
          done(err);
        });
      });

      test("should fetch bought levels of the current organization", function(done) {
        client.my_organization_levels().then(function(res) {
          try {
            (res.statusCode).should.equal(200);
            res.body._items.should.be.an('array');
            done();
          } catch (e) {
            done(e);
          }
        }).catch(function(err) {
          done(err);
        });
      });
    });
  });
});
