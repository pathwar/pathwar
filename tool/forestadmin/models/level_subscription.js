module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('level_subscription', {
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
    tournamentTeamId: {
      type: DataTypes.STRING,
    },
    levelFlavorId: {
      type: DataTypes.STRING,
    },
  }, {
    tableName: 'level_subscription',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

