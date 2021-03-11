package pwdb

import (
	"fmt"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
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

func (entity *OrganizationMember) BeforeSave(db *gorm.DB) error {
	if entity.Slug == "" {
		var user User
		if entity.User == nil {
			err := db.First(&user, "id = ?", entity.UserID).Error
			if err != nil {
				return GormToErrcode(err)
			}
		} else {
			user = *entity.User
		}
		var organization Organization
		if entity.Organization == nil {
			err := db.First(&organization, "id = ?", entity.OrganizationID).Error
			if err != nil {
				return GormToErrcode(err)
			}
		} else {
			organization = *entity.Organization
		}
		entity.Slug = fmt.Sprintf("%s@%s", user.Slug, organization.Slug)
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

func (entity *Team) BeforeSave(db *gorm.DB) error {
	if entity.Slug == "" {
		var organization Organization
		if entity.Organization == nil {
			err := db.First(&organization, "id = ?", entity.OrganizationID).Error
			if err != nil {
				return GormToErrcode(err)
			}
		} else {
			organization = *entity.Organization
		}
		var season Season
		if entity.Season == nil {
			err := db.First(&season, "id = ?", entity.SeasonID).Error
			if err != nil {
				return GormToErrcode(err)
			}
		} else {
			season = *entity.Season
		}
		entity.Slug = fmt.Sprintf("%s@%s", organization.Slug, season.Slug)
	}
	return nil
}

func (entity *TeamMember) BeforeSave(db *gorm.DB) error {
	if entity.Slug == "" {
		var user User
		if entity.User == nil {
			err := db.First(&user, "id = ?", entity.UserID).Error
			if err != nil {
				return GormToErrcode(err)
			}
		} else {
			user = *entity.User
		}
		var team Team
		if entity.Team == nil {
			err := db.First(&team, "id = ?", entity.TeamID).Error
			if err != nil {
				return GormToErrcode(err)
			}
		} else {
			team = *entity.Team
		}
		entity.Slug = fmt.Sprintf("%s@%s", user.Slug, team.Slug)
	}
	return nil
}

func (entity *TeamInvite) BeforeSave(db *gorm.DB) error {
	if entity.Slug == "" {
		var user User
		if entity.User == nil {
			err := db.First(&user, "id = ?", entity.UserID).Error
			if err != nil {
				return GormToErrcode(err)
			}
		} else {
			user = *entity.User
		}
		var team Team
		if entity.Team == nil {
			err := db.First(&team, "id = ?", entity.TeamID).Error
			if err != nil {
				return GormToErrcode(err)
			}
		} else {
			team = *entity.Team
		}
		entity.Slug = fmt.Sprintf("%s@%s", user.Slug, team.Slug)
	}
	return nil
}

func (entity *SeasonChallenge) BeforeSave(db *gorm.DB) error {
	if entity.Slug == "" {
		var flavor ChallengeFlavor
		err := db.First(&flavor, "id = ?", entity.FlavorID).Error
		if err != nil {
			return GormToErrcode(err)
		}
		var season Season
		err = db.First(&season, "id = ?", entity.SeasonID).Error
		if err != nil {
			return GormToErrcode(err)
		}
		entity.Slug = fmt.Sprintf("%s@%s", flavor.Slug, season.Slug)
	}
	return nil
}

func (entity *ChallengeInstance) BeforeSave(db *gorm.DB) error {
	if entity.Slug == "" {
		var flavor ChallengeFlavor
		err := db.First(&flavor, "id = ?", entity.FlavorID).Error
		if err != nil {
			return GormToErrcode(err)
		}
		var agent Agent
		err = db.First(&agent, "id = ?", entity.AgentID).Error
		if err != nil {
			return GormToErrcode(err)
		}
		entity.Slug = fmt.Sprintf("%s@%s", flavor.Slug, agent.Slug)
	}
	return nil
}

func (entity *ChallengeSubscription) BeforeSave(db *gorm.DB) error {
	if entity.Slug == "" {
		var user User
		err := db.First(&user, "id = ?", entity.BuyerID).Error
		if err != nil {
			return GormToErrcode(err)
		}
		var seasonChallenge SeasonChallenge
		err = db.First(&seasonChallenge, "id = ?", entity.SeasonChallengeID).Error
		if err != nil {
			return GormToErrcode(err)
		}
		entity.Slug = fmt.Sprintf("%s@%s", user.Slug, seasonChallenge.Slug)
	}
	return nil
}

func (entity *ChallengeValidation) BeforeSave(db *gorm.DB) error {
	if entity.Slug == "" {
		var subscription ChallengeSubscription
		err := db.First(&subscription, "id = ?", entity.ChallengeSubscriptionID).Error
		if err != nil {
			return GormToErrcode(err)
		}
		entity.Slug = subscription.Slug
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
