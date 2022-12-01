package pwes

import (
	"time"

	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

type Event interface {
	getID() int64
	getCreatedAt() *time.Time
	execute() error
}

type EventChallengeSubscriptionValidate struct {
	ID              int64
	CreatedAt       *time.Time
	SeasonChallenge *pwdb.SeasonChallenge
}

func (e EventChallengeSubscriptionValidate) getID() int64 {
	return e.ID
}

func (e EventChallengeSubscriptionValidate) getCreatedAt() *time.Time {
	return e.CreatedAt
}
