package pwapi

import (
	"context"

	"pathwar.land/v2/go/pkg/errcode"
	"pathwar.land/v2/go/pkg/pwdb"
)

func (svc *service) AdminChallengeFlavorAdd(ctx context.Context, in *AdminChallengeFlavorAdd_Input) (*AdminChallengeFlavorAdd_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	var challenge pwdb.Challenge
	err := svc.db.
		Where(pwdb.Challenge{ID: in.ChallengeFlavor.ChallengeID}).
		First(&challenge).
		Error
	if err != nil {
		return nil, errcode.ErrChallengeFlavorAdd.Wrap(err)
	}

	challengeFlavor := in.ChallengeFlavor

	err = svc.db.Create(challengeFlavor).Error
	if err != nil {
		return nil, errcode.ErrChallengeFlavorAdd.Wrap(err)
	}

	out := AdminChallengeFlavorAdd_Output{
		ChallengeFlavor: challengeFlavor,
	}
	return &out, nil
}
