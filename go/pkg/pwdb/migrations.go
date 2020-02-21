package pwdb

import (
	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
	"pathwar.land/go/pkg/errcode"
)

func migrate(db *gorm.DB, sfn *snowflake.Node, opts Opts) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{})

	// only called on fresh database
	m.InitSchema(func(tx *gorm.DB) error {
		err := tx.AutoMigrate(All()...).Error
		if err != nil {
			tx.Rollback()
			return errcode.ErrDBAutoMigrate.Wrap(err)
		}

		if !opts.skipFK {
			for _, fk := range ForeignKeys() {
				e := ByName(fk[0])
				if err := tx.Model(e).AddForeignKey(fk[1], fk[2], "RESTRICT", "RESTRICT").Error; err != nil {
					tx.Rollback()
					return errcode.ErrDBAddForeignKey.Wrap(err)
				}
			}
		}

		err = createFirstEntities(tx, sfn, opts)
		if err != nil {
			return GormToErrcode(err)
		}

		return nil
	})

	// FIXME: add new migrations here...

	err := m.Migrate()
	if err != nil {
		return errcode.ErrDBRunMigrations.Wrap(err)
	}

	// anyway, call db.automigrate
	err = db.AutoMigrate(All()...).Error
	if err != nil {
		return errcode.ErrDBAutoMigrate.Wrap(err)
	}

	return nil
}

func createFirstEntities(tx *gorm.DB, sfn *snowflake.Node, opts Opts) error {
	//
	// seasons
	//
	solo := &Season{
		// ID:         "solo-season",
		Name:       "Solo Mode",
		Status:     Season_Started,
		Visibility: Season_Public,
		IsDefault:  true,
	}
	testSeason := &Season{
		Name:       "Test Season",
		Status:     Season_Started,
		Visibility: Season_Public,
	}
	for _, season := range []*Season{solo, testSeason} {
		if err := tx.Create(season).Error; err != nil {
			return GormToErrcode(err)
		}
	}

	//
	// staff team & org
	//

	staffOrg := &Organization{
		Name:           "Staff",
		DeletionStatus: DeletionStatus_Active,
		// GravatarURL: staff
	}
	staffTeam := &Team{
		IsDefault:      true,
		Season:         solo,
		Organization:   staffOrg,
		DeletionStatus: DeletionStatus_Active,
		// GravatarURL: staff
	}
	hackSparrow := &User{
		Username:                "Hack Sparrow",
		OAuthSubject:            "Hack Sparrow",
		OrganizationMemberships: []*OrganizationMember{{Organization: staffOrg}},
		TeamMemberships:         []*TeamMember{{Team: staffTeam}},
		DeletionStatus:          DeletionStatus_Active,
		// State: special
		// GravatarURL: m1ch3l
	}
	err := tx.
		Set("gorm:association_autoupdate", true).
		Create(hackSparrow).
		Error
	if err != nil {
		return GormToErrcode(err)
	}

	//
	// agents
	//

	localhost := &Agent{
		Name:     "localhost",
		Hostname: "localhost",
		Status:   Agent_Active, // only useful during dev
	}
	localhost2 := &Agent{
		Name:     "localhost-2",
		Hostname: "localhost",
		Status:   Agent_Active,
	}
	localhost3 := &Agent{
		Name:     "localhost-3",
		Hostname: "localhost",
		Status:   Agent_Inactive,
	}
	for _, agent := range []*Agent{localhost, localhost2, localhost3} {
		err = tx.Create(agent).Error
		if err != nil {
			return GormToErrcode(err)
		}
	}

	//
	// challenges
	//
	bundle := `version: "3.7"
networks: {}
volumes: {}
services:
    front:
        image: pathwar/helloworld@sha256:00247fcdcad9336f9cbfcde74a56650b6ffd7c27037957a1e8048d02eb7bdbe3
        ports:
            - "80"
        labels:
            land.pathwar.compose.challenge-name: helloworld
            land.pathwar.compose.challenge-version: 1.0.0
            land.pathwar.compose.origin: was-built
            land.pathwar.compose.service-name: front
`
	helloworld := newOfficialChallengeWithFlavor("Hello World", "https://github.com/pathwar/pathwar/tree/master/challenges/web/helloworld", bundle)
	helloworld.addSeasonChallengeByID(solo.ID)
	helloworld.addSeasonChallengeByID(testSeason.ID)

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
	trainingSQLI.addSeasonChallengeByID(solo.ID)

	nopBundle := ``

	trainingHTTP := newOfficialChallengeWithFlavor("Training HTTP", "https://github.com/pathwar/pathwar/tree/master/challenges/web/training-http", nopBundle)
	trainingHTTP.addSeasonChallengeByID(solo.ID)

	trainingInclude := newOfficialChallengeWithFlavor("Training Include", "https://github.com/pathwar/pathwar/tree/master/challenges/web/training-include", nopBundle)
	trainingInclude.addSeasonChallengeByID(solo.ID)

	trainingBrute := newOfficialChallengeWithFlavor("Training Brute", "https://github.com/pathwar/pathwar/tree/master/challenges/web/training-brute", nopBundle)
	trainingBrute.addSeasonChallengeByID(solo.ID)

	captchaLuigi := newOfficialChallengeWithFlavor("Captcha Luigi", "https://github.com/pathwar/pathwar/tree/master/challenges/web/captcha-luigi", nopBundle)
	captchaLuigi.addSeasonChallengeByID(testSeason.ID)

	captchaMario := newOfficialChallengeWithFlavor("Captcha Mario", "https://github.com/pathwar/pathwar/tree/master/challenges/web/captcha-mario", nopBundle)
	captchaMario.addSeasonChallengeByID(testSeason.ID)

	uploadHi := newOfficialChallengeWithFlavor("Upload HI", "https://github.com/pathwar/pathwar/tree/master/challenges/web/upload-hi", nopBundle)
	uploadHi.addSeasonChallengeByID(testSeason.ID)

	imageboard := newOfficialChallengeWithFlavor("Image Board", "https://github.com/pathwar/pathwar/tree/master/challenges/web/imageboard", nopBundle)
	imageboard.addSeasonChallengeByID(testSeason.ID)

	for _, flavor := range []*ChallengeFlavor{
		helloworld, trainingSQLI, trainingHTTP, trainingInclude, trainingBrute,
		captchaLuigi, captchaMario, uploadHi, imageboard,
	} {
		err := tx.
			Set("gorm:association_autoupdate", true).
			Create(flavor).
			Error
		if err != nil {
			return GormToErrcode(err)
		}

		// FIXME: should not be necessary, should be done automatically thanks to association_autoupdate
		for _, seasonChallenge := range flavor.SeasonChallenges {
			seasonChallenge.FlavorID = flavor.ID
			err := tx.
				Set("gorm:association_autoupdate", true).
				Create(seasonChallenge).
				Error
			if err != nil {
				return GormToErrcode(err)
			}
		}
	}

	//// Challenge Instances
	devConfig := []byte(`{"passphrases": ["a", "b", "c", "d"]}`)
	instances := []*ChallengeInstance{
		{Status: ChallengeInstance_Available, AgentID: localhost.ID, FlavorID: trainingSQLI.ID, InstanceConfig: devConfig},
		{Status: ChallengeInstance_Available, AgentID: localhost2.ID, FlavorID: trainingSQLI.ID, InstanceConfig: devConfig},
		{Status: ChallengeInstance_Available, AgentID: localhost3.ID, FlavorID: helloworld.ID, InstanceConfig: devConfig},
		{Status: ChallengeInstance_Disabled, AgentID: localhost.ID, FlavorID: trainingSQLI.ID, InstanceConfig: devConfig},
		{Status: ChallengeInstance_Disabled, AgentID: localhost2.ID, FlavorID: trainingSQLI.ID, InstanceConfig: devConfig},
		{Status: ChallengeInstance_Disabled, AgentID: localhost3.ID, FlavorID: helloworld.ID, InstanceConfig: devConfig},
	}
	for _, instance := range instances {
		err := tx.Set("gorm:association_autoupdate", true).
			Create(instance).
			Error
		if err != nil {
			return GormToErrcode(err)
		}
	}

	// challenge subscription
	subscription := ChallengeSubscription{
		SeasonChallengeID: trainingSQLI.SeasonChallenges[0].ID,
		TeamID:            staffTeam.ID,
		BuyerID:           hackSparrow.ID,
		Status:            ChallengeSubscription_Active,
	}
	err = tx.Set("gorm:association_autoupdate", true).
		Create(&subscription).
		Error
	if err != nil {
		return GormToErrcode(err)
	}

	//
	// Achievements
	//

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

	//
	// coupons
	//
	coupons := []*Coupon{
		{Hash: "test-coupon-1", Value: 42, MaxValidationCount: 1, SeasonID: solo.ID},
		{Hash: "test-coupon-2", Value: 42, MaxValidationCount: 1, SeasonID: testSeason.ID},
		{Hash: "test-coupon-3", Value: 42, MaxValidationCount: 0, SeasonID: solo.ID},
		{Hash: "test-coupon-4", Value: 42, MaxValidationCount: 2, SeasonID: solo.ID},
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
}
