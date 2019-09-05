module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('notification', {
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
    isRead: {
      type: DataTypes.INTEGER,
    },
    clickUrl: {
      type: DataTypes.STRING,
    },
    msg: {
      type: DataTypes.STRING,
    },
    args: {
      type: DataTypes.STRING,
    },
    userId: {
      type: DataTypes.STRING,
    },
  }, {
    tableName: 'notification',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

