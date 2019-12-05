package pwapi

import (
	"context"

	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwdb"
)

func (svc *service) SeasonChallengeBuy(ctx context.Context, in *SeasonChallengeBuy_Input) (*SeasonChallengeBuy_Output, error) {
	if in == nil || in.SeasonChallengeID == 0 || in.TeamID == 0 {
		return nil, errcode.ErrMissingInput
	}

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated.Wrap(err)
	}

	// check if user belongs to team
	// FIXME: or is admin
	var team pwdb.Team
	err = svc.db.
		Joins("JOIN team_member ON team_member.team_id = team.id AND team_member.user_id = ?", userID).
		Preload("Members").
		First(&team, in.TeamID).
		Error
	if err != nil {
		return nil, errcode.ErrInvalidTeam.Wrap(err)
	}

	// check if season is valid
	var seasonChallenge pwdb.SeasonChallenge
	err = svc.db.First(&seasonChallenge, in.SeasonChallengeID).Error
	if err != nil {
		return nil, errcode.ErrInvalidSeason.Wrap(err)
	}

	// check if challenge and team belongs to the same season
	if seasonChallenge.SeasonID != team.SeasonID {
		return nil, errcode.ErrTeamNotInSeason
	}

	// check for duplicate
	var c int
	err = svc.db.
		Model(pwdb.ChallengeSubscription{}).
		Where(pwdb.ChallengeSubscription{
			SeasonChallengeID: in.SeasonChallengeID,
			TeamID:            in.TeamID,
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
		SeasonChallengeID: in.SeasonChallengeID,
		TeamID:            in.TeamID,
		BuyerID:           userID,
		Status:            pwdb.ChallengeSubscription_Active,
	}
	err = svc.db.Create(&subscription).Error
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
