package pwapi

import (
	"context"
	"time"

	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwdb"
)

func (svc *service) ChallengeSubscriptionClose(ctx context.Context, in *ChallengeSubscriptionClose_Input) (*ChallengeSubscriptionClose_Output, error) {
	// validation
	if in == nil || in.ChallengeSubscriptionID == 0 {
		return nil, errcode.ErrMissingInput
	}

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrGetUserIDFromContext.Wrap(err)
	}

	// fetch subscription
	var subscription pwdb.ChallengeSubscription
	err = svc.db.
		Preload("Team", "team.deletion_status = ?", pwdb.DeletionStatus_Active).
		Preload("SeasonChallenge").
		Preload("Validations").
		Joins("JOIN team ON team.id = challenge_subscription.team_id").
		Joins("JOIN team_member ON team_member.team_id = team.id AND team_member.user_id = ?", userID).
		First(&subscription, in.ChallengeSubscriptionID).
		Error
	if err != nil {
		return nil, errcode.ErrGetChallengeSubscription.Wrap(err)
	}

	if subscription.Status != pwdb.ChallengeSubscription_Active {
		return nil, errcode.ErrChallengeAlreadyClosed
	}

	// check for required validations
	// FIXME: check for required validation, not just for amount of validations
	if len(subscription.Validations) == 0 {
		return nil, errcode.ErrMissingChallengeValidation
	}

	// FIXME: add more validations

	// update challenge subscription
	now := time.Now()
	err = svc.db.
		Model(&subscription).
		Updates(pwdb.ChallengeSubscription{
			Status:   pwdb.ChallengeSubscription_Closed,
			ClosedAt: &now,
			CloserID: userID,
		}).Error
	if err != nil {
		return nil, errcode.ErrUpdateChallengeSubscription.Wrap(err)
	}

	ret := ChallengeSubscriptionClose_Output{ChallengeSubscription: &subscription}
	return &ret, nil
}
