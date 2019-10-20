package pwengine

import (
	"context"
)

func (e *engine) ListTeams(context.Context, *Void) (*ListTeamsOutput, error) {
	var teams ListTeamsOutput
	if err := e.db.Set("gorm:auto_preload", true).Find(&teams.Items).Error; err != nil {
		return nil, err
	}

	return &teams, nil
}
