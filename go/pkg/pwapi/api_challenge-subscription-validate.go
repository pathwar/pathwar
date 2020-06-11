package pwapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
	"pathwar.land/pathwar/v2/go/pkg/pwinit"
)

func (svc *service) ChallengeSubscriptionValidate(ctx context.Context, in *ChallengeSubscriptionValidate_Input) (*ChallengeSubscriptionValidate_Output, error) {
	// validation
	if in == nil || in.ChallengeSubscriptionID == 0 || in.Passphrases == nil {
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
		Preload("SeasonChallenge.Flavor").
		Preload("SeasonChallenge.Flavor.Instances").
		Joins("JOIN team ON team.id = challenge_subscription.team_id").
		Joins("JOIN team_member ON team_member.team_id = team.id AND team_member.user_id = ?", userID).
		First(&challengeSubscription, in.ChallengeSubscriptionID).
		Error
	if err != nil {
		return nil, errcode.ErrGetChallengeSubscription.Wrap(err)
	}

	// check if challenge subscription is still open
	if challengeSubscription.Status != pwdb.ChallengeSubscription_Active {
		return nil, errcode.ErrChallengeInactiveValidation.Wrap(errors.New("challenge is disabled"))
	}

	instances := challengeSubscription.SeasonChallenge.Flavor.Instances
	if len(instances) == 0 {
		return nil, errcode.ErrChallengeInactiveValidation.Wrap(errors.New("challenge has no instances"))
	}

	// compare input and instances' passphrases
	var configData pwinit.InitConfig
	err = json.Unmarshal(instances[0].GetInstanceConfig(), &configData)
	if err != nil {
		return nil, errcode.ErrParseInitConfig.Wrap(err)
	}
	amountExpected := len(configData.Passphrases)
	if amountExpected == 0 {
		return nil, errcode.ErrChallengeInactiveValidation.Wrap(errors.New("challenge config is invalid"))
	}
	if amountExpected < len(in.Passphrases) {
		return nil, errcode.ErrChallengeIncompleteValidation.Wrap(fmt.Errorf("too many passphrases"))
	}

	// FIXME: revalidation

	validPassphrases := make([]bool, amountExpected)
	usedInstances := make(map[int64]bool, len(instances))
	for _, instance := range instances {
		var configData pwinit.InitConfig
		err = json.Unmarshal(instance.GetInstanceConfig(), &configData)
		if err != nil {
			return nil, errcode.ErrParseInitConfig.Wrap(err)
		}
		for index, passphrase := range configData.Passphrases {
			for _, userPassphrase := range in.Passphrases {
				if passphrase == userPassphrase {
					validPassphrases[index] = true
					usedInstances[instance.ID] = true
				}
			}
		}
	}
	amountValid := 0
	for _, valid := range validPassphrases {
		if valid {
			amountValid++
		}
	}

	if amountValid > 0 && amountValid < amountExpected {
		return nil, errcode.ErrChallengeIncompleteValidation.Wrap(fmt.Errorf("%d/%d valid passphrases", amountValid, amountExpected))
	}

	if amountValid == 0 {
		return nil, errcode.ErrChallengeIncompleteValidation.Wrap(errors.New("invalid passphrase(s)"))
	}

	// create validation
	validPassphraseIndices := []int{}
	for index, valid := range validPassphrases {
		if valid {
			validPassphraseIndices = append(validPassphraseIndices, index)
		}
	}
	passphrases, err := json.Marshal(validPassphraseIndices)
	if err != nil {
		return nil, errcode.ErrChallengeJSONMarshalPassphrases.Wrap(err)
	}
	validation := pwdb.ChallengeValidation{
		ChallengeSubscriptionID: in.ChallengeSubscriptionID,
		Passphrases:             string(passphrases),
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

	// mark used instances as needing a redump
	usedInstanceIDs := make([]int64, len(usedInstances))
	i := 0
	for id := range usedInstances {
		usedInstanceIDs[i] = id
		i++
	}
	err = svc.db.
		Model(&instances[0]).
		Where("id IN (?)", usedInstanceIDs).
		Update(pwdb.ChallengeInstance{Status: pwdb.ChallengeInstance_NeedRedump}).
		Error
	if err != nil {
		return nil, errcode.ErrAgentUpdateState.Wrap(err)
	}

	ret := ChallengeSubscriptionValidate_Output{
		ChallengeValidation: &validation,
	}
	return &ret, nil
}
