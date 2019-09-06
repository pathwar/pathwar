module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('level_flavor', {
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
    levelVersionId: {
      type: DataTypes.STRING,
    },
  }, {
    tableName: 'level_flavor',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

