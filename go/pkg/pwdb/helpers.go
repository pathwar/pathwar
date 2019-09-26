package pwdb

import (
	"encoding/base64"
	fmt "fmt"
	"math/rand"

	"github.com/brianvoe/gofakeit"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

func GetInfo(db *gorm.DB, logger *zap.Logger) (*Info, error) {
	info := Info{
		TableRows: make(map[string]uint32),
	}
	for _, model := range All() {
		var count uint32
		tableName := db.NewScope(model).TableName()
		if err := db.Model(model).Count(&count).Error; err != nil {
			logger.Warn("failed to get table rows", zap.String("table", tableName), zap.Error(err))
			continue
		}
		info.TableRows[tableName] = count
	}

	return &info, nil
}

func GetDump(db *gorm.DB) (*Dump, error) {
	dump := Dump{}
	if err := db.Find(&dump.Levels).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.LevelVersions).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.LevelFlavors).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.LevelInstances).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.Hypervisors).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.Users).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.Teams).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.TeamMembers).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.Tournaments).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.TournamentTeams).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.TournamentMembers).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&dump.Coupons).Error; err != nil {
		return nil, err
	}
	return &dump, nil
}

func GenerateFakeData(db *gorm.DB, logger *zap.Logger) error {
	hypervisors := []*Hypervisor{}
	for i := 0; i < 3; i++ {
		hypervisor := &Hypervisor{
			Name:    gofakeit.HipsterWord(),
			Address: gofakeit.IPv4Address(),
			Status:  Hypervisor_Active,
		}
		hypervisors = append(hypervisors, hypervisor)
	}

	levels := []*Level{}
	for i := 0; i < 5; i++ {
		level := &Level{
			Name:        gofakeit.HipsterWord(),
			Description: gofakeit.HipsterSentence(10),
			Author:      gofakeit.Name(),
			Locale:      "fr_FR",
			IsDraft:     false,
			Versions:    []*LevelVersion{},
		}
		for i := 0; i < 2; i++ {
			version := &LevelVersion{
				Driver:    LevelVersion_Docker,
				Version:   gofakeit.IPv4Address(),
				Changelog: gofakeit.HipsterSentence(5),
				IsDraft:   false,
				IsLatest:  i == 0,
				SourceURL: gofakeit.URL(),
				Flavors:   []*LevelFlavor{},
			}
			for j := 0; j < 2; j++ {
				flavor := &LevelFlavor{
					Instances: []*LevelInstance{},
				}
				for k := 0; k < 2; k++ {
					instance := &LevelInstance{
						HypervisorID: hypervisors[rand.Int()%len(hypervisors)].ID,
						Status:       LevelInstance_Active,
					}
					flavor.Instances = append(flavor.Instances, instance)
				}
				version.Flavors = append(version.Flavors, flavor)
			}
			level.Versions = append(level.Versions, version)
		}
		levels = append(levels, level)
	}

	teams := []*Team{}
	for i := 0; i < 5; i++ {
		team := &Team{
			Name:        gofakeit.HipsterWord(),
			GravatarURL: gofakeit.ImageURL(400, 400) + "?" + gofakeit.HipsterWord(),
			Locale:      "fr_FR",
		}
		teams = append(teams, team)
	}

	users := []*User{}
	for i := 0; i < 10; i++ {
		id, err := uuid.NewV4().MarshalBinary()
		if err != nil {
			return fmt.Errorf("failed to generate uuid: %w", err)
		}
		out := base64.StdEncoding.EncodeToString(id)
		user := &User{
			ID:          out,
			Username:    gofakeit.Name(),
			GravatarURL: gofakeit.ImageURL(400, 400) + "?" + gofakeit.HipsterWord(),
			WebsiteURL:  gofakeit.URL(),
			Locale:      "fr_FR",
			Memberships: []*TeamMember{},
		}
		users = append(users, user)
	}

	tournaments := []*Tournament{}
	for i := 0; i < 3; i++ {
		tournament := &Tournament{
			Name:       gofakeit.HipsterWord(),
			Status:     Tournament_Started,
			Visibility: Tournament_Public,
			IsDefault:  false,
		}
		tournaments = append(tournaments, tournament)
	}
	tournaments[0].IsDefault = true

	logger.Debug("Generating hypervisors")
	for _, entity := range hypervisors {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return fmt.Errorf("failed to create hypervisors: %w", err)
		}
	}
	logger.Debug("Generating users")
	for _, entity := range users {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return fmt.Errorf("failed to create users: %w", err)
		}
	}
	logger.Debug("Generating levels")
	for _, entity := range levels {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return fmt.Errorf("failed to create levels: %w", err)
		}
	}
	logger.Debug("Generating tournaments")
	for _, entity := range tournaments {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return fmt.Errorf("failed to create tournaments: %w", err)
		}
	}
	logger.Debug("Generating teams")
	for _, entity := range teams {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return fmt.Errorf("failed to create teams: %w", err)
		}
	}

	coupons := []*Coupon{}
	for i := 0; i < 3; i++ {
		coupon := &Coupon{
			Hash:               gofakeit.UUID(),
			MaxValidationCount: int32(rand.Int() % 5),
			Value:              int32(rand.Int() % 10),
			TournamentID:       tournaments[rand.Int()%len(tournaments)].ID,
		}
		coupons = append(coupons, coupon)
	}

	memberships := []*TeamMember{}
	for _, user := range users {
		for i := 0; i < 2; i++ {
			memberships = append(
				memberships,
				&TeamMember{
					TeamID: teams[rand.Int()%len(teams)].ID,
					UserID: user.ID,
					Role:   TeamMember_Member, // or Owner
				},
			)
		}
	}

	logger.Debug("Generating memberships")
	for _, entity := range memberships {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return fmt.Errorf("failed to create memberships: %w", err)
		}
	}
	logger.Debug("Generating coupons")
	for _, entity := range coupons {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return fmt.Errorf("failed to create coupons: %w", err)
		}
	}
	return nil
}
