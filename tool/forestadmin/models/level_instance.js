module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('level_instance', {
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
    hypervisorId: {
      type: DataTypes.STRING,
    },
    flavorId: {
      type: DataTypes.STRING,
    },
  }, {
    tableName: 'level_instance',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

