package pwengine

import (
	"context"
	"fmt"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) ChallengeSubscriptionValidate(ctx context.Context, in *ChallengeSubscriptionValidateInput) (*ChallengeSubscriptionValidateOutput, error) {
	// validation
	if in == nil || in.ChallengeSubscriptionID == 0 || in.Passphrase == "" {
		return nil, ErrMissingArgument
	}

	userID, err := userIDFromContext(ctx, e.db)
	if err != nil {
		return nil, fmt.Errorf("get userid from context: %w", err)
	}

	// check input challenge subscription
	// FIXME: or is admin
	var challengeSubscription pwdb.ChallengeSubscription
	err = e.db.
		Preload("Team").
		Preload("SeasonChallenge").
		Joins("JOIN team ON team.id = challenge_subscription.team_id").
		Joins("JOIN team_member ON team_member.team_id = team.id AND team_member.user_id = ?", userID).
		First(&challengeSubscription, in.ChallengeSubscriptionID).
		Error
	if err != nil {
		return nil, ErrInvalidArgument // fmt.Errorf("fetch challenge subscription: %w", err)
	}

	// FIXME: check if passphrase is valid
	// FIXME: check if passphrase_key wasn't already validated for this team ? or let it
	// FIXME: check if challenge subscription is still open

	// create subscription
	validation := pwdb.ChallengeValidation{
		ChallengeSubscriptionID: in.ChallengeSubscriptionID,
		Passphrase:              in.Passphrase,
		PassphraseKey:           "test",
		AuthorID:                userID,
		AuthorComment:           in.Comment,
		Status:                  pwdb.ChallengeValidation_NeedReview,
	}
	err = e.db.Create(&validation).Error
	if err != nil {
		return nil, fmt.Errorf("create challenge validation: %w", err)
	}

	// load and return the freshly inserted entry
	err = e.db.
		Preload("Author").
		Preload("ChallengeSubscription").
		Preload("ChallengeSubscription.SeasonChallenge").
		Preload("ChallengeSubscription.Team").
		First(&validation, validation.ID).
		Error
	if err != nil {
		return nil, fmt.Errorf("fetch challenge validation: %w", err)
	}

	ret := ChallengeSubscriptionValidateOutput{ChallengeValidation: &validation}
	return &ret, nil
}
