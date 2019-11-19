package pwengine

import (
	"context"
	"fmt"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) UserSetPreferences(ctx context.Context, in *UserSetPreferences_Input) (*UserSetPreferences_Output, error) {
	userID, err := userIDFromContext(ctx, e.db)
	if err != nil {
		return nil, fmt.Errorf("get userid from context: %w", err)
	}

	var (
		hasChanges = false
		updates    = map[string]interface{}{}
	)

	// update active season
	if in.ActiveSeasonID != 0 {
		hasChanges = true

		// get active season
		var seasonIDs []string
		err := e.db.
			Table("season").
			Where("id = ?", in.ActiveSeasonID).
			Pluck("id", &seasonIDs).
			Error
		switch {
		case err == nil && len(seasonIDs) == 1:
			updates["active_season_id"] = seasonIDs[0]
		case err == nil && len(seasonIDs) == 0:
			return nil, ErrInvalidArgument
		default:
			return nil, fmt.Errorf("get season: %w", err)
		}

		// get active season membership (optional)
		var seasonMemberIDs []int64
		err = e.db.
			Table("team_member").
			Joins("left join team ON team.id = team_member.team_id").
			Where("team_member.user_id = ?", userID).
			Where("team.season_id = ?", in.ActiveSeasonID).
			Pluck("team_member.id", &seasonMemberIDs).
			Error
		switch {
		case err == nil && len(seasonMemberIDs) == 1:
			updates["active_team_member_id"] = seasonMemberIDs[0]
		case err == nil && len(seasonMemberIDs) == 0:
			updates["active_team_member_id"] = 0 // nil instead?
		default:
			return nil, fmt.Errorf("get season organization: %w", err)
		}
	}

	if !hasChanges {
		return nil, ErrMissingArgument
	}

	err = e.db.Model(pwdb.User{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	// FIXME: check amount of updated rows

	return nil, nil
}
