package pwapi

import (
	"context"

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

	err := svc.db.Create(in.ChallengeFlavor).Error
	if err != nil {
		return nil, errcode.ErrChallengeFlavorAdd.Wrap(err)
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
