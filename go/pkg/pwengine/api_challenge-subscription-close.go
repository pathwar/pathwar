package pwengine

import (
	"context"
	"fmt"
	"time"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) ChallengeSubscriptionClose(ctx context.Context, in *ChallengeSubscriptionCloseInput) (*ChallengeSubscriptionCloseOutput, error) {
	// validation
	if in == nil || in.ChallengeSubscriptionID == 0 {
		return nil, ErrMissingArgument
	}

	userID, err := userIDFromContext(ctx, e.db)
	if err != nil {
		return nil, fmt.Errorf("get userid from context: %w", err)
	}

	// fetch subscription
	var subscription pwdb.ChallengeSubscription
	err = e.db.
		Preload("Team").
		Preload("SeasonChallenge").
		Preload("Validations").
		Joins("JOIN team ON team.id = challenge_subscription.team_id").
		Joins("JOIN team_member ON team_member.team_id = team.id AND team_member.user_id = ?", userID).
		Where(pwdb.ChallengeSubscription{
			Status: pwdb.ChallengeSubscription_Active,
		}).
		First(&subscription, in.ChallengeSubscriptionID).
		Error
	if err != nil {
		return nil, ErrInvalidArgument // fmt.Errorf("fetch challenge subscription: %w", err)
	}

	// check for required validations
	// FIXME: check for required validation, not just for amount of validations
	if len(subscription.Validations) == 0 {
		return nil, ErrMissingRequiredValidation
	}

	// FIXME: add more validations

	// update challenge subscription
	now := time.Now()
	err = e.db.
		Model(&subscription).
		Updates(pwdb.ChallengeSubscription{
			Status:   pwdb.ChallengeSubscription_Closed,
			ClosedAt: &now,
			CloserID: userID,
		}).Error
	if err != nil {
		return nil, fmt.Errorf("update challenge subscription: %w", err)
	}

	ret := ChallengeSubscriptionCloseOutput{ChallengeSubscription: &subscription}
	return &ret, nil
}
