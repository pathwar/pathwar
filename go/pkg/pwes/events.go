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

type EventUserRegister struct {
	ID        int64
	CreatedAt *time.Time
}

func (e EventUserRegister) getID() int64 {
	return e.ID
}

func (e EventUserRegister) getCreatedAt() *time.Time {
	return e.CreatedAt
}

type EventUserLogin struct {
	ID        int64
	CreatedAt *time.Time
}

func (e EventUserLogin) getID() int64 {
	return e.ID
}

func (e EventUserLogin) getCreatedAt() *time.Time {
	return e.CreatedAt
}

type EventUserSetPreferences struct {
	ID        int64
	CreatedAt *time.Time
}

func (e EventUserSetPreferences) getID() int64 {
	return e.ID
}

func (e EventUserSetPreferences) getCreatedAt() *time.Time {
	return e.CreatedAt
}

type EventUserDeleteAccount struct {
	ID        int64
	CreatedAt *time.Time
}

func (e EventUserDeleteAccount) getID() int64 {
	return e.ID
}

func (e EventUserDeleteAccount) getCreatedAt() *time.Time {
	return e.CreatedAt
}

type EventSeasonChallengeBuy struct {
	ID        int64
	CreatedAt *time.Time
}

func (e EventSeasonChallengeBuy) getID() int64 {
	return e.ID
}

func (e EventSeasonChallengeBuy) getCreatedAt() *time.Time {
	return e.CreatedAt
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

type EventCouponValidate struct {
	ID        int64
	CreatedAt *time.Time
}

func (e EventCouponValidate) getID() int64 {
	return e.ID
}

func (e EventCouponValidate) getCreatedAt() *time.Time {
	return e.CreatedAt
}

type EventAgentRegister struct {
	ID        int64
	CreatedAt *time.Time
}

func (e EventAgentRegister) getID() int64 {
	return e.ID
}

func (e EventAgentRegister) getCreatedAt() *time.Time {
	return e.CreatedAt
}

type EventAgentChallengeInstanceCreate struct {
	ID        int64
	CreatedAt *time.Time
}

func (e EventAgentChallengeInstanceCreate) getID() int64 {
	return e.ID
}

func (e EventAgentChallengeInstanceCreate) getCreatedAt() *time.Time {
	return e.CreatedAt
}

type EventAgentChallengeInstanceUpdate struct {
	ID        int64
	CreatedAt *time.Time
}

func (e EventAgentChallengeInstanceUpdate) getID() int64 {
	return e.ID
}

func (e EventAgentChallengeInstanceUpdate) getCreatedAt() *time.Time {
	return e.CreatedAt
}

type EventTeamCreation struct {
	ID        int64
	CreatedAt *time.Time
}

func (e EventTeamCreation) getID() int64 {
	return e.ID
}

func (e EventTeamCreation) getCreatedAt() *time.Time {
	return e.CreatedAt
}

type EventTeamInviteSend struct {
	ID        int64
	CreatedAt *time.Time
}

func (e EventTeamInviteSend) getID() int64 {
	return e.ID
}

func (e EventTeamInviteSend) getCreatedAt() *time.Time {
	return e.CreatedAt
}

type EventTeamInviteAccept struct {
	ID        int64
	CreatedAt *time.Time
}

func (e EventTeamInviteAccept) getID() int64 {
	return e.ID
}

func (e EventTeamInviteAccept) getCreatedAt() *time.Time {
	return e.CreatedAt
}
