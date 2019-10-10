module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('team', {
    id: {
      type: DataTypes.STRING,
      primaryKey: true,
    },
    createdAt: {
      type: DataTypes.DATE,
    },
    updatedAt: {
      type: DataTypes.DATE,
    },
    name: {
      type: DataTypes.STRING,
    },
    gravatarUrl: {
      type: DataTypes.STRING,
    },
    locale: {
      type: DataTypes.STRING,
    },
  }, {
    tableName: 'team',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

