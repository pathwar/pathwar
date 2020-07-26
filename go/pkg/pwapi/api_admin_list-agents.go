package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminListAgents(ctx context.Context, in *AdminListAgents_Input) (*AdminListAgents_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	var agents []*pwdb.Agent
	err := svc.db.
		Preload("ChallengeInstances").
		Find(&agents).Error
	if err != nil {
		return nil, errcode.ErrListAgents.Wrap(err)
	}

	out := AdminListAgents_Output{Agents: agents}
	return &out, nil
}
