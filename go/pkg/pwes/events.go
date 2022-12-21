package pwes

import (
	"context"
	"time"

	"pathwar.land/pathwar/v2/go/pkg/pwapi"

	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

type Event interface {
	getID() int64
	getCreatedAt() *time.Time
	execute(ctx context.Context, apiClient *pwapi.HTTPClient) error
}

type EventChallengeSubscriptionValidate struct {
	ID              int64
	CreatedAt       *time.Time
	SeasonChallenge *pwdb.SeasonChallenge
	Team            *pwdb.Team
}

func (e EventChallengeSubscriptionValidate) getID() int64 {
	return e.ID
}

func (e EventChallengeSubscriptionValidate) getCreatedAt() *time.Time {
	return e.CreatedAt
}
