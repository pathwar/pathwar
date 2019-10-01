package pwdb

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func migrate(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{})

	// only called on fresh database
	m.InitSchema(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(All()...).Error; err != nil {
			tx.Rollback()
			return err
		}
		for _, fk := range ForeignKeys() {
			e := ByName(fk[0])
			if err := tx.Model(e).AddForeignKey(fk[1], fk[2], "RESTRICT", "RESTRICT").Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		for _, entity := range firstEntities() {
			if err := tx.Create(entity).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		return nil
	})

	// FIXME: add new migrations here...

	if err := m.Migrate(); err != nil {
		return err
	}

	// anyway, call db.automigrate
	if err := db.AutoMigrate(All()...).Error; err != nil {
		return err
	}

	return nil
}

func firstEntities() []interface{} {
	solo := &Tournament{
		// ID:         "solo-tournament",
		Name:       "Solo Mode",
		Status:     Tournament_Started,
		Visibility: Tournament_Public,
		IsDefault:  true,
	}
	testTournament := &Tournament{
		Name:       "Test Tournament",
		Status:     Tournament_Started,
		Visibility: Tournament_Public,
	}
	m1ch3l := &User{
		ID:       "m1ch3l",
		Username: "m1ch3l",
		// State: special
	}
	staff := &Team{
		Name: "Staff",
	}
	soloStaff := &TournamentTeam{
		IsDefault:  true,
		Tournament: solo,
		Team:       staff,
	}
	localhost := &Hypervisor{
		Name:    "localhost",
		Address: "127.0.0.1",
		Status:  Hypervisor_Active, // only useful during dev
	}
	helloWorld := &Level{
		Name:     "Hello World (test)",
		IsDraft:  false,
		Author:   "m1ch3l",
		Homepage: "https://github.com/pathwar/pathwar/tree/master/level/example/hello-world",
	}
	helloWorldLatest := &LevelVersion{
		Level:     helloWorld,
		SourceURL: "https://github.com/pathwar/pathwar/tree/master/level/example/hello-world",
		IsLatest:  true,
		IsDraft:   false,
		Changelog: "Lorem Ipsum",
		Version:   "latest",
		Driver:    LevelVersion_DockerCompose,
	}

	return []interface{}{
		solo,
		testTournament,
		m1ch3l,
		staff,
		soloStaff,
		localhost,
		helloWorld,
		helloWorldLatest,
	}
}
