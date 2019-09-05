module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('user', {
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
    username: {
      type: DataTypes.STRING,
    },
    gravatarUrl: {
      type: DataTypes.STRING,
    },
    websiteUrl: {
      type: DataTypes.STRING,
    },
    locale: {
      type: DataTypes.STRING,
    },
    isStaff: {
      type: DataTypes.INTEGER,
    },
  }, {
    tableName: 'user',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

