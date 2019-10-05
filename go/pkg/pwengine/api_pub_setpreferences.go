package pwengine

import (
	"context"
	"fmt"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) SetPreferences(ctx context.Context, in *SetPreferencesInput) (*Void, error) {
	userID, err := subjectFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get subject from context: %w", err)
	}

	var (
		hasChanges = false
		updates    = map[string]interface{}{}
	)

	// update active tournament
	if in.ActiveTournamentID != "" {
		hasChanges = true

		// get active tournament
		var tournamentIDs []string
		err := e.db.
			Table("tournament").
			Where("ID = ?", in.ActiveTournamentID).
			Pluck("id", &tournamentIDs).
			Error
		switch {
		case err == nil && len(tournamentIDs) == 1:
			updates["active_tournament_id"] = tournamentIDs[0]
		case err == nil && len(tournamentIDs) == 0:
			return nil, ErrInvalidArgument
		default:
			return nil, fmt.Errorf("failed to get tournament: %w", err)
		}

		// get active tournament membership (optional)
		var tournamentMemberIDs []string
		err = e.db.
			Table("tournament_member").
			Joins("left join tournament_team ON tournament_team.id = tournament_member.tournament_team_id").
			Where("tournament_member.user_id = ?", userID).
			Where("tournament_team.tournament_id = ?", in.ActiveTournamentID).
			Pluck("tournament_member.id", &tournamentMemberIDs).
			Error
		switch {
		case err == nil && len(tournamentMemberIDs) == 1:
			updates["active_tournament_member_id"] = tournamentMemberIDs[0]
		case err == nil && len(tournamentMemberIDs) == 0:
			updates["active_tournament_member_id"] = "" // nil instead?
		default:
			return nil, fmt.Errorf("failed to get tournament team: %w", err)
		}
	}

	if !hasChanges {
		return nil, ErrMissingArgument
	}

	err = e.db.Model(pwdb.User{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	// FIXME: check amount of updated rows

	return nil, nil
}
