package pwengine

import (
	"context"
	"fmt"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) SeasonChallengeBuy(ctx context.Context, in *SeasonChallengeBuy_Input) (*SeasonChallengeBuy_Output, error) {
	// validation
	if in == nil || in.SeasonChallengeID == 0 || in.TeamID == 0 {
		return nil, ErrMissingArgument
	}

	userID, err := userIDFromContext(ctx, e.db)
	if err != nil {
		return nil, fmt.Errorf("get userid from context: %w", err)
	}

	// check if user belongs to team
	// FIXME: or is admin
	var team pwdb.Team
	err = e.db.
		Joins("JOIN team_member ON team_member.team_id = team.id AND team_member.user_id = ?", userID).
		Preload("Members").
		First(&team, in.TeamID).
		Error
	if err != nil {
		return nil, ErrInvalidArgument // fmt.Errorf("fetch team: %w", err)
	}

	// check if season is valid
	var seasonChallenge pwdb.SeasonChallenge
	err = e.db.First(&seasonChallenge, in.SeasonChallengeID).Error
	if err != nil {
		return nil, ErrInvalidArgument // fmt.Errorf("fetch season challenge: %w", err)
	}

	// check if challenge and team belongs to the same season
	if seasonChallenge.SeasonID != team.SeasonID {
		return nil, fmt.Errorf("team and challenge should be on the same season")
	}

	// check for duplicate
	var c int
	err = e.db.
		Model(pwdb.ChallengeSubscription{}).
		Where(pwdb.ChallengeSubscription{
			SeasonChallengeID: in.SeasonChallengeID,
			TeamID:            in.TeamID,
		}).
		Count(&c).
		Error
	if err != nil {
		return nil, fmt.Errorf("check for duplicate: %w", err)
	}
	if c > 0 {
		return nil, ErrDuplicate
	}

	// FIXME: validate if team has enough money

	// create subscription
	subscription := pwdb.ChallengeSubscription{
		SeasonChallengeID: in.SeasonChallengeID,
		TeamID:            in.TeamID,
		BuyerID:           userID,
		Status:            pwdb.ChallengeSubscription_Active,
	}
	err = e.db.Create(&subscription).Error
	if err != nil {
		return nil, fmt.Errorf("create challenge subscription: %w", err)
	}

	// load and return the freshly inserted entry
	err = e.db.
		Preload("Team", "team.deletion_status = ?", pwdb.DeletionStatus_Active).
		Preload("Team.Season").
		Preload("Buyer").
		Preload("SeasonChallenge").
		Preload("SeasonChallenge.Flavor").
		Preload("SeasonChallenge.Flavor.Challenge").
		First(&subscription, subscription.ID).
		Error
	if err != nil {
		return nil, fmt.Errorf("fetch challenge subscription: %w", err)
	}

	ret := SeasonChallengeBuy_Output{ChallengeSubscription: &subscription}
	return &ret, nil
}
