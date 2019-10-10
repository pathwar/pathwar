const models = require('../../models');
const requireAll = require('require-all');

module.exports = function (app) {
  require('lumber-forestadmin').run(app, {
    modelsDir: __dirname + '/../../models',
    envSecret: process.env.FOREST_ENV_SECRET,
    authSecret: process.env.FOREST_AUTH_SECRET,
    sequelize: models.sequelize, 
  });

  requireAll({
    dirname: __dirname + '/../../routes',
    recursive: true,
    resolve: Module => app.use('/forest', Module)
  });
};
