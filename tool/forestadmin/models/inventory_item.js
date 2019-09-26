module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('inventory_item', {
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
    item: {
      type: DataTypes.INTEGER,
    },
    ownerId: {
      type: DataTypes.STRING,
    },
  }, {
    tableName: 'inventory_item',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

