module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('coupon_validation', {
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
    comment: {
      type: DataTypes.STRING,
    },
    authorId: {
      type: DataTypes.STRING,
    },
    couponId: {
      type: DataTypes.STRING,
    },
  }, {
    tableName: 'coupon_validation',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

