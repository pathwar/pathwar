module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('auth_method', {
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
    identifier: {
      type: DataTypes.STRING,
    },
    emailAddress: {
      type: DataTypes.STRING,
    },
    passwordHash: {
      type: DataTypes.STRING,
    },
    salt: {
      type: DataTypes.STRING,
    },
    totpToken: {
      type: DataTypes.STRING,
    },
    url: {
      type: DataTypes.STRING,
    },
    isVerified: {
      type: DataTypes.INTEGER,
    },
    provider: {
      type: DataTypes.INTEGER,
    },
    userId: {
      type: DataTypes.STRING,
    },
  }, {
    tableName: 'auth_method',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

