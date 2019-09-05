module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('team_member', {
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
    role: {
      type: DataTypes.INTEGER,
    },
    userId: {
      type: DataTypes.STRING,
    },
    teamId: {
      type: DataTypes.STRING,
    },
  }, {
    tableName: 'team_member',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

