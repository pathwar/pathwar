package pwapi

import (
	"context"
	"strings"
	"time"

	"pathwar.land/v2/go/pkg/errcode"
	"pathwar.land/v2/go/pkg/pwdb"
)

func (svc *service) AgentRegister(ctx context.Context, in *AgentRegister_Input) (*AgentRegister_Output, error) {
	if !isAgentContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil || in.Name == "" {
		return nil, errcode.ErrMissingInput
	}

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
	agent.NginxPort = in.NginxPort
	agent.Metadata = in.Metadata
	agent.DomainSuffix = in.DomainSuffix
	agent.AuthSalt = in.AuthSalt
	agent.Status = pwdb.Agent_Active
	now := time.Now()
	agent.LastRegistrationAt = &now
	agent.LastSeenAt = &now
	agent.TimesSeen++
	agent.TimesRegistered++

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
