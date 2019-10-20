package pwengine

import (
	"context"
	"fmt"
)

func (e *engine) ListTeams(context.Context, *Void) (*ListTeamsOutput, error) {
	var teams ListTeamsOutput
	err := e.db.
		Set("gorm:auto_preload", true). // FIXME: explicit preloading
		Find(&teams.Items).Error
	if err != nil {
		return nil, fmt.Errorf("query teams: %w", err)
	}

	return &teams, nil
}
