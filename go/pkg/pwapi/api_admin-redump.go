package pwapi

import (
	"context"
	"fmt"

	"go.uber.org/multierr"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminRedump(ctx context.Context, in *AdminRedump_Input) (*AdminRedump_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	var errs error
	for _, identifier := range in.Identifiers {
		instances := []int64{}
		if id, err := pwdb.GetIDBySlugAndKind(svc.db, identifier, "challenge-instance"); err == nil {
			instances = append(instances, id)
		} else if id, err := pwdb.GetIDBySlugAndKind(svc.db, identifier, "challenge-flavor"); err == nil {
			var flavor pwdb.ChallengeFlavor
			if err = svc.db.Preload("Instances").First(&flavor, id).Error; err != nil {
				errs = multierr.Append(errs, err)
				continue
			}
			for _, instance := range flavor.Instances {
				instances = append(instances, instance.ID)
			}
		}

		if len(instances) == 0 {
			errs = multierr.Append(errs, fmt.Errorf("no such entry for %q", identifier))
			continue
		}

		err := svc.db.
			Model(pwdb.ChallengeInstance{}).
			Where("id IN (?)", instances).
			Updates(&pwdb.ChallengeInstance{Status: pwdb.ChallengeInstance_NeedRedump}).
			Error
		if err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	out := AdminRedump_Output{}
	return &out, errs
}
