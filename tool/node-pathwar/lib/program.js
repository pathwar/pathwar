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
  .alias('select')
  .description('list objects')
  .option('--no-trunc', "don't truncate output")
  .option('-f, --field <field>', 'fields to print', utils.collect, [])
  //.option('-q, --quiet', 'only print ids')
  .action(function(type, conditions, options) {
    var client = utils.newApi(options);

    type = type || '';
    var url = '/' + type;
    if (conditions.length) {
      var where = utils.castFields(type, conditions);
      url = url + '?where=' + JSON.stringify(where);
    }
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
  .alias('show')
  .description('show object')
  .action(function(item, options) {
    var client = utils.newApi(options);

    utils.searchItems(item, client, function(items) {
      // FIXME: handle --format option
      console.log(items[0]);
    }, utils.panic);
  });


program
  .command('rm <item>')
  .alias('delete')
  .description('remove an item')
  .action(function(item, options) {
    var client = utils.newApi(options);

    // FIXME: add warning !

    utils.searchItems(item, client, function(items) {
      client.delete(items[0]).then(function(res) {
        console.log('done');
      }).catch(utils.panic);
    });
  });


program
  .command('touch <type> [fields...]')
  .alias('add')
  .description('create an item')
  .action(function(type, fields, options) {
    var client = utils.newApi(options);

    var input = utils.castFields(type, fields);

    client.post('/' + type, input).then(function(res) {
      console.log(res.body._id);
    }).catch(utils.panic);
  });


program.command('ed');


module.exports = program;


module.exports.run = function() {
  program.parse(process.argv);
  if (!process.argv.slice(2).length) {
    program.outputHelp();
  }
};
