package pwapi

import (
	"context"

	"github.com/jinzhu/gorm"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminChallengeFlavorAdd(ctx context.Context, in *AdminChallengeFlavorAdd_Input) (*AdminChallengeFlavorAdd_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	in.ApplyDefaults()
	if in == nil || (in.ChallengeID == "" && in.ChallengeFlavor.ChallengeID == 0) {
		return nil, errcode.ErrMissingInput
	}

	if in.ChallengeID != "" && in.ChallengeFlavor.ChallengeID == 0 {
		var err error
		in.ChallengeFlavor.ChallengeID, err = pwdb.GetIDBySlugAndKind(svc.db, in.ChallengeID, "challenge")
		if err != nil {
			return nil, err
		}
	}

	err := svc.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(in.ChallengeFlavor).Error
		if err != nil {
			return errcode.ErrChallengeFlavorAdd.Wrap(err)
		}

		var agentsToInstanciate []*pwdb.Agent
		if err = tx.Where(pwdb.Agent{DefaultAgent: true}).Find(&agentsToInstanciate).Error; err != nil {
			return err
		}

		for _, agent := range agentsToInstanciate {
			instance := pwdb.ChallengeInstance{
				Status:         pwdb.ChallengeInstance_IsNew,
				AgentID:        agent.ID,
				FlavorID:       in.ChallengeFlavor.ID,
				InstanceConfig: []byte(`{"passphrases": ["a", "b", "c", "d"]}`),
			}
			err = tx.Create(&instance).Error
			if err != nil {
				return pwdb.GormToErrcode(err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	out := AdminChallengeFlavorAdd_Output{
		ChallengeFlavor: in.ChallengeFlavor,
	}
	return &out, nil
}

func (in *AdminChallengeFlavorAdd_Input) ApplyDefaults() {
	if in == nil {
		return
	}
	if in.ChallengeFlavor == nil {
		in.ChallengeFlavor = &pwdb.ChallengeFlavor{}
	}
	if in.ChallengeFlavor.Version == "" {
		in.ChallengeFlavor.Version = "v1.0.0"
	}
}
