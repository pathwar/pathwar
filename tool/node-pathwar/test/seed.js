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

    test("should create a user-tokens(is_session=true) as user", function(done) {
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

    suite('#as-admin', function() {

      test("should create some sessions as admin", function(done) {
        var objects = [{
          name: 'Worldwide',
          public: true
        }, {
          name: 'Staff',
          public: false
        }, {
          name: 'Beta',
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

      test("should create some users as admin+user (strange)", function(done) {
        this.timeout(5000);
        var objects = [{
          login: 'joe',
          email: 'joe@pathwar.net',
          active: true,
          password: 'secure'
        }, {
          login: 'm1ch3l',
          email: 'm1ch3l@pathwar.net',
          active: true,
          //available_sessions: [
          //  refs['sessions'][0]['_id'],
          //  refs['sessions'][1]['_id']
          //],
          password: 'super-secure'
        }, {
          login: 'test-moderator',
          email: 'test-moderator@pathwar.net',
          role: 'moderator',
          active: true,
          password: 'super-secure'
        }, {
          login: 'test-admin',
          email: 'test-admin@pathwar.net',
          role: 'admin',
          active: true,
          password: 'super-secure'
        }, {
          login: 'test-user',
          email: 'test-user@pathwar.net',
          role: 'user',
          active: true,
          password: 'super-secure'
        }, {
          login: 'test-level-server',
          email: 'test-level-server@pathwar.net',
          role: 'level-server',
          active: true,
          password: 'super-secure'
        }, {
          login: 'moul',
          email: 'm@42.am',
          active: true,
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

      test("should create some coupons as admin", function(done) {
        var objects = [{
          hash: '1234567890',
          value: 1,
          session: refs.sessions[0]
        }, {
          hash: '000987654321',
          value: 1,
          session: refs.sessions[1]
        }, {
          hash: '000987654321',
          value: 2,
          session: refs.sessions[1]
        }, {
          hash: '000987654321',
          value: 3,
          session: refs.sessions[1]
        }, {
          hash: 'multi-session-trap',
          value: -100
        }, {
          hash: 'multi-session-10',
          value: 10,
          validations_limit: 100
        }];
        for (var i = 0; i < 100; i++) {
          objects.push({
            hash: '10-' + i,
            value: 10,
            session: refs.sessions[0]
          });
        }
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
              refs['coupons'] = ids;
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

      test("should create some servers as admin", function(done) {
        var objects = [{
          name: 'fake-server',
          ip_address: '1.2.3.4',
          active: true,
          token: '1234567890',
          tags: ['fake', 'dummy', 'example']
        }, {
          name: 'dedi-moul',
          ip_address: '195.154.233.249',
          active: true,
          token: '0987654321',
          tags: ['docker', 'x86_64', 'dedibox']
        }];
        client.post("/servers", objects).then(
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
                (item._links.self.title).should.equal('server');
              }
              refs['servers'] = ids;
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

      test("should create some organizations as admin/user (strange)", function(done) {
        var objects = [{
          name: 'pwn-around-the-world',
          session: refs.sessions[0]
        }, {
          name: 'staff',
          session: refs.sessions[1],
          owner: refs.users[0]
        }];
        client.post("/organizations", objects).then(
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
                (item._links.self.title).should.equal('organization');
              }
              refs['organizations'] = ids;
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

      test("should create some achievements as admin", function(done) {
        var objects = [{
          name: 'flash-gordon',
          description: 'Validate a level in less than a minute'
        }, {
          name: 'hack-the-planet',
          description: 'Finish 50 levels'
        }, {
          name: 'API-hacker',
          description: 'Hack the API'
        }];
        client.post("/achievements", objects).then(
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
                (item._links.self.title).should.equal('achievement');
              }
              refs['achievements'] = ids;
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

      test("should create some levels as admin", function(done) {
        var objects = [{
          name: 'welcome',
          description: 'An easy welcome level',
          price: 0,
          tags: ['easy', 'welcome', 'official', 'free'],
          author: 'Pathwar Team'
        }, {
          name: 'pnu',
          description: 'Possible not upload',
          price: 42,
          tags: ['php', 'advanced'],
          author: 'Pathwar Team'
        }, {
          name: 'calc',
          price: 42,
          tags: ['python', 'sql', 'easy'],
          author: 'Pathwar Team'
        }, {
          name: 'upload-hi',
          price: 10,
          tags: ['easy', 'upload', 'php'],
          author: 'Pathwar Team'
        }, {
          name: 'training-http',
          price: 0,
          tags: ['tutorial', 'easy', 'free'],
          author: 'Pathwar Team'
        }, {
          name: 'training-sqli',
          price: 0,
          tags: ['tutorial', 'easy', 'free'],
          author: 'Pathwar Team'
        }, {
          name: 'training-brute',
          price: 0,
          tags: ['tutorial', 'easy', 'free'],
          author: 'Pathwar Team'
        }, {
          name: 'training-include',
          price: 0,
          tags: ['tutorial', 'easy', 'free'],
          author: 'Pathwar Team'
        }, {
          name: 'training-tools',
          price: 0,
          tags: ['tutorial', 'easy', 'free'],
          author: 'Pathwar Team'
        }, {
          name: 'captcha',
          price: 0,
          tags: ['tutorial', 'easy', 'free'],
          author: 'Pathwar Team'
        }];

        for (var i = 1; i < 50; i++) {
          objects.push({
            name: 'fake-level-' + i,
            price: Math.ceil(Math.random() * 30),
            tags: ['fake', 'dummy'],
            author: 'Pathwar Team'
          });
        }
        client.post("/levels", objects).then(
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
                (item._links.self.title).should.equal('level');
              }
              refs['levels'] = ids;
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

      test("should create some items as admin", function(done) {
        var objects = [{
          name: 'spiderpig-glasses',
          description: 'Unlock all level hints',
          price: 4242,
          quantity: 1000
        }, {
          name: 'whoswho shield',
          description: 'Cannot be attacked on whoswho',
          price: 500,
          quantity: 1
        }];
        client.post("/items", objects).then(
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
                (item._links.self.title).should.equal('item');
              }
              refs['items'] = ids;
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

      test("should create some level-hints as admin", function(done) {
        var objects = [{
          level: refs.levels[0],
          name: 'level sources',
          price: 42
        }, {
          level: refs.levels[0],
          name: 'full solution',
          price: 420
        }, {
          level: refs.levels[1],
          name: 'level sources',
          price: 42
        }];
        client.post("/level-hints", objects).then(
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
                (item._links.self.title).should.equal('level hint');
              }
              refs['level-hints'] = ids;
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

      test("should create some level-instances as admin", function(done) {
        var objects = [{
          level: refs.levels[0],
          server: refs.servers[0],
          passphrases: [
            { 'key': '0', 'value': '1234567890' },
            { 'key': '1', 'value': '0987654321' }
          ],
          urls: [
            { 'name': '80', 'url': 'http://1.2.3.4:1234/' },
            { 'name': '22', 'url': 'http://1.2.3.4:1235/' }
          ]
        }];
        for (var i = 0; i < 50; i++) {
          objects.push({
            level: refs.levels[i % refs.servers.length],
            server: refs.servers[i % refs.servers.length],
            urls: [{
              'name': '80',
              'url': 'http://1.2.3.4:' + Math.ceil(Math.random() * 60000 + 2000).toString()
            }],
            passphrases: [{
              'key': '0',
              'value': Math.ceil(Math.random() * 1000).toString()
            }]
          });
        }
        client.post("/level-instances", objects).then(
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
                (item._links.self.title).should.equal('level instance');
              }
              refs['level-instances'] = ids;
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

      // TODO
      // ----
      // POST organization-users as admin
      // POST user-notifications as admin
      // UPDATE organization-statistics as admin

    });

    suite('#as-user', function() {

      test("should create some organization-levels (buy) as user", function(done) {
        var objects = [{
          organization: refs.organizations[0],
          level: refs.levels[0]
        }, {
          organization: refs.organizations[0],
          level: refs.levels[1]
        }, {
          organization: refs.organizations[1],
          level: refs.levels[0]
        }];
        client.post("/organization-levels", objects).then(
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
                (item._links.self.title).should.equal('organization bought level');
              }
              refs['organization-levels'] = ids;
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

      test("should create some organization-achievements as user", function(done) {
        var objects = [{
          organization: refs.organizations[0],
          achievement: refs.achievements[0]
        }, {
          organization: refs.organizations[0],
          achievement: refs.achievements[1]
        }, {
          organization: refs.organizations[1],
          achievement: refs.achievements[0]
        }];
        client.post("/organization-achievements", objects).then(
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
                (item._links.self.title).should.equal('organization earned achievement');
              }
              refs['organization-achievements'] = ids;
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

      // TODO
      // ----
      // POST organization-coupons as user
      // POST organization-items as user
      // POST organization-level-validations as user
      // POST user-organization-invites as user

    });
  });
});
