package pwes

import (
	"context"
	"strconv"
	"time"

	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func EventHandler(ctx context.Context, apiClient *pwapi.HTTPClient, timestamp *time.Time, logger *zap.Logger) error {
	if apiClient == nil || timestamp == nil {
		return errcode.ErrMissingInput
	}

	res, err := apiClient.AdminListActivities(ctx, &pwapi.AdminListActivities_Input{Since: timestamp})
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	activities := res.GetActivities()
	if len(activities) == 0 {
		return nil
	}

	//TODO: Handle other events
	for _, activity := range activities {
		var e Event
		switch activity.Kind {
		case pwdb.Activity_ChallengeSubscriptionValidate:
			e = &EventChallengeSubscriptionValidate{ID: activity.ID, CreatedAt: activity.CreatedAt, SeasonChallenge: activity.SeasonChallenge}
		}

		err = e.execute()
		if err != nil {
			logger.Debug("The event : " + strconv.Itoa(int(e.getID())) + " failed to execute.")
		}
		*timestamp = *e.getCreatedAt()
	}

	return nil
}
