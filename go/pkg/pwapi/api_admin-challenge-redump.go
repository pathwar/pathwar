package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminChallengeRedump(ctx context.Context, in *AdminChallengeRedump_Input) (*AdminChallengeRedump_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	challengeID, err := pwdb.GetIDBySlugAndKind(svc.db, in.ChallengeID, "challenge")
	if err != nil {
		return nil, err
	}

	challengeFlavors := []pwdb.ChallengeFlavor{}
	err = svc.db.Where(&pwdb.ChallengeFlavor{ChallengeID: challengeID}).Find(&challengeFlavors).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	flavorIDs := []int64{}
	for _, challengeFlavor := range challengeFlavors {
		flavorIDs = append(flavorIDs, challengeFlavor.ID)
	}

	out := &AdminChallengeRedump_Output{}
	err = svc.db.
		Model(pwdb.ChallengeInstance{}).
		Where("flavor_id IN (?)", flavorIDs).
		Updates(&pwdb.ChallengeInstance{
			Status:         pwdb.ChallengeInstance_NeedRedump,
			InstanceConfig: []byte{},
		}).
		Find(&out.ChallengeInstances).
		Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	return out, nil
}
