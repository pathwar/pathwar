package pwapi

import (
	"github.com/jinzhu/gorm"
	"pathwar.land/go/pkg/pwdb"
)

func userTeamForSeason(db *gorm.DB, userID, seasonID int64) (*pwdb.Team, error) {
	var team pwdb.Team

	err := db.
		Where(pwdb.Team{
			SeasonID:       seasonID,
			DeletionStatus: pwdb.DeletionStatus_Active,
		}).
		Joins("JOIN team_member ON team.id = team_member.team_id AND team_member.user_id = ?", userID).
		First(&team).
		Error

	return &team, err
}

func seasonFromSeasonChallengeID(db *gorm.DB, seasonChallengeID int64) (*pwdb.Season, error) {
	var seasonChallenge pwdb.SeasonChallenge

	err := db.
		Preload("Season").
		First(&seasonChallenge, seasonChallengeID).
		Error

	if err != nil {
		return nil, err
	}

	return seasonChallenge.Season, nil
}

func seasonIDExists(db *gorm.DB, seasonID int64) (bool, error) {
	var c int
	err := db.
		Table("season").
		Select("id").
		Where(&pwdb.Season{ID: seasonID}).
		Count(&c).
		Error
	if err != nil {
		return false, err
	}
	if c == 0 {
		return false, nil
	}
	return true, nil
}
