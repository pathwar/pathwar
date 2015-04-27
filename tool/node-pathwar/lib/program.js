var _ = require('lodash'),
    Q = require('q'),
    program = require('commander'),
    utils = require('./utils');


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
  .command('ls [type]')
  .description('list objects')
  //.option('-k, --keys <keys...>', 'keys to print')
  //.option('-q, --quiet', 'only print ids')
  //.option('-f, --filter <filter...>', 'filter constraints')
  .action(function(type, options) {
    var client = utils.newApi(options);

    type = type || '';
    client.get('/' + type)
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

        var table = utils.newTable({
          head: keys
        });

        _.forEach(res.body._items, function(item) {
          var row = [];
          _.forEach(keys, function(key) {
            // FIXME: special print for uuids
            // FIXME: special print for dates
            row.push(item[key] || '');
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
