var _ = require('lodash'),
    Q = require('q'),
    moment = require('moment'),
    program = require('commander'),
    utils = require('./utils'),
    validator = require('validator');


program
  .version(utils.getVersion('..'))
  .option('--api-endpoint <url>', 'set the API endpoint')
  .option('--token <token>', 'set the token');
//.option('--dry-run', 'do not execute actions')
//.option('-D, --debug', 'enable debug mode')


program._events.version = null;
program
  .command('version')
  .description('show the version information')
  .action(function() {
    console.log('Client version: ' + utils.getVersion('..'));
    console.log('Node.js version (client): ' + process.version);
    console.log('OS/Arch (client): ' + process.platform + '/' + process.arch);
    // FIXME: add information about server
  });


program
  .command('ls <type> [conditions...]')
  .description('list objects')
  .option('--no-trunc', "don't truncate output")
  .option('-f, --field <field>', 'fields to print', utils.collect, [])
  //.option('-q, --quiet', 'only print ids')
  .action(function(type, conditions, options) {
    var client = utils.newApi(options);

    type = type || '';
    var url = '/' + type;
    if (conditions.length) {
      var where = {};
      _.forEach(conditions, function(condition) {
        var split = condition.split('=');
        var key = split[0], value = split[1];

        if (['true', 'false', '1', '0', 'True', 'False'].indexOf(value) >= 0) {
          value = validator.toBoolean(value.toLowerCase());
        }

        if (validator.isNumeric(value)) {
          value = parseInt(value);
        }

        where[key] = value;
        // FIXME: cast values accordingly to the resources

      });
      url = url + '?where=' + JSON.stringify(where);
    }
    console.log('url', url);
    client.get(url)
      .then(function(res) {
        if (!res.body._items.length) {
          console.error('No entries');
          return;
        }

        // get all keys
        var keys = _.difference(
          _.union(_.flatten(_.map(res.body._items, _.keys))),
          ['_links', '_etag']
        );
        if (options.field.length) {
          keys = _.intersection(keys, options.field);
          var difference = _.difference(options.field, keys);
          if (difference.length) {
            utils.error('Unknown fields: ' + difference);
          }
        }
        if (!keys.length) {
          utils.panic('No fields to print');
        }

        var table = utils.newTable({
          head: keys
        });

        _.forEach(res.body._items, function(item) {
          var row = [];
          _.forEach(keys, function(key) {

            switch (key) {

              // Dates
            case '_updated':
            case '_created':
              row.push(moment(item[key]).fromNow());
              break;

              // UUID
            case '_id':
              row.push(utils.truncateUUID(item[key], options.trunc));
              break;

            default:
              var value = (item[key] || '').toString();

              if (validator.isUUID(value)) {
                value = utils.truncateUUID(value, options.trunc);
              }

              if (value.substring(0, 4) == '$2a$') {
                value = '<blowfish>';
              }

              row.push(value);
              break;
            }

          });
          table.push(row);
        });

        console.log(table.toString());
      })
      .catch(utils.panic);
  });


program
  .command('cat <item>')
  .description('show object')
  .action(function(item, options) {
    var client = utils.newApi(options);

    var once = function(item) {
      // FIXME: if object is resolved, only create a request for his type
      return [
        client.get('/servers/' + item._id),
        client.get('/users/' + item._id),
        client.get('/coupons/' + item._id)
        // FIXME: add all kind of items
      ];
    };

    var promises = once({_id: item});
    Q.allSettled(promises).then(function(results) {
      var items = _.compact(_.pluck(_.pluck(results, 'value'), 'body'));
      // FIXME: handle --format option
      console.log(items[0]);
    }, utils.panic);
  });


program.command('touch');
program.command('rm');
program.command('ed');


module.exports = program;


module.exports.run = function() {
  program.parse(process.argv);
  if (!process.argv.slice(2).length) {
    program.outputHelp();
  }
};
