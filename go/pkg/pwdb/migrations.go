package pwdb

import (
	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
	"pathwar.land/v2/go/pkg/errcode"
)

func migrate(db *gorm.DB, sfn *snowflake.Node, opts Opts) error {
	migrateOpts := gormigrate.DefaultOptions
	migrateOpts.UseTransaction = true
	m := gormigrate.New(db, migrateOpts, []*gormigrate.Migration{})

	// only called on fresh database
	m.InitSchema(func(tx *gorm.DB) error {
		tx.Set("gorm:table_options", "charset=utf8mb4")
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

		err = createFirstEntities(tx)
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

func createFirstEntities(tx *gorm.DB) error {
	// FIXME: replace those direct DB inserts by API calls to admin endpoints
	// default season
	solo := &Season{
		Name:       "Solo Mode",
		Status:     Season_Started,
		Visibility: Season_Public,
		IsDefault:  true,
	}
	err := tx.Create(solo).Error
	if err != nil {
		return GormToErrcode(err)
	}

	// staff org, team and members
	staffOrg := &Organization{
		Name:           "Staff",
		DeletionStatus: DeletionStatus_Active,
	}
	staffTeam := &Team{
		IsDefault:      true,
		Season:         solo,
		Organization:   staffOrg,
		DeletionStatus: DeletionStatus_Active,
	}
	hackSparrow := &User{
		Username:                "Hack Sparrow",
		OAuthSubject:            "Hack Sparrow",
		OrganizationMemberships: []*OrganizationMember{{Organization: staffOrg}},
		TeamMemberships:         []*TeamMember{{Team: staffTeam}},
		DeletionStatus:          DeletionStatus_Active,
	}
	err = tx.Set("gorm:association_autoupdate", true).Create(hackSparrow).Error
	if err != nil {
		return GormToErrcode(err)
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
	challengeDebug.addSeasonChallengeByID(solo.ID)

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
	helloworld.addSeasonChallengeByID(solo.ID)

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
	trainingHTTP.addSeasonChallengeByID(solo.ID)

	for _, flavor := range []*ChallengeFlavor{
		challengeDebug, helloworld, trainingSQLI, trainingHTTP,
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

	return nil
}
