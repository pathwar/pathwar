package pwapi

import (
	"context"

	"github.com/jinzhu/gorm"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) UserSetPreferences(ctx context.Context, in *UserSetPreferences_Input) (*UserSetPreferences_Output, error) {
	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated
	}

	var (
		hasChanges = false
		updates    = map[string]interface{}{}
		activity   = pwdb.Activity{}
	)

	// get ID thanks to the slug
	if in.ActiveSeasonSlug != "" {
		in.ActiveSeasonID, err = pwdb.GetIDBySlugAndKind(svc.db, in.ActiveSeasonSlug, "season")
		if err != nil {
			return nil, err
		}
	}
	// update active season
	if in.ActiveSeasonID != 0 {
		hasChanges = true

		exists, err := seasonIDExists(svc.db, in.ActiveSeasonID)
		if err != nil || !exists {
			return nil, errcode.ErrInvalidSeasonID.Wrap(err)
		}
		activity.SeasonID = in.ActiveSeasonID
		updates["active_season_id"] = in.ActiveSeasonID

		// get active season membership (if user already has a team for this season)
		var seasonMemberIDs []int64
		var teamIDs []int64
		err = svc.db.
			Table("team_member").
			Joins("left join team ON team.id = team_member.team_id").
			Where("team_member.user_id = ?", userID).
			Where("team.season_id = ?", in.ActiveSeasonID).
			Pluck("team.id", &teamIDs).
			Pluck("team_member.id", &seasonMemberIDs).
			Error
		if err != nil || len(seasonMemberIDs) > 1 {
			return nil, errcode.ErrGetActiveSeasonMembership.Wrap(err)
		}
		if len(seasonMemberIDs) == 1 {
			updates["active_team_member_id"] = seasonMemberIDs[0]
			activity.TeamMemberID = seasonMemberIDs[0]
			activity.TeamID = teamIDs[0]
		}
		if len(seasonMemberIDs) == 0 {
			return nil, errcode.ErrUserHasNoTeamForSeason // FIXME: should we create a team for the user?
		}
	}

	if !hasChanges {
		return nil, errcode.ErrMissingInput
	}

	err = svc.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(pwdb.User{}).Where("id = ?", userID).Updates(updates).Error
		if err != nil {
			return err
		}

		activity.Kind = pwdb.Activity_UserSetPreferences
		activity.AuthorID = userID
		activity.UserID = userID
		return tx.Create(&activity).Error
	})
	if err != nil {
		return nil, errcode.ErrUpdateUser.Wrap(err)
	}

	ret := UserSetPreferences_Output{}
	return &ret, nil
}
