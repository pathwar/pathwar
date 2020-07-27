package pwdb

import (
	"fmt"
	"strings"
)

func (entity *Team) ASCIIID() string {
	if entity == nil {
		return "-"
	}
	return entity.Slug
}

func (entity *User) ASCIIID() string {
	if entity == nil {
		return "-"
	}
	return entity.Slug
}

func (entity *Organization) ASCIIID() string {
	if entity == nil {
		return "-"
	}
	return entity.Slug
}

func (entity *Season) ASCIIID() string {
	if entity == nil {
		return "-"
	}
	return entity.Slug
}

func (entity *Challenge) ASCIIID() string {
	if entity == nil {
		return "-"
	}
	return entity.Slug
}

func (entity *Coupon) ASCIIID() string {
	if entity == nil {
		return "-"
	}
	return fmt.Sprintf("%d", entity.ID)
}

func (entity *Agent) ASCIIID() string {
	if entity == nil {
		return "-"
	}
	return entity.Slug
}

func (entity *SeasonChallenge) ASCIIID() string {
	if entity == nil {
		return "-"
	}
	return entity.Slug
}

func (entity *TeamMember) ASCIIID() string {
	if entity == nil {
		return "-"
	}
	return entity.Slug
}

func (entity *ChallengeSubscription) ASCIIID() string {
	if entity == nil {
		return "-"
	}
	return entity.Slug
}

func (entity *ChallengeFlavor) ASCIIID() string {
	if entity == nil {
		return "-"
	}
	return strings.TrimSuffix(entity.Slug, "@default")
}
