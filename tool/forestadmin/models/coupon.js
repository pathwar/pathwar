module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('coupon', {
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
    hash: {
      type: DataTypes.STRING,
    },
    value: {
      type: DataTypes.INTEGER,
    },
    maxValidationCount: {
      type: DataTypes.INTEGER,
    },
    tournamentId: {
      type: DataTypes.STRING,
    },
  }, {
    tableName: 'coupon',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

