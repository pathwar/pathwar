package pwapi

import (
	"context"
	"strings"

	"gorm.io/gorm"
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
		switch {
		case err != nil && strings.Contains(err.Error(), "Error 1062: Duplicate entry"):
			// flavor already exists, update it
			var existing pwdb.ChallengeFlavor
			if err := svc.db.First(&existing, "slug = ?", in.ChallengeFlavor.Slug).Error; err != nil {
				return pwdb.GormToErrcode(err)
			}
			in.ChallengeFlavor.ID = 0
			in.ChallengeFlavor.CreatedAt = nil
			in.ChallengeFlavor.UpdatedAt = nil
			in.ChallengeFlavor.Slug = ""
			in.ChallengeFlavor.ChallengeID = 0
			in.ChallengeFlavor.Version = ""
			if err = svc.db.Model(&existing).Updates(in.ChallengeFlavor).Error; err != nil {
				return pwdb.GormToErrcode(err)
			}
			// FIXME: need redump if compose bundle changes
			in.ChallengeFlavor = &existing
			return nil
		case err != nil:
			return errcode.ErrChallengeFlavorAdd.Wrap(err)
		}

		// testing seasons
		{
			var testingSeasons []*pwdb.Season
			if err = tx.Where(pwdb.Season{IsTesting: true}).Find(&testingSeasons).Error; err != nil {
				return err
			}

			for _, season := range testingSeasons {
				seasonChallenge := pwdb.SeasonChallenge{
					SeasonID: season.ID,
					FlavorID: in.ChallengeFlavor.ID,
				}
				err = tx.Create(&seasonChallenge).Error
				if err != nil {
					return pwdb.GormToErrcode(err)
				}
			}
		}

		// default agents
		{
			var agentsToInstanciate []*pwdb.Agent
			if err = tx.Where(pwdb.Agent{DefaultAgent: true}).Find(&agentsToInstanciate).Error; err != nil {
				return err
			}

			for _, agent := range agentsToInstanciate {
				instance := pwdb.ChallengeInstance{
					Status:   pwdb.ChallengeInstance_IsNew,
					AgentID:  agent.ID,
					FlavorID: in.ChallengeFlavor.ID,
				}
				err = tx.Create(&instance).Error
				if err != nil {
					return pwdb.GormToErrcode(err)
				}
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
		in.ChallengeFlavor.Version = "default"
	}
	if in.ChallengeFlavor.Category == "" {
		in.ChallengeFlavor.Category = "uncategorized"
	}
	if in.ChallengeFlavor.RedumpPolicyConfig == "" {
		in.ChallengeFlavor.RedumpPolicyConfig = `{"strategy":"on-validation"}`
	}
	if in.ChallengeFlavor.Passphrases == 0 {
		in.ChallengeFlavor.Passphrases = 1
	}
}
