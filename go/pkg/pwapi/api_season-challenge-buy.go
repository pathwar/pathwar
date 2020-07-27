package pwapi

import (
	"context"

	"github.com/jinzhu/gorm"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) SeasonChallengeBuy(ctx context.Context, in *SeasonChallengeBuy_Input) (*SeasonChallengeBuy_Output, error) {
	in.ApplyDefaults()
	if in == nil || in.FlavorID == "" || in.SeasonID == "" {
		return nil, errcode.ErrMissingInput
	}

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated.Wrap(err)
	}

	flavorID, err := pwdb.GetIDBySlugAndKind(svc.db, in.FlavorID, "challenge-flavor")
	if err != nil {
		return nil, errcode.ErrInvalidFlavor.Wrap(err)
	}
	seasonID, err := pwdb.GetIDBySlugAndKind(svc.db, in.SeasonID, "season")
	if err != nil {
		return nil, errcode.ErrInvalidSeason.Wrap(err)
	}

	// fetch team for this season
	var team pwdb.Team
	{
		err = svc.db.
			Joins("JOIN team_member ON team_member.team_id = team.id AND team_member.user_id = ?", userID).
			Preload("Members").
			Where(pwdb.Team{SeasonID: seasonID}).
			First(&team).
			Error
		if err != nil {
			return nil, errcode.ErrInvalidTeam.Wrap(err)
		}
	}

	// check if season is valid
	var seasonChallenge pwdb.SeasonChallenge
	{
		err = svc.db.
			Preload("Flavor").
			Where(pwdb.SeasonChallenge{FlavorID: flavorID, SeasonID: seasonID}).
			First(&seasonChallenge).
			Error
		if err != nil {
			return nil, errcode.ErrInvalidSeason.Wrap(err)
		}
	}

	// check if challenge and team belongs to the same season
	if seasonChallenge.SeasonID != team.SeasonID {
		return nil, errcode.ErrTeamNotInSeason
	}

	if seasonChallenge.Flavor.PurchasePrice > team.Cash {
		return nil, errcode.ErrNotEnoughCash
	}

	// check for duplicate
	var c int
	err = svc.db.
		Model(pwdb.ChallengeSubscription{}).
		Where(pwdb.ChallengeSubscription{
			SeasonChallengeID: seasonChallenge.ID,
			TeamID:            team.ID,
		}).
		Count(&c).
		Error
	if err != nil {
		return nil, errcode.ErrChallengeAlreadySubscribed.Wrap(err)
	}
	if c > 0 {
		return nil, errcode.ErrChallengeAlreadySubscribed
	}

	// FIXME: validate if team has enough money

	// create subscription
	subscription := pwdb.ChallengeSubscription{
		SeasonChallengeID: seasonChallenge.ID,
		TeamID:            team.ID,
		BuyerID:           userID,
		Status:            pwdb.ChallengeSubscription_Active,
	}
	err = svc.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&subscription).Error; err != nil {
			return err
		}

		// update team cash
		if cash := seasonChallenge.Flavor.PurchasePrice; cash != 0 {
			err = tx.Model(&pwdb.Team{}).
				Where("id = ?", team.ID).
				UpdateColumn("cash", gorm.Expr("cash - ?", cash)).
				Error
			if err != nil {
				return err
			}
		}

		activity := pwdb.Activity{
			Kind:                    pwdb.Activity_SeasonChallengeBuy,
			AuthorID:                userID,
			TeamID:                  team.ID,
			SeasonChallengeID:       seasonChallenge.ID,
			ChallengeSubscriptionID: subscription.ID,
			SeasonID:                seasonChallenge.SeasonID,
		}
		return tx.Create(&activity).Error
	})
	if err != nil {
		return nil, errcode.ErrCreateChallengeSubscription.Wrap(err)
	}

	// load and return the freshly inserted entry
	err = svc.db.
		Preload("Team", "team.deletion_status = ?", pwdb.DeletionStatus_Active).
		Preload("Team.Season").
		Preload("Buyer").
		Preload("SeasonChallenge").
		Preload("SeasonChallenge.Flavor").
		Preload("SeasonChallenge.Flavor.Challenge").
		First(&subscription, subscription.ID).
		Error
	if err != nil {
		return nil, errcode.ErrGetChallengeSubscription.Wrap(err)
	}

	ret := SeasonChallengeBuy_Output{ChallengeSubscription: &subscription}
	return &ret, nil
}

func (in *SeasonChallengeBuy_Input) ApplyDefaults() {
	if in == nil {
		return
	}
	if in.SeasonID == "" {
		in.SeasonID = "solo-mode"
	}
}
