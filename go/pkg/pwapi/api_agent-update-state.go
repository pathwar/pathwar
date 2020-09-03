package pwapi

import (
	"context"
	"reflect"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AgentUpdateState(ctx context.Context, in *AgentUpdateState_Input) (*AgentUpdateState_Output, error) {
	if !isAgentContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrGetUserIDFromContext.Wrap(err)
	}

	var dbInstances []*pwdb.ChallengeInstance
	err = svc.db.Find(&dbInstances).Error
	if err != nil {
		return nil, errcode.ErrAgentUpdateState.Wrap(err)
	}
	for _, challengeInstance := range in.Instances {
		var dbInstance *pwdb.ChallengeInstance
		for _, instance := range dbInstances {
			if instance.ID == challengeInstance.ID {
				dbInstance = instance
			}
		}
		updated := false
		if !reflect.DeepEqual(dbInstance, challengeInstance) {
			updated = true
		}
		cpy := challengeInstance
		err := svc.db.Model(&cpy).
			Update(pwdb.ChallengeInstance{
				Status:         challengeInstance.Status,
				InstanceConfig: challengeInstance.InstanceConfig,
			}).
			Error
		if err != nil {
			return nil, errcode.ErrAgentUpdateState.Wrap(err)
		}
		if updated {
			activity := pwdb.Activity{
				Kind:                pwdb.Activity_AgentChallengeInstanceUpdate,
				AuthorID:            userID,
				AgentID:             challengeInstance.AgentID,
				ChallengeInstanceID: challengeInstance.ID,
				ChallengeFlavorID:   challengeInstance.FlavorID,
			}
			if err := svc.db.Create(&activity).Error; err != nil {
				return nil, errcode.ErrAgentUpdateState.Wrap(err)
			}
		}
	}

	if err != nil {
		return nil, errcode.ErrCommitUserTransaction.Wrap(err)
	}

	ret := &AgentUpdateState_Output{}
	return ret, nil
}
