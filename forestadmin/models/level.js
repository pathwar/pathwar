module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('level', {
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
    description: {
      type: DataTypes.STRING,
    },
    author: {
      type: DataTypes.STRING,
    },
    locale: {
      type: DataTypes.STRING,
    },
    isDraft: {
      type: DataTypes.INTEGER,
    },
    previewUrl: {
      type: DataTypes.STRING,
    },
  }, {
    tableName: 'level',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

