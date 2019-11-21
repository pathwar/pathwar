package pwdb

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func migrate(db *gorm.DB, sfn *snowflake.Node, opts Opts) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{})

	// only called on fresh database
	m.InitSchema(func(tx *gorm.DB) error {
		err := tx.AutoMigrate(All()...).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("automigrate: %w", err)
		}

		if !opts.skipFK {
			for _, fk := range ForeignKeys() {
				e := ByName(fk[0])
				if err := tx.Model(e).AddForeignKey(fk[1], fk[2], "RESTRICT", "RESTRICT").Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("addforeignkey %q %q: %w", fk[1], fk[2], err)
				}
			}
		}

		err = createFirstEntities(tx, sfn)
		if err != nil {
			return fmt.Errorf("create first entities: %w", err)
		}

		return nil
	})

	// FIXME: add new migrations here...

	if err := m.Migrate(); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	// anyway, call db.automigrate
	if err := db.AutoMigrate(All()...).Error; err != nil {
		return fmt.Errorf("automigrate: %w", err)
	}

	return nil
}

func createFirstEntities(tx *gorm.DB, sfn *snowflake.Node) error {
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
			return err
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
		return err
	}

	//
	// hypervisors
	//

	localhost := &Hypervisor{
		Name:    "default",
		Address: "default-hypervisor.pathwar.land",
		Status:  Hypervisor_Active, // only useful during dev
	}
	err = tx.Create(localhost).Error
	if err != nil {
		return err
	}

	//
	// challenges
	//

	helloworld := newOfficialChallengeWithFlavor("Hello World", "https://github.com/pathwar/pathwar/tree/master/challenges/web/helloworld")
	helloworld.addSeasonChallengeByID(solo.ID)
	helloworld.addSeasonChallengeByID(testSeason.ID)

	trainingHTTP := newOfficialChallengeWithFlavor("Training HTTP", "https://github.com/pathwar/pathwar/tree/master/challenges/web/training-http")
	trainingHTTP.addSeasonChallengeByID(solo.ID)

	trainingSQLI := newOfficialChallengeWithFlavor("Training SQLI", "https://github.com/pathwar/pathwar/tree/master/challenges/web/training-sqli")
	trainingSQLI.addSeasonChallengeByID(solo.ID)

	trainingInclude := newOfficialChallengeWithFlavor("Training Include", "https://github.com/pathwar/pathwar/tree/master/challenges/web/training-include")
	trainingInclude.addSeasonChallengeByID(solo.ID)

	trainingBrute := newOfficialChallengeWithFlavor("Training Brute", "https://github.com/pathwar/pathwar/tree/master/challenges/web/training-brute")
	trainingBrute.addSeasonChallengeByID(solo.ID)

	captchaLuigi := newOfficialChallengeWithFlavor("Captcha Luigi", "https://github.com/pathwar/pathwar/tree/master/challenges/web/captcha-luigi")
	captchaLuigi.addSeasonChallengeByID(testSeason.ID)

	captchaMario := newOfficialChallengeWithFlavor("Captcha Mario", "https://github.com/pathwar/pathwar/tree/master/challenges/web/captcha-mario")
	captchaMario.addSeasonChallengeByID(testSeason.ID)

	uploadHi := newOfficialChallengeWithFlavor("Upload HI", "https://github.com/pathwar/pathwar/tree/master/challenges/web/upload-hi")
	uploadHi.addSeasonChallengeByID(testSeason.ID)

	imageboard := newOfficialChallengeWithFlavor("Image Board", "https://github.com/pathwar/pathwar/tree/master/challenges/web/imageboard")
	imageboard.addSeasonChallengeByID(testSeason.ID)

	for _, flavor := range []*ChallengeFlavor{
		helloworld, trainingHTTP, trainingSQLI, trainingInclude, trainingBrute,
		captchaLuigi, captchaMario, uploadHi, imageboard,
	} {
		err := tx.
			Set("gorm:association_autoupdate", true).
			Create(flavor).
			Error
		if err != nil {
			return err
		}

		// FIXME: should not be necessary, should be done automatically thanks to association_autoupdate
		for _, seasonChallenge := range flavor.SeasonChallenges {
			seasonChallenge.FlavorID = flavor.ID
			err := tx.
				Set("gorm:association_autoupdate", true).
				Create(seasonChallenge).
				Error
			if err != nil {
				return err
			}
		}
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
			return err
		}
	}

	return nil
}
