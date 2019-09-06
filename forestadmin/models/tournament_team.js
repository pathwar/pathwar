module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  const Model = sequelize.define('tournament_team', {
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
    isDefault: {
      type: DataTypes.INTEGER,
    },
    tournamentId: {
      type: DataTypes.STRING,
    },
    teamId: {
      type: DataTypes.STRING,
    },
  }, {
    tableName: 'tournament_team',
    underscored: true,
  });

  Model.associate = (models) => {
  };

  return Model;
};

