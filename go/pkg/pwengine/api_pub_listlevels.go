package pwengine

import (
	"context"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) ListLevels(context.Context, *Void) (*pwdb.LevelList, error) {
	var levels pwdb.LevelList
	if err := e.db.Set("gorm:auto_preload", true).Find(&levels.Items).Error; err != nil {
		return nil, err
	}

	return &levels, nil
}
