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

type Event interface {
	GetID() int64
	GetCreatedAt() *time.Time
	execute(ctx context.Context, apiClient *pwapi.HTTPClient, logger *zap.Logger) error
}

func EventHandler(ctx context.Context, apiClient *pwapi.HTTPClient, timestamp *time.Time, logger *zap.Logger) error {
	if apiClient == nil || timestamp == nil {
		return errcode.ErrMissingInput
	}

	// Get all events from timestamp to now -1 seconds in order to avoid to miss an event
	to := time.Now()
	to = to.Add(-time.Second)
	logger.Info("event handler started", zap.Time("timestamp", *timestamp))
	res, err := apiClient.AdminListActivities(ctx, &pwapi.AdminListActivities_Input{Since: timestamp, To: &to})
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
		case pwdb.Activity_UserRegister:
			e = &EventUserRegister{ID: activity.ID, CreatedAt: activity.CreatedAt}
		case pwdb.Activity_UserLogin:
			e = &EventUserLogin{ID: activity.ID, CreatedAt: activity.CreatedAt}
		case pwdb.Activity_UserSetPreferences:
			e = &EventUserSetPreferences{ID: activity.ID, CreatedAt: activity.CreatedAt}
		case pwdb.Activity_UserDeleteAccount:
			e = &EventUserDeleteAccount{ID: activity.ID, CreatedAt: activity.CreatedAt}
		case pwdb.Activity_SeasonChallengeBuy:
			e = &EventSeasonChallengeBuy{ID: activity.ID, CreatedAt: activity.CreatedAt}
		case pwdb.Activity_ChallengeSubscriptionValidate:
			e = &EventChallengeSubscriptionValidate{ID: activity.ID, CreatedAt: activity.CreatedAt, SeasonChallenge: activity.SeasonChallenge, Team: activity.Team}
		case pwdb.Activity_CouponValidate:
			e = &EventCouponValidate{ID: activity.ID, CreatedAt: activity.CreatedAt}
		case pwdb.Activity_AgentRegister:
			e = &EventAgentRegister{ID: activity.ID, CreatedAt: activity.CreatedAt}
		case pwdb.Activity_AgentChallengeInstanceCreate:
			e = &EventAgentChallengeInstanceCreate{ID: activity.ID, CreatedAt: activity.CreatedAt}
		case pwdb.Activity_AgentChallengeInstanceUpdate:
			e = &EventAgentChallengeInstanceUpdate{ID: activity.ID, CreatedAt: activity.CreatedAt}
		case pwdb.Activity_TeamCreation:
			e = &EventTeamCreation{ID: activity.ID, CreatedAt: activity.CreatedAt}
		case pwdb.Activity_TeamInviteSend:
			e = &EventTeamInviteSend{ID: activity.ID, CreatedAt: activity.CreatedAt}
		case pwdb.Activity_TeamInviteAccept:
			e = &EventTeamInviteAccept{ID: activity.ID, CreatedAt: activity.CreatedAt}
		case pwdb.Activity_Unknown:
			logger.Debug("The event : " + strconv.Itoa(int(e.GetID())) + " is unknown kind.")
			continue
		default:
			continue
		}

		err = e.execute(ctx, apiClient, logger)
		if err != nil {
			logger.Debug("The event : " + strconv.Itoa(int(e.GetID())) + " failed to execute.")
		}
		*timestamp = *e.GetCreatedAt()
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
