module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('tournament_member', {
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
    tournamentTeamId: {
      type: DataTypes.STRING,
    },
  }, {
    tableName: 'tournament_member',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

