// This model was generated by Lumber. However, you remain in control of your models.
// Learn how here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/models/enrich-your-models
module.exports = (sequelize, DataTypes) => {
  const { Sequelize } = sequelize;
  // This section contains the fields of your model, mapped to your table's columns.
  // Learn more here: https://docs.forestadmin.com/documentation/v/v6/reference-guide/models/enrich-your-models#declaring-a-new-field-in-a-model
  const User = sequelize.define('user', {
    createdAt: {
      type: DataTypes.DATE,
    },
    updatedAt: {
      type: DataTypes.DATE,
    },
    deletedAt: {
      type: DataTypes.DATE,
    },
    username: {
      type: DataTypes.STRING,
    },
    email: {
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
    oAuthSubject: {
      type: DataTypes.STRING,
    },
    deletionReason: {
      type: DataTypes.STRING,
    },
    deletionStatus: {
      type: DataTypes.INTEGER,
    },
    activeTeamMemberId: {
      type: DataTypes.BIGINT,
    },
    activeSeasonId: {
      type: DataTypes.BIGINT,
    },
  }, {
    tableName: 'user',
    underscored: true,
  });

  // This section contains the relationships for this model. See: https://docs.forestadmin.com/documentation/v/v6/reference-guide/relationships#adding-relationships.
  User.associate = (models) => {
    User.hasMany(models.couponValidation, {
      foreignKey: {
        name: 'authorIdKey',
        field: 'author_id',
      },
      as: 'authorCouponValidations',
    });
    User.hasMany(models.notification, {
      foreignKey: {
        name: 'userIdKey',
        field: 'user_id',
      },
      as: 'notifications',
    });
    User.hasMany(models.organizationMember, {
      foreignKey: {
        name: 'userIdKey',
        field: 'user_id',
      },
      as: 'organizationMembers',
    });
    User.hasMany(models.challengeSubscription, {
      foreignKey: {
        name: 'buyerIdKey',
        field: 'buyer_id',
      },
      as: 'buyerChallengeSubscriptions',
    });
    User.hasMany(models.challengeValidation, {
      foreignKey: {
        name: 'authorIdKey',
        field: 'author_id',
      },
      as: 'authorChallengeValidations',
    });
    User.hasMany(models.whoswhoAttempt, {
      foreignKey: {
        name: 'targetUserIdKey',
        field: 'target_user_id',
      },
      as: 'targetUserWhoswhoAttempts',
    });
    User.hasMany(models.whoswhoAttempt, {
      foreignKey: {
        name: 'authorIdKey',
        field: 'author_id',
      },
      as: 'authorWhoswhoAttempts',
    });
    User.hasMany(models.achievement, {
      foreignKey: {
        name: 'authorIdKey',
        field: 'author_id',
      },
      as: 'authorAchievements',
    });
    User.hasMany(models.teamMember, {
      foreignKey: {
        name: 'userIdKey',
        field: 'user_id',
      },
      as: 'teamMembers',
    });
  };

  return User;
};