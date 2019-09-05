module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('level_version', {
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
    version: {
      type: DataTypes.STRING,
    },
    changelog: {
      type: DataTypes.STRING,
    },
    isDraft: {
      type: DataTypes.INTEGER,
    },
    isLatest: {
      type: DataTypes.INTEGER,
    },
    sourceUrl: {
      type: DataTypes.STRING,
    },
    driver: {
      type: DataTypes.INTEGER,
    },
    levelId: {
      type: DataTypes.STRING,
    },
  }, {
    tableName: 'level_version',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

