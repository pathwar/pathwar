package pwdb

import (
	"fmt"

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
	if entity.Slug == "" {
		var challenge Challenge
		err := db.First(&challenge, "id = ?", entity.ChallengeID).Error
		if err != nil {
			return GormToErrcode(err)
		}
		entity.Slug = fmt.Sprintf("%s@%s", challenge.Slug, entity.Version)
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
