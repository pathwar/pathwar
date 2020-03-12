package pwapi

import (
	"context"

	"pathwar.land/v2/go/pkg/errcode"
	"pathwar.land/v2/go/pkg/pwdb"
)

func (svc *service) ChallengeSubscriptionValidate(ctx context.Context, in *ChallengeSubscriptionValidate_Input) (*ChallengeSubscriptionValidate_Output, error) {
	// validation
	if in == nil || in.ChallengeSubscriptionID == 0 || in.Passphrase == "" {
		return nil, errcode.ErrMissingInput
	}

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrGetUserIDFromContext.Wrap(err)
	}

	// check input challenge subscription
	// FIXME: or is admin
	var challengeSubscription pwdb.ChallengeSubscription
	err = svc.db.
		Preload("Team", "team.deletion_status = ?", pwdb.DeletionStatus_Active).
		Preload("SeasonChallenge").
		Joins("JOIN team ON team.id = challenge_subscription.team_id").
		Joins("JOIN team_member ON team_member.team_id = team.id AND team_member.user_id = ?", userID).
		First(&challengeSubscription, in.ChallengeSubscriptionID).
		Error
	if err != nil {
		return nil, errcode.ErrGetChallengeSubscription.Wrap(err)
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
	err = svc.db.Create(&validation).Error
	if err != nil {
		return nil, errcode.ErrCreateChallengeValidation.Wrap(err)
	}

	// load and return the freshly inserted entry
	err = svc.db.
		Preload("Author").
		Preload("ChallengeSubscription").
		Preload("ChallengeSubscription.SeasonChallenge").
		Preload("ChallengeSubscription.Validations").
		Preload("ChallengeSubscription.Team").
		First(&validation, validation.ID).
		Error
	if err != nil {
		return nil, errcode.ErrGetChallengeValidation.Wrap(err)
	}

	// FIXME: only redump the validated instance
	for _, instance := range challengeSubscription.SeasonChallenge.Instances {
		err = svc.db.Model(&instance).
			Update(pwdb.ChallengeInstance{
				Status: pwdb.ChallengeInstance_NeedRedump,
			}).
			Error
		if err != nil {
			return nil, errcode.ErrAgentUpdateState.Wrap(err)
		}
	}

	ret := ChallengeSubscriptionValidate_Output{ChallengeValidation: &validation}
	return &ret, nil
}
