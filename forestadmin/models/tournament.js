module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('tournament', {
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
    status: {
      type: DataTypes.INTEGER,
    },
    visibility: {
      type: DataTypes.INTEGER,
    },
    isDefault: {
      type: DataTypes.INTEGER,
    },
  }, {
    tableName: 'tournament',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

