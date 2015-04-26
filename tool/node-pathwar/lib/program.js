var _ = require('lodash'),
    program = require('commander'),
    utils = require('./utils');


program
  .version(utils.getVersion('..'))
  .option('--api-endpoint <url>', 'set the API endpoint');
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
  .action(function(type, options) {
    var client = utils.newApi(options);

    type = type || '';
    client.get('/' + type)
      .then(function(res) {
        var keys = _.keys(res.body._items[0]);
        // FIXME: aggregate keys from all items and remove dummy fields
        var table = utils.newTable({
          head: keys
        });
        _.forEach(res.body._items, function(item) {
          var row = [];
          _.forEach(keys, function(key) {
            row.push(item[key] || '');
          });
          table.push(row);
        });
        console.log(table.toString());
      })
      .catch(function(err) {
        console.error(err);
      });
  });


program.command('add');
program.command('rm');
program.command('update');


module.exports = program;


module.exports.run = function() {
  program.parse(process.argv);
  if (!process.argv.slice(2).length) {
    program.outputHelp();
  }
};
