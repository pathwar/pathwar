module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('level_validation', {
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
    status: {
      type: DataTypes.INTEGER,
    },
    authorComment: {
      type: DataTypes.STRING,
    },
    correctorComment: {
      type: DataTypes.STRING,
    },
    levelSubscriptionId: {
      type: DataTypes.STRING,
    },
  }, {
    tableName: 'level_validation',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

