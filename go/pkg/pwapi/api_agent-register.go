package pwapi

import (
	"context"
	"strings"
	"time"

	"pathwar.land/go/v2/pkg/errcode"
	"pathwar.land/go/v2/pkg/pwdb"
)

func (svc *service) AgentRegister(ctx context.Context, in *AgentRegister_Input) (*AgentRegister_Output, error) {
	if in == nil || in.Name == "" {
		return nil, errcode.ErrMissingInput
	}

	// FIXME: check if client is agent

	// check if agent already exists
	var agent pwdb.Agent
	err := svc.db.
		Where(pwdb.Agent{Name: in.Name}).
		First(&agent).
		Error
	if err != nil && !pwdb.IsRecordNotFoundError(err) {
		return nil, errcode.ErrGetAgent.Wrap(err)
	}

	// override it with input
	agent.Name = in.Name
	agent.Hostname = in.Hostname
	agent.OS = in.OS
	agent.Arch = in.Arch
	agent.Version = in.Version
	agent.Tags = strings.Join(in.Tags, ", ")
	now := time.Now()
	agent.LastRegistrationAt = &now
	agent.LastSeenAt = &now

	// save last object with updated last_seen etc
	err = svc.db.Save(&agent).Error
	if err != nil {
		return nil, errcode.ErrSaveAgent.Wrap(err)
	}

	// return the object
	err = svc.db.First(&agent, agent.ID).Error
	if err != nil {
		return nil, errcode.ErrGetAgent.Wrap(err)
	}

	out := &AgentRegister_Output{
		Agent: &agent,
	}
	return out, nil
}
