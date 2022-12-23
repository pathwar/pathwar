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

	logger.Info("event handler started", zap.Time("timestamp", *timestamp))
	res, err := apiClient.AdminListActivities(ctx, &pwapi.AdminListActivities_Input{Since: timestamp})
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	activities := res.GetActivities()
	if len(activities) == 0 {
		logger.Info("no activities to handle")
		return nil
	}

	// Use rebuild to process all events to be up-to-date in an efficient way
	if timestamp.IsZero() {
		logger.Info("Recompute all events from the beginning")
		*timestamp = *activities[len(activities)-1].CreatedAt
		err := Rebuild(ctx, apiClient, Opts{WithoutScore: false, From: "", To: timestamp.Format(TimeLayout), Logger: logger})
		if err != nil && err != errcode.ErrNothingToRebuild {
			return errcode.TODO.Wrap(err)
		}
		return nil
	}

	//TODO: Handle other events
	var e Event
	for _, activity := range activities {
		switch activity.Kind {
		case pwdb.Activity_ChallengeSubscriptionValidate:
			e = &EventChallengeSubscriptionValidate{ID: activity.ID, CreatedAt: activity.CreatedAt, SeasonChallenge: activity.SeasonChallenge, Team: activity.Team}
		case pwdb.Activity_Unknown:
			logger.Debug("The event : " + strconv.Itoa(int(e.getID())) + " is unknown kind.")
			continue
		default:
			continue
		}

		err = e.execute(ctx, apiClient)
		if err != nil {
			logger.Debug("The event : " + strconv.Itoa(int(e.getID())) + " failed to execute.")
		}
		*timestamp = *e.getCreatedAt()
	}

	return nil
}

// NewOpts returns same default values for development
func NewOpts() Opts {
	return Opts{
		WithoutScore: false,
		From:         "",
		To:           "",
		RefreshRate:  5,
		Logger:       zap.NewNop(),
	}
}

type Opts struct {
	WithoutScore bool
	From         string
	To           string
	RefreshRate  int
	Logger       *zap.Logger
}
