package pwdb

import (
	"testing"

	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

func TestingSqliteDB(t *testing.T, logger *zap.Logger) *gorm.DB {
	t.Helper()

	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("init in-memory sqlite server: %v", err)
	}

	sfn, err := snowflake.NewNode(1)
	if err != nil {
		t.Fatalf("init snowflake generator: %v", err)
	}

	opts := Opts{
		Logger: logger,
		skipFK: true, // required for sqlite :(
	}

	db, err = Configure(db, sfn, opts)
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
		soloSeason := Season{Name: "Solo Mode"}
		err := tx.Where(&soloSeason).First(&soloSeason).Error
		if err != nil {
			return GormToErrcode(err)
		}
		staffOrg := Organization{Name: "Staff"}
		err = tx.Where(&staffOrg).First(&staffOrg).Error
		if err != nil {
			return GormToErrcode(err)
		}
		staffTeam := Team{SeasonID: soloSeason.ID, OrganizationID: staffOrg.ID}
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
			Name:       "Test Season",
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

		dummyChallenge1 := newOfficialChallengeWithFlavor("dummy challenge 1", "https://...", "")
		dummyChallenge1.addSeasonChallengeByID(soloSeason.ID)
		dummyChallenge1.addSeasonChallengeByID(testSeason.ID)
		dummyChallenge2 := newOfficialChallengeWithFlavor("dummy challenge 1", "https://...", "")
		dummyChallenge2.addSeasonChallengeByID(soloSeason.ID)
		dummyChallenge2.addSeasonChallengeByID(testSeason.ID)
		dummyChallenge3 := newOfficialChallengeWithFlavor("dummy challenge 1", "https://...", "")
		dummyChallenge3.addSeasonChallengeByID(soloSeason.ID)
		dummyChallenge3.addSeasonChallengeByID(testSeason.ID)

		for _, flavor := range []*ChallengeFlavor{
			dummyChallenge1, dummyChallenge2, dummyChallenge3,
		} {
			err := tx.Set("gorm:association_autoupdate", true).Create(flavor).Error
			if err != nil {
				return GormToErrcode(err)
			}

			// FIXME: should not be necessary, should be done automatically thanks to association_autoupdate
			for _, seasonChallenge := range flavor.SeasonChallenges {
				seasonChallenge.FlavorID = flavor.ID
				err := tx.Set("gorm:association_autoupdate", true).Create(seasonChallenge).Error
				if err != nil {
					return GormToErrcode(err)
				}
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
			{Hash: "test-coupon-1", Value: 42, MaxValidationCount: 1, SeasonID: soloSeason.ID},
			{Hash: "test-coupon-2", Value: 42, MaxValidationCount: 1, SeasonID: testSeason.ID},
			{Hash: "test-coupon-3", Value: 42, MaxValidationCount: 0, SeasonID: soloSeason.ID},
			{Hash: "test-coupon-4", Value: 42, MaxValidationCount: 2, SeasonID: soloSeason.ID},
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
