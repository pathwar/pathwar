package pwdb

import (
	"fmt"
	"math/rand"

	"github.com/brianvoe/gofakeit"
	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"pathwar.land/go/internal/randstring"
)

func GetInfo(db *gorm.DB, logger *zap.Logger) (*Info, error) {
	info := Info{
		TableRows: make(map[string]uint32),
	}
	for _, model := range All() {
		var count uint32
		tableName := db.NewScope(model).TableName()
		if err := db.Model(model).Count(&count).Error; err != nil {
			logger.Warn("get table rows", zap.String("table", tableName), zap.Error(err))
			continue
		}
		info.TableRows[tableName] = count
	}

	return &info, nil
}

func GetDump(db *gorm.DB) (*Dump, error) {
	dump := Dump{}
	if err := db.Find(&dump.Challenges).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.ChallengeFlavors).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.ChallengeInstances).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.Hypervisors).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.Users).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.Organizations).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.OrganizationMembers).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.Seasons).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.Teams).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.TeamMembers).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.Coupons).Error; err != nil {
		return nil, err
	}
	return &dump, nil
}

func GenerateFakeData(db *gorm.DB, sfn *snowflake.Node, logger *zap.Logger) error {
	//
	// hypervisors
	//

	hypervisors := []*Hypervisor{}
	for i := 0; i < 3; i++ {
		hypervisor := &Hypervisor{
			Name:    gofakeit.HipsterWord(),
			Address: gofakeit.IPv4Address(),
			Status:  Hypervisor_Active,
		}
		hypervisors = append(hypervisors, hypervisor)
	}
	logger.Debug("Generating hypervisors")
	for _, entity := range hypervisors {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return fmt.Errorf("create hypervisors: %w", err)
		}
	}

	//
	// seasons
	//

	seasons := []*Season{}
	for i := 0; i < 3; i++ {
		season := &Season{
			Name:       gofakeit.HipsterWord(),
			Status:     Season_Started,
			Visibility: Season_Public,
			IsDefault:  false,
		}
		seasons = append(seasons, season)
	}
	seasons[0].IsDefault = true
	logger.Debug("Generating seasons")
	for _, entity := range seasons {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return fmt.Errorf("create seasons: %w", err)
		}
	}

	//
	// challenges
	//

	challenges := []*Challenge{}
	for i := 0; i < 5; i++ {
		challenge := &Challenge{
			Name:        gofakeit.HipsterWord(),
			Description: gofakeit.HipsterSentence(10),
			Author:      gofakeit.Name(),
			Locale:      "fr_FR",
			IsDraft:     false,
			Flavors:     []*ChallengeFlavor{},
		}
		for i := 0; i < 2; i++ {
			flavor := &ChallengeFlavor{
				Driver:           ChallengeFlavor_Docker,
				Version:          gofakeit.IPv4Address(),
				Changelog:        gofakeit.HipsterSentence(5),
				IsDraft:          false,
				IsLatest:         i == 0,
				SourceURL:        gofakeit.URL(),
				SeasonChallenges: []*SeasonChallenge{},
			}
			for j := 0; j < 2; j++ {
				seasonChallenge := &SeasonChallenge{
					SeasonID: seasons[rand.Int()%len(seasons)].ID,
				}
				for k := 0; k < 2; k++ {
					instance := &ChallengeInstance{
						HypervisorID: hypervisors[rand.Int()%len(hypervisors)].ID,
						Status:       ChallengeInstance_Active,
					}
					seasonChallenge.Instances = append(seasonChallenge.Instances, instance)
				}
				flavor.SeasonChallenges = append(flavor.SeasonChallenges, seasonChallenge)
			}
			challenge.Flavors = append(challenge.Flavors, flavor)
		}
		challenges = append(challenges, challenge)
	}

	logger.Debug("Generating challenges")
	for _, entity := range challenges {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return fmt.Errorf("create challenges: %w", err)
		}
	}

	//
	// organizations
	//

	organizations := []*Organization{}
	for i := 0; i < 5; i++ {
		organization := &Organization{
			Name:        gofakeit.HipsterWord(),
			GravatarURL: gofakeit.ImageURL(400, 400) + "?" + gofakeit.HipsterWord(),
			Locale:      "fr_FR",
		}
		organizations = append(organizations, organization)
	}

	logger.Debug("Generating organizations")
	for _, entity := range organizations {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return fmt.Errorf("create organizations: %w", err)
		}
	}

	//
	// users
	//

	users := []*User{}
	for i := 0; i < 10; i++ {
		user := &User{
			Username:                gofakeit.Name(),
			GravatarURL:             gofakeit.ImageURL(400, 400) + "?" + gofakeit.HipsterWord(),
			WebsiteURL:              gofakeit.URL(),
			Locale:                  "fr_FR",
			OrganizationMemberships: []*OrganizationMember{},
			OAuthSubject:            randstring.RandString(10),
		}
		users = append(users, user)
	}
	logger.Debug("Generating users")
	for _, entity := range users {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return fmt.Errorf("create users: %w", err)
		}
	}

	//
	// coupons
	//

	coupons := []*Coupon{}
	for i := 0; i < 3; i++ {
		coupon := &Coupon{
			Hash:               gofakeit.UUID(),
			MaxValidationCount: int32(rand.Int() % 5),
			Value:              int32(rand.Int() % 10),
			SeasonID:           seasons[rand.Int()%len(seasons)].ID,
		}
		coupons = append(coupons, coupon)
	}

	logger.Debug("Generating coupons")
	for _, entity := range coupons {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return fmt.Errorf("create coupons: %w", err)
		}
	}

	//
	// organization members
	//

	memberships := []*OrganizationMember{}
	for _, user := range users {
		for i := 0; i < 2; i++ {
			memberships = append(
				memberships,
				&OrganizationMember{
					OrganizationID: organizations[rand.Int()%len(organizations)].ID,
					UserID:         user.ID,
					Role:           OrganizationMember_Member, // or Owner
				},
			)
		}
	}

	logger.Debug("Generating memberships")
	for _, entity := range memberships {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return fmt.Errorf("create memberships: %w", err)
		}
	}

	return nil
}
