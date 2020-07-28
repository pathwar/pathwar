package pwdb

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func migrate(db *gorm.DB, opts Opts) error {
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
	globalSeason := &Season{
		Name:       "Global",
		Status:     Season_Started,
		Visibility: Season_Public,
		IsDefault:  true,
	}
	err := tx.Create(globalSeason).Error
	if err != nil {
		return GormToErrcode(err)
	}
	testingSeason := &Season{
		Name:       "Testing",
		Status:     Season_Started,
		Visibility: Season_Private,
		IsDefault:  false,
	}
	err = tx.Create(testingSeason).Error
	if err != nil {
		return GormToErrcode(err)
	}

	// staff org, team and members
	staffOrg := &Organization{
		Name:           "Staff",
		DeletionStatus: DeletionStatus_Active,
	}
	staffTeamGlobal := &Team{
		IsDefault:      true,
		Season:         globalSeason,
		Organization:   staffOrg,
		DeletionStatus: DeletionStatus_Active,
	}
	staffTeamTesting := &Team{
		IsDefault:      true,
		Season:         testingSeason,
		Organization:   staffOrg,
		DeletionStatus: DeletionStatus_Active,
	}
	hackSparrow := &User{
		Username:                "Hack Sparrow",
		OAuthSubject:            "Hack Sparrow",
		OrganizationMemberships: []*OrganizationMember{{Organization: staffOrg}},
		TeamMemberships:         []*TeamMember{{Team: staffTeamGlobal}, {Team: staffTeamTesting}},
		DeletionStatus:          DeletionStatus_Active,
	}
	err = tx.Set("gorm:association_autoupdate", true).Create(hackSparrow).Error
	if err != nil {
		return GormToErrcode(err)
	}
	return nil
}
