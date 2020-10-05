package pwapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
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
	var subscription pwdb.ChallengeSubscription
	err = svc.db.
		Preload("Team", "team.deletion_status = ?", pwdb.DeletionStatus_Active).
		Preload("SeasonChallenge").
		Preload("SeasonChallenge.Flavor").
		Preload("SeasonChallenge.Flavor.Instances").
		Preload("SeasonChallenge.Season").
		Joins("JOIN team ON team.id = challenge_subscription.team_id").
		Joins("JOIN team_member ON team_member.team_id = team.id AND team_member.user_id = ?", userID).
		First(&subscription, in.ChallengeSubscriptionID).
		Error
	if err != nil {
		return nil, errcode.ErrGetChallengeSubscription.Wrap(err)
	}

	// check if challengesubscription subscription is still open
	if subscription.Status != pwdb.ChallengeSubscription_Active {
		return nil, errcode.ErrChallengeInactiveValidation.Wrap(errors.New("challenge is disabled"))
	}

	instances := subscription.SeasonChallenge.Flavor.Instances
	if len(instances) == 0 {
		return nil, errcode.ErrChallengeInactiveValidation.Wrap(errors.New("challenge has no instances"))
	}

	// compare input and instances' passphrases
	var amountExpected int
	{
		configData, err := instances[0].ParseInstanceConfig()
		if err != nil {
			return nil, err
		}
		amountExpected = len(configData.Passphrases)
		if amountExpected == 0 {
			return nil, errcode.ErrChallengeInactiveValidation.Wrap(errors.New("challenge config is invalid"))
		}
		if amountExpected < len(in.Passphrases) {
			return nil, errcode.ErrChallengeIncompleteValidation.Wrap(fmt.Errorf("too many passphrases"))
		}
	}

	// FIXME: revalidation

	validPassphrases := make([]bool, amountExpected)
	usedInstances := make(map[int64]bool, len(instances))
	for _, instance := range instances {
		configData, err := instance.ParseInstanceConfig()
		if err != nil {
			return nil, err
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

	if subscription.SeasonChallenge.Season.IsGlobal {
		validation.Status = pwdb.ChallengeValidation_AutoAccepted
	}

	// update DB
	err = svc.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&validation).Error
		if err != nil {
			return errcode.ErrCreateChallengeValidation.Wrap(err)
		}

		// update challenge subscription
		now := time.Now()
		err = tx.
			Model(&subscription).
			Updates(pwdb.ChallengeSubscription{
				Status:   pwdb.ChallengeSubscription_Closed,
				ClosedAt: &now,
				CloserID: userID,
			}).Error
		if err != nil {
			return errcode.ErrUpdateChallengeSubscription.Wrap(err)
		}

		// mark used instances as needing a redump
		usedInstanceIDs := make([]int64, len(usedInstances))
		i := 0
		for id := range usedInstances {
			usedInstanceIDs[i] = id
			i++
		}
		err = tx.
			Model(&instances[0]).
			Where("id IN (?)", usedInstanceIDs).
			Update(pwdb.ChallengeInstance{Status: pwdb.ChallengeInstance_NeedRedump, InstanceConfig: []byte{}}).
			Error
		if err != nil {
			return errcode.ErrAgentUpdateState.Wrap(err)
		}

		// update team cash
		err = tx.Model(&pwdb.Team{}).
			Where("id = ?", subscription.TeamID).
			UpdateColumn("cash", gorm.Expr("cash + ?", subscription.SeasonChallenge.Flavor.ValidationReward)).
			Error
		if err != nil {
			return err
		}

		activity := pwdb.Activity{
			Kind:                    pwdb.Activity_ChallengeSubscriptionValidate,
			AuthorID:                userID,
			ChallengeSubscriptionID: subscription.ID,
			TeamID:                  subscription.TeamID,
			SeasonChallengeID:       subscription.SeasonChallengeID,
			ChallengeFlavorID:       subscription.SeasonChallenge.FlavorID,
			SeasonID:                subscription.SeasonChallenge.SeasonID,
		}
		return tx.Create(&activity).Error
	})
	if err != nil {
		return nil, err
	}

	// load updated challenge subscription with validations
	err = svc.db.
		Preload("Validations").
		Where(pwdb.ChallengeSubscription{ID: subscription.ID}).
		First(&subscription).
		Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	// load freshly inserted entry
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
	ret := ChallengeSubscriptionValidate_Output{
		ChallengeValidation: &validation,
	}
	return &ret, nil
}
