package pwapi

import (
	"context"

	"pathwar.land/v2/go/pkg/errcode"
	"pathwar.land/v2/go/pkg/pwdb"
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
	)

	// update active season
	if in.ActiveSeasonID != 0 {
		hasChanges = true

		exists, err := seasonIDExists(svc.db, in.ActiveSeasonID)
		if err != nil || !exists {
			return nil, errcode.ErrInvalidSeasonID.Wrap(err)
		}
		updates["active_season_id"] = in.ActiveSeasonID

		// get active season membership (if user already has a team for this season)
		var seasonMemberIDs []int64
		err = svc.db.
			Table("team_member").
			Joins("left join team ON team.id = team_member.team_id").
			Where("team_member.user_id = ?", userID).
			Where("team.season_id = ?", in.ActiveSeasonID).
			Pluck("team_member.id", &seasonMemberIDs).
			Error
		if err != nil || len(seasonMemberIDs) > 1 {
			return nil, errcode.ErrGetActiveSeasonMembership.Wrap(err)
		}
		if len(seasonMemberIDs) == 1 {
			updates["active_team_member_id"] = seasonMemberIDs[0]
		}
		if len(seasonMemberIDs) == 0 {
			updates["active_team_member_id"] = 0 // nil instead?
		}
	}

	if !hasChanges {
		return nil, errcode.ErrMissingInput
	}

	err = svc.db.Model(pwdb.User{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		return nil, errcode.ErrUpdateUser.Wrap(err)
	}
	// FIXME: check amount of updated rows

	return nil, nil
}
