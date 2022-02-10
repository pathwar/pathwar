package pwdb

import (
	"math/rand"

	"github.com/brianvoe/gofakeit"
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"pathwar.land/pathwar/v2/go/internal/randstring"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func GetInfo(db *gorm.DB, logger *zap.Logger) (*Info, error) {
	info := Info{
		TableRows: make(map[string]int64),
	}
	stmt := &gorm.Statement{DB: db}
	for _, model := range All() {
		var count int64
		err := stmt.Parse(model)
		if err != nil {
			return nil, GormToErrcode(err)
		}
		tableName := stmt.Schema.Table
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
		return nil, GormToErrcode(err)
	}
	if err := db.Find(&dump.ChallengeFlavors).Error; err != nil {
		return nil, GormToErrcode(err)
	}
	if err := db.Find(&dump.ChallengeInstances).Error; err != nil {
		return nil, GormToErrcode(err)
	}
	if err := db.Find(&dump.Agents).Error; err != nil {
		return nil, GormToErrcode(err)
	}
	if err := db.Find(&dump.Users).Error; err != nil {
		return nil, GormToErrcode(err)
	}
	if err := db.Find(&dump.Organizations).Error; err != nil {
		return nil, GormToErrcode(err)
	}
	if err := db.Find(&dump.OrganizationMembers).Error; err != nil {
		return nil, GormToErrcode(err)
	}
	if err := db.Find(&dump.Seasons).Error; err != nil {
		return nil, GormToErrcode(err)
	}
	if err := db.Find(&dump.Teams).Error; err != nil {
		return nil, GormToErrcode(err)
	}
	if err := db.Find(&dump.TeamMembers).Error; err != nil {
		return nil, GormToErrcode(err)
	}
	if err := db.Find(&dump.Coupons).Error; err != nil {
		return nil, GormToErrcode(err)
	}
	return &dump, nil
}

func GenerateFakeData(db *gorm.DB, sfn *snowflake.Node, logger *zap.Logger) error {
	//
	// agents
	//

	agents := []*Agent{}
	for i := 0; i < 3; i++ {
		agent := &Agent{
			Name:     gofakeit.HipsterWord(),
			Hostname: gofakeit.IPv4Address(),
			Status:   Agent_Active,
		}
		agents = append(agents, agent)
	}
	logger.Debug("Generating agents")
	for _, entity := range agents {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return GormToErrcode(err)
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
			IsGlobal:   false,
		}
		seasons = append(seasons, season)
	}
	seasons[0].IsGlobal = true
	logger.Debug("Generating seasons")
	for _, entity := range seasons {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return GormToErrcode(err)
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
				SourceURL:        gofakeit.URL(),
				SeasonChallenges: []*SeasonChallenge{},
			}
			for j := 0; j < 2; j++ {
				seasonChallenge := &SeasonChallenge{
					SeasonID: seasons[rand.Int()%len(seasons)].ID,
				}
				flavor.SeasonChallenges = append(flavor.SeasonChallenges, seasonChallenge)
			}
			for j := 0; j < 2; j++ {
				instance := &ChallengeInstance{
					AgentID: agents[rand.Int()%len(agents)].ID,
					Status:  ChallengeInstance_Available,
				}
				flavor.Instances = append(flavor.Instances, instance)
			}
			challenge.Flavors = append(challenge.Flavors, flavor)
		}
		challenges = append(challenges, challenge)
	}

	logger.Debug("Generating challenges")
	for _, entity := range challenges {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return GormToErrcode(err)
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
			return GormToErrcode(err)
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
			return GormToErrcode(err)
		}
	}

	//
	// coupons
	//

	coupons := []*Coupon{}
	for i := 0; i < 3; i++ {
		coupon := &Coupon{
			Hash:               gofakeit.UUID(),
			MaxValidationCount: int64(rand.Int() % 5),
			Value:              int64(rand.Int() % 10),
			SeasonID:           seasons[rand.Int()%len(seasons)].ID,
		}
		coupons = append(coupons, coupon)
	}

	logger.Debug("Generating coupons")
	for _, entity := range coupons {
		if err := db.Set("gorm:association_autoupdate", true).Create(entity).Error; err != nil {
			return GormToErrcode(err)
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
			return GormToErrcode(err)
		}
	}

	return nil
}

func GetIDBySlugAndKind(db *gorm.DB, slug string, kind string) (int64, error) {
	var (
		ids []int64
		err error
	)
	switch kind {
	case "user":
		err = db.
			Model(User{}).
			Where("id = ? OR slug = ?", slug, slug).
			Pluck("id", &ids).Error
	case "organization":
		err = db.
			Model(Organization{}).
			Where("id = ? OR slug = ?", slug, slug).
			Pluck("id", &ids).Error
	case "team":
		err = db.
			Model(Team{}).
			Where("id = ? OR slug = ?", slug, slug).
			Pluck("id", &ids).Error
	case "team-invite":
		err = db.
			Model(TeamInvite{}).
			Where("id = ? OR slug = ?", slug, slug).
			Pluck("id", &ids).Error
	case "season":
		err = db.
			Model(Season{}).
			Where("id = ? OR slug = ?", slug, slug).
			Pluck("id", &ids).Error
	case "challenge":
		err = db.
			Model(Challenge{}).
			Where("id = ? OR slug = ?", slug, slug).
			Pluck("id", &ids).Error
	case "challenge-flavor":
		err = db.
			Model(ChallengeFlavor{}).
			Where("id = ? OR slug = ? OR slug = ?", slug, slug, slug+"@default").
			Pluck("id", &ids).Error
	case "challenge-instance":
		err = db.
			Model(ChallengeInstance{}).
			Where("id = ?", slug).
			Pluck("id", &ids).Error
	default:
		return 0, errcode.ErrUnknownDBKind
	}

	if err != nil {
		return 0, GormToErrcode(err)
	}
	if len(ids) == 0 {
		return 0, errcode.ErrNoSuchSlug
	}
	if len(ids) > 1 {
		return 0, errcode.ErrAmbiguousSlug
	}
	return ids[0], nil
}
