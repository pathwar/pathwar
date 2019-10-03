package pwengine

import (
	"context"

	"pathwar.land/go/pkg/pwdb"
)

func (c *client) ListLevels(context.Context, *Void) (*pwdb.LevelList, error) {
	var levels pwdb.LevelList
	if err := c.db.Set("gorm:auto_preload", true).Find(&levels.Items).Error; err != nil {
		return nil, err
	}

	return &levels, nil
}
