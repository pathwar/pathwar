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
		if err := tx.AutoMigrate(All()...).Error; err != nil {
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

		for _, entity := range firstEntities(sfn) {
			if err := tx.Create(entity).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("create first entities: %w", err)
			}
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

func firstEntities(sfn *snowflake.Node) []interface{} {
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
	m1ch3l := &User{
		Username:     "m1ch3l",
		OAuthSubject: "m1ch3l",
		// State: special
	}
	staff := &Organization{
		Name: "Staff",
	}
	soloStaff := &Team{
		IsDefault:    true,
		Season:       solo,
		Organization: staff,
	}
	localhost := &Hypervisor{
		Name:    "localhost",
		Address: "127.0.0.1",
		Status:  Hypervisor_Active, // only useful during dev
	}
	helloWorld := &Challenge{
		Name:     "Hello World (test)",
		IsDraft:  false,
		Author:   "m1ch3l",
		Homepage: "https://github.com/pathwar/pathwar/tree/master/challenge/example/hello-world",
	}
	helloWorldLatest := &ChallengeFlavor{
		Challenge: helloWorld,
		SourceURL: "https://github.com/pathwar/pathwar/tree/master/challenge/example/hello-world",
		IsLatest:  true,
		IsDraft:   false,
		Changelog: "Lorem Ipsum",
		Version:   "latest",
		Driver:    ChallengeFlavor_DockerCompose,
	}

	return []interface{}{
		solo,
		testSeason,
		m1ch3l,
		staff,
		soloStaff,
		localhost,
		helloWorld,
		helloWorldLatest,
	}
}
