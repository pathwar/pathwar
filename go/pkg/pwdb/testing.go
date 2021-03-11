package pwdb

import (
	"testing"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

func TestingSqliteDB(t *testing.T) *gorm.DB {
	t.Helper()

	// logger
	logger := zapgorm2.New(zap.L())
	logger.SetAsDefault()
	DefaultGormConfig.Logger = logger

	// disable foreignKey
	//DefaultGormConfig.DisableForeignKeyConstraintWhenMigrating = true

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &DefaultGormConfig)
	if err != nil {
		t.Fatalf("init in-memory sqlite server: %v", err)
	}

	sfn, err := snowflake.NewNode(1)
	if err != nil {
		t.Fatalf("init snowflake generator: %v", err)
	}

	db, err = Configure(db, sfn)
	if err != nil {
		t.Fatalf("init pwdb: %v", err)
	}

	TestingCreateEntities(t, db)

	return db
}

// FIXME: func TestingMySQLDB(t *testing.T, logger *zap.Logger) *gorm.DB { }

func TestingCreateEntities(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.Transaction(func(tx *gorm.DB) error {
		// load base objects
		globalSeason := Season{IsGlobal: true}
		err := tx.Where(&globalSeason).First(&globalSeason).Error
		if err != nil {
			return GormToErrcode(err)
		}
		staffOrg := Organization{Name: "Staff"}
		err = tx.Where(&staffOrg).First(&staffOrg).Error
		if err != nil {
			return GormToErrcode(err)
		}
		staffTeam := Team{SeasonID: globalSeason.ID, OrganizationID: staffOrg.ID}
		err = tx.Where(&staffTeam).First(&staffTeam).Error
		if err != nil {
			return GormToErrcode(err)
		}
		hackSparrow := &User{OAuthSubject: "Hack Sparrow"}
		err = tx.Where(&hackSparrow).First(&hackSparrow).Error
		if err != nil {
			return GormToErrcode(err)
		}

		// seasons
		testSeason := &Season{
			Name:       "Unit Test Season",
			Status:     Season_Started,
			Visibility: Season_Public,
		}
		err = tx.Create(testSeason).Error
		if err != nil {
			return GormToErrcode(err)
		}

		// agents
		dummyAgent1 := &Agent{
			Name:         "dummy-agent-1",
			Hostname:     "dummy-agent-1.com",
			Status:       Agent_Active, // only useful during dev
			DomainSuffix: "local",
			AuthSalt:     "bluh",
		}
		dummyAgent2 := &Agent{
			Name:         "dummy-agent-2",
			Hostname:     "dummy-agent-2.com",
			Status:       Agent_Active,
			DomainSuffix: "local",
			NginxPort:    4242,
			AuthSalt:     "blah",
		}
		dummyAgent3 := &Agent{
			Name:         "dummy-agent-3",
			Hostname:     "dummy-agent-3.com",
			Status:       Agent_Inactive,
			DomainSuffix: "local",
			AuthSalt:     "blih",
		}
		for _, agent := range []*Agent{dummyAgent1, dummyAgent2, dummyAgent3} {
			err := tx.Create(agent).Error
			if err != nil {
				return GormToErrcode(err)
			}
		}

		// challenges
		bundle := `version: "3.7"
networks: {}
volumes: {}
services:
  gotty:
    image: pathwar/challenge-debug@sha256:a5f48f8c2eaf7f5cd106b2d115cfce351bb64653211f45cfe4829b668d3547f3
    ports:
      - "8080"
    labels:
      land.pathwar.compose.challenge-name: challenge-debug
      land.pathwar.compose.challenge-version: 1.0.0
      land.pathwar.compose.origin: was-built
      land.pathwar.compose.service-name: gotty
`
		challengeDebug := newOfficialChallengeWithFlavor("Debug", "https://github.com/pathwar/challenge-debug", bundle)
		challengeDebug.PurchasePrice = 0
		challengeDebug.addSeasonChallengeByID(globalSeason.ID)

		bundle = `version: "3.7"
networks: {}
volumes: {}
services:
    front:
        image: pathwar/helloworld@sha256:bf7a6384b4f19127ca1dd3c383d695030478e6d68ec27f24bb83edc42a5f3d26
        ports:
            - "80"
        labels:
            land.pathwar.compose.challenge-name: helloworld
            land.pathwar.compose.challenge-version: 1.0.0
            land.pathwar.compose.origin: was-built
            land.pathwar.compose.service-name: front
`
		helloworld := newOfficialChallengeWithFlavor("Hello World", "https://github.com/pathwar/pathwar/tree/master/challenges/web/helloworld", bundle)
		helloworld.PurchasePrice = 0
		helloworld.addSeasonChallengeByID(globalSeason.ID)

		bundle = `version: "3.7"
networks: {}
volumes: {}
services:
    front:
        image: pathwar/training-sqli@sha256:77c49c7907e19cd92baf2d6278dd017d2f5f6b9d6214d308694fba1572693545
        ports:
            - "80"
        depends_on:
            - mysql
        labels:
            land.pathwar.compose.challenge-name: training-sqli
            land.pathwar.compose.challenge-version: 1.0.0
            land.pathwar.compose.origin: was-built
            land.pathwar.compose.service-name: front
    mysql:
        image: pathwar/training-sqli@sha256:914ee0d8bf48e176b378c43ad09751c341d0266381e76ae12c385fbc6beb5983
        expose:
            - "3306"
        labels:
            land.pathwar.compose.challenge-name: training-sqli
            land.pathwar.compose.challenge-version: 1.0.0
            land.pathwar.compose.origin: was-built
            land.pathwar.compose.service-name: mysql
`
		trainingSQLI := newOfficialChallengeWithFlavor("Training SQLI", "https://github.com/pathwar/pathwar/tree/master/challenges/web/training-sqli", bundle)
		trainingSQLI.addSeasonChallengeByID(globalSeason.ID)

		bundle = `version: "3.7"
networks: {}
volumes: {}
services:
    front:
        image: pathwar/training-http@sha256:92c46270f8d7be9d927345353b7ea49b37dbf6c82ab6b2da3bc401f9fbacf5e5
        ports:
          - "80"
        labels:
            land.pathwar.compose.challenge-name: training-http
            land.pathwar.compose.challenge-version: 1.0.0
            land.pathwar.compose.origin: was-built
            land.pathwar.compose.service-name: front
`

		trainingHTTP := newOfficialChallengeWithFlavor("Training HTTP", "https://github.com/pathwar/pathwar/tree/master/challenges/web/training-http", bundle)
		trainingHTTP.addSeasonChallengeByID(globalSeason.ID)

		for _, flavor := range []*ChallengeFlavor{
			challengeDebug, helloworld, trainingSQLI, trainingHTTP,
		} {
			err = tx.Set("gorm:association_autoupdate", false).Create(flavor.Challenge).Error
			if err != nil {
				return GormToErrcode(err)
			}
			flavor.ChallengeID = flavor.Challenge.ID

			err = tx.Set("gorm:association_autoupdate", false).Create(flavor).Error
			if err != nil {
				return GormToErrcode(err)
			}
		}

		dummyChallenge1 := newOfficialChallengeWithFlavor("dummy challenge 1", "https://...", "")
		dummyChallenge1.PurchasePrice = 0
		dummyChallenge1.addSeasonChallengeByID(globalSeason.ID)
		dummyChallenge1.addSeasonChallengeByID(testSeason.ID)
		dummyChallenge2 := newOfficialChallengeWithFlavor("dummy challenge 2", "https://...", "")
		dummyChallenge2.addSeasonChallengeByID(globalSeason.ID)
		dummyChallenge2.addSeasonChallengeByID(testSeason.ID)
		dummyChallenge3 := newOfficialChallengeWithFlavor("dummy challenge 3", "https://...", "")
		dummyChallenge3.addSeasonChallengeByID(globalSeason.ID)
		dummyChallenge3.addSeasonChallengeByID(testSeason.ID)

		for _, flavor := range []*ChallengeFlavor{
			dummyChallenge1, dummyChallenge2, dummyChallenge3,
		} {
			err := tx.Create(flavor.Challenge).Error
			if err != nil {
				return GormToErrcode(err)
			}
			flavor.ChallengeID = flavor.Challenge.ID
			err = tx.Create(flavor).Error
			if err != nil {
				return GormToErrcode(err)
			}
		}

		// Challenge Instances
		devConfig := []byte(`{"passphrases": ["a", "b", "c", "d"]}`)
		instances := []*ChallengeInstance{
			{Status: ChallengeInstance_Available, AgentID: dummyAgent1.ID, FlavorID: dummyChallenge1.ID, InstanceConfig: devConfig},
			{Status: ChallengeInstance_Available, AgentID: dummyAgent2.ID, FlavorID: dummyChallenge2.ID, InstanceConfig: devConfig},
			{Status: ChallengeInstance_Available, AgentID: dummyAgent3.ID, FlavorID: dummyChallenge3.ID, InstanceConfig: devConfig},
			{Status: ChallengeInstance_Disabled, AgentID: dummyAgent1.ID, FlavorID: dummyChallenge1.ID, InstanceConfig: devConfig},
			{Status: ChallengeInstance_Disabled, AgentID: dummyAgent2.ID, FlavorID: dummyChallenge1.ID, InstanceConfig: devConfig},
			{Status: ChallengeInstance_Disabled, AgentID: dummyAgent3.ID, FlavorID: dummyChallenge2.ID, InstanceConfig: devConfig},
		}
		for _, instance := range instances {
			err = tx.Set("gorm:association_autoupdate", true).Create(instance).Error
			if err != nil {
				return GormToErrcode(err)
			}
		}

		// challenge subscription
		subscription := ChallengeSubscription{
			SeasonChallengeID: dummyChallenge1.SeasonChallenges[0].ID,
			TeamID:            staffTeam.ID,
			BuyerID:           hackSparrow.ID,
			Status:            ChallengeSubscription_Active,
		}
		err = tx.Set("gorm:association_autoupdate", true).Create(&subscription).Error
		if err != nil {
			return GormToErrcode(err)
		}

		// Achievements
		achievements := []*Achievement{
			{
				AuthorID: hackSparrow.ID,
				TeamID:   staffTeam.ID,
				IsGlobal: true,
				Comment:  ":)",
				Type:     Achievement_Staff,
			}, {
				AuthorID: hackSparrow.ID,
				TeamID:   staffTeam.ID,
				Type:     Achievement_Moderator,
			},
		}
		for _, achievement := range achievements {
			err = tx.Create(achievement).Error
			if err != nil {
				return GormToErrcode(err)
			}
		}

		// coupons
		coupons := []*Coupon{
			{Hash: "test-coupon-1", Value: 42, MaxValidationCount: 1, SeasonID: globalSeason.ID},
			{Hash: "test-coupon-2", Value: 42, MaxValidationCount: 1, SeasonID: testSeason.ID},
			{Hash: "test-coupon-3", Value: 42, MaxValidationCount: 0, SeasonID: globalSeason.ID},
			{Hash: "test-coupon-4", Value: 42, MaxValidationCount: 2, SeasonID: globalSeason.ID},
		}
		for _, coupon := range coupons {
			err := tx.
				Set("gorm:association_autoupdate", true).
				Create(coupon).
				Error
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		t.Fatalf("create testing entities: %v", err)
	}
}

func newOfficialChallengeWithFlavor(name string, homepage string, composeBundle string) *ChallengeFlavor {
	return &ChallengeFlavor{
		Challenge: &Challenge{
			Name:     name,
			Author:   "Staff Team",
			Homepage: homepage,
			IsDraft:  false,
		},
		SourceURL:        homepage,
		Version:          "default",
		ComposeBundle:    composeBundle,
		PurchasePrice:    5,
		ValidationReward: 10,
		Driver:           ChallengeFlavor_DockerCompose,
	}
}
