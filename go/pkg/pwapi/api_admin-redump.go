package pwapi

import (
	"context"
	"strconv"

	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminRedump(ctx context.Context, in *AdminRedump_Input) (*AdminRedump_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	for _, identifier := range in.Identifiers {
		nb, err := strconv.Atoi(identifier)
		if err != nil {
			// for now, we only accept numerical identifiers, but the plan is to also search for names
			return nil, errcode.TODO.Wrap(err)
		}

		// FIXME: support passing IDs for other entities too
		var instance pwdb.ChallengeInstance
		err = svc.db.First(&instance, nb).Error
		if err != nil {
			return nil, errcode.TODO.Wrap(err)
		}

		switch instance.Status {
		case pwdb.ChallengeInstance_NeedRedump:
			svc.logger.Debug("level already marked as needing a redump", zap.Int64("instance-id", instance.ID))
		default:
			err = svc.db.Model(&instance).Updates(&pwdb.ChallengeInstance{Status: pwdb.ChallengeInstance_NeedRedump}).Error
			if err != nil {
				return nil, errcode.TODO.Wrap(err)
			}
			svc.logger.Debug("level marked as needing a redump", zap.Int64("instance-id", instance.ID))
		}
	}

	out := AdminRedump_Output{}
	return &out, nil
}
