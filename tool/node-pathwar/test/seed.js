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


suite("[seed]", function() {
  var client;

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

  suite('#checks', function() {
    test("should have an empty database", function(done) {
      client.get("/levels").then(
        function(res) {
          inspect('res', res);
          try {
            (res.body._meta.total).should.equal(0);
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
    test("should successfully create a user-tokens(is_session=true)", function(done) {
      client.post("/user-tokens", {
        is_session: true
      }).then(
        function(res) {
          inspect('res', res);
          try {
            (res.statusCode).should.equal(201);
            (res.body._status).should.equal('OK');
            (res.body._links.self.title).should.equal('user token');
            done();
          } catch (e) {
            done(e);
          }
        },
        function(err) {
          inspect('err', err);
          done(err);
        }
      );
    });
  });

  suite('#seed', function() {
    var refs = {};
    test("should successfully create some sessions", function(done) {
      var objects = [{
        name: 'world',
        public: true
      }, {
        name: 'beta',
        public: false
      }, {
        name: 'Epitech2015',
        public: false
      }];
      client.post("/sessions", objects).then(
        function(res) {
          inspect('res', res);
          try {
            (res.statusCode).should.equal(201);
            (res.body._status).should.equal('OK');
            (res.body._items.length).should.equal(objects.length);
            var ids = [];
            for (var idx in res.body._items) {
              var item = res.body._items[idx];
              ids.push(item._id);
              (item._status).should.equal('OK');
              (item._links.self.title).should.equal('session');
            }
            refs['sessions'] = ids;
            done();
          } catch (e) {
            done(e);
          }
        },
        function(err) {
          inspect('err', err);
          done(err);
        }
      );
    });
    test("should successfully create some users", function(done) {
      var objects = [{
        login: 'joe',
        email: 'joe@pathwar.net',
        password: 'secure'
      }, {
        login: 'm1ch3l',
        email: 'm1ch3l@pathwar.net',
        role: 'superuser',
        active: true,
        //available_sessions: [
        //  refs['sessions'][0]['_id'],
        //  refs['sessions'][1]['_id']
        //],
        password: 'super-secure'
      }];
      client.post("/users", objects).then(
        function(res) {
          inspect('res', res);
          try {
            (res.statusCode).should.equal(201);
            (res.body._status).should.equal('OK');
            (res.body._items.length).should.equal(objects.length);
            var ids = [];
            for (var idx in res.body._items) {
              var item = res.body._items[idx];
              ids.push(item._id);
              (item._status).should.equal('OK');
              (item._links.self.title).should.equal('user');
            }
            refs['users'] = ids;
            done();
          } catch (e) {
            done(e);
          }
        },
        function(err) {
          inspect('err', err);
          done(err);
        }
      );
    });
    test("should successfully create some coupons", function(done) {
      var objects = [{
        hash: '1234567890',
        value: 42,
        session: refs.sessions[0]
      }, {
        hash: '000987654321',
        value: 24,
        session: refs.sessions[1]
      }];
      client.post("/coupons", objects).then(
        function(res) {
          inspect('res', res);
          try {
            (res.statusCode).should.equal(201);
            (res.body._status).should.equal('OK');
            (res.body._items.length).should.equal(objects.length);
            var ids = [];
            for (var idx in res.body._items) {
              var item = res.body._items[idx];
              ids.push(item._id);
              (item._status).should.equal('OK');
              (item._links.self.title).should.equal('coupon');
            }
            refs['users'] = ids;
            done();
          } catch (e) {
            done(e);
          }
        },
        function(err) {
          inspect('err', err);
          done(err);
        }
      );
    });
  });
});
