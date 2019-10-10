module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('achievement', {
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
    type: {
      type: DataTypes.INTEGER,
    },
    isGlobal: {
      type: DataTypes.INTEGER,
    },
    comment: {
      type: DataTypes.STRING,
    },
    argument: {
      type: DataTypes.STRING,
    },
    authorId: {
      type: DataTypes.STRING,
    },
    levelValidationId: {
      type: DataTypes.STRING,
    },
  }, {
    tableName: 'achievement',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

