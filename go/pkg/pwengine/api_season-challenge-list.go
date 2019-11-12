package pwengine

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) SeasonChallengeList(ctx context.Context, in *SeasonChallengeListInput) (*SeasonChallengeListOutput, error) {
	if in == nil || in.SeasonID == 0 {
		return nil, ErrMissingArgument
	}

	userID, err := userIDFromContext(ctx, e.db)
	if err != nil {
		return nil, fmt.Errorf("get userid from context: %w", err)
	}

	var c int
	err = e.db.
		Table("season").
		Select("id").
		Where(&pwdb.Season{ID: in.SeasonID}).
		Count(&c).
		Error
	if err != nil {
		return nil, fmt.Errorf("fetch season: %w", err)
	}
	if c == 0 {
		return nil, ErrInvalidArgument // invalid in.SeasonID
	}

	team, err := userTeamForSeason(e.db, userID, in.SeasonID)
	if err != nil {
		return nil, ErrInvalidArgument // user does not have team for this season
	}

	var ret SeasonChallengeListOutput
	err = e.db.
		Preload("Season").
		Preload("Flavor").
		Preload("Flavor.Challenge").
		Preload("Subscriptions", "team_id = ?", team.ID).
		Preload("Subscriptions.Team").
		Where(pwdb.SeasonChallenge{SeasonID: in.SeasonID}).
		Find(&ret.Items).
		Error
	if err != nil {
		return nil, fmt.Errorf("fetch season challenges: %w", err)
	}

	return &ret, nil
}

func userTeamForSeason(db *gorm.DB, userID, seasonID int64) (*pwdb.Team, error) {
	var team pwdb.Team

	err := db.
		Where(pwdb.Team{SeasonID: seasonID}).
		Joins("JOIN team_member ON team.id = team_member.team_id AND team_member.user_id = ?", userID).
		First(&team).
		Error

	return &team, err
}
