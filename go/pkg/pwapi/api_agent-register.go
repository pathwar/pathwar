package pwapi

import (
	"context"
	"strings"
	"time"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AgentRegister(ctx context.Context, in *AgentRegister_Input) (*AgentRegister_Output, error) {
	if !isAgentContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil || in.Name == "" {
		return nil, errcode.ErrMissingInput
	}

	userID, _ := userIDFromContext(ctx, svc.db)
	// userID is only used for activity logging, we don't care that it returns an error

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
	agent.NginxPort = int64(in.NginxPort)
	agent.Metadata = in.Metadata
	agent.DomainSuffix = in.DomainSuffix
	agent.AuthSalt = in.AuthSalt
	agent.Status = pwdb.Agent_Active
	now := time.Now()
	agent.LastRegistrationAt = &now
	agent.LastSeenAt = &now
	agent.TimesSeen++
	agent.TimesRegistered++
	agent.DefaultAgent = in.DefaultAgent

	// save last object with updated last_seen etc
	err = svc.db.Save(&agent).Error
	if err != nil {
		return nil, errcode.ErrSaveAgent.Wrap(err)
	}

	var challengeFlavorsToInstanciate []*pwdb.ChallengeFlavor
	// check if agent is registering for the first time
	if agent.TimesSeen == 1 {
		// if default flag, start an instance for each challenge flavor in database
		if agent.DefaultAgent {
			err = svc.db.
				Find(&challengeFlavorsToInstanciate).
				Error
			if err != nil {
				return nil, pwdb.GormToErrcode(err)
			}
		} else {
			// else just automatically start a challenge debug instance
			var debugFlavor pwdb.ChallengeFlavor
			err = svc.db.
				Where(pwdb.ChallengeFlavor{SourceURL: "https://github.com/pathwar/challenge-debug"}).
				First(&debugFlavor).
				Error
			if err != nil {
				return nil, pwdb.GormToErrcode(err)
			}
			challengeFlavorsToInstanciate = append(challengeFlavorsToInstanciate, &debugFlavor)
		}

		activity := pwdb.Activity{
			Kind:     pwdb.Activity_AgentRegister,
			AuthorID: userID,
			AgentID:  agent.ID,
		}
		err = svc.db.Create(&activity).Error
		if err != nil {
			return nil, pwdb.GormToErrcode(err)
		}
	}

	for _, challengeFlavor := range challengeFlavorsToInstanciate {
		instance := pwdb.ChallengeInstance{
			Status:   pwdb.ChallengeInstance_IsNew,
			AgentID:  agent.ID,
			FlavorID: challengeFlavor.ID,
		}
		err = svc.db.Create(&instance).Error
		if err != nil {
			return nil, pwdb.GormToErrcode(err)
		}
		activity := pwdb.Activity{
			Kind:                pwdb.Activity_AgentChallengeInstanceCreate,
			AuthorID:            userID,
			AgentID:             agent.ID,
			ChallengeInstanceID: instance.ID,
			ChallengeFlavorID:   challengeFlavor.ID,
		}
		if err := svc.db.Create(&activity).Error; err != nil {
			return nil, errcode.ErrAgentUpdateState.Wrap(err)
		}
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
