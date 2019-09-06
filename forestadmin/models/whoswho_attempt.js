module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('whoswho_attempt', {
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
    success: {
      type: DataTypes.INTEGER,
    },
    authorId: {
      type: DataTypes.STRING,
    },
    targetMemberId: {
      type: DataTypes.STRING,
    },
    targetTeamId: {
      type: DataTypes.STRING,
    },
  }, {
    tableName: 'whoswho_attempt',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

