package pwdb

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

func (entity *Challenge) BeforeSave(db *gorm.DB) error {
	if entity.Slug == "" {
		entity.Slug = slug.Make(entity.Name)
	}
	return nil
}

func (entity *Season) BeforeSave(db *gorm.DB) error {
	if entity.Slug == "" {
		entity.Slug = slug.Make(entity.Name)
	}
	return nil
}

func (entity *Agent) BeforeSave(db *gorm.DB) error {
	if entity.Slug == "" {
		entity.Slug = slug.Make(entity.Name)
	}
	return nil
}

func (entity *Organization) BeforeSave(db *gorm.DB) error {
	if entity.Slug == "" {
		entity.Slug = slug.Make(entity.Name)
	}
	return nil
}

func (entity *ChallengeFlavor) BeforeSave(db *gorm.DB) error {
	// FIXME: make a mix of challenge and flavor
	if entity.Slug == "" {
		entity.Slug = entity.Version
	}
	return nil
}

func (entity *User) BeforeSave(db *gorm.DB) error {
	if entity.Slug == "" {
		entity.Slug = slug.Make(entity.Username)
	}
	return nil
}

/*
func (entity *Team) BeforeSave(db *gorm.DB) error {
        // FIXME: make a join of orga and season
	if entity.Slug == "" {
		entity.Slug = slug.Make(entity.Username)
	}
	return nil
}
*/
