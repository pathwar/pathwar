package pwengine

import (
	"context"

	"pathwar.land/go/pkg/pwdb"
)

func (c *client) ListTeams(context.Context, *Void) (*pwdb.TeamList, error) {
	var teams pwdb.TeamList
	if err := c.db.Set("gorm:auto_preload", true).Find(&teams.Items).Error; err != nil {
		return nil, err
	}

	return &teams, nil
}
