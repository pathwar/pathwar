package pwdb

import (
	"gorm.io/gorm"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(All()...)
	if err != nil {
		return errcode.ErrDBAutoMigrate.Wrap(err)
	}

	// create first entities only on fresh database
	err = db.First(&Season{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = createFirstEntities(db)
			if err != nil {
				return GormToErrcode(err)
			}
		} else {
			return GormToErrcode(err)
		}
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
		IsGlobal:   true,
	}
	err := tx.Create(globalSeason).Error
	if err != nil {
		return GormToErrcode(err)
	}
	testingSeason := &Season{
		Name:       "Testing",
		Status:     Season_Started,
		Visibility: Season_Private,
		IsTesting:  true,
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
	err = tx.Create(staffOrg).Error
	if err != nil {
		return GormToErrcode(err)
	}
	staffTeamGlobal := &Team{
		IsGlobal:       true,
		Season:         globalSeason,
		Organization:   staffOrg,
		DeletionStatus: DeletionStatus_Active,
		Slug:           "staff@global",
	}
	staffTeamTesting := &Team{
		Season:         testingSeason,
		Organization:   staffOrg,
		DeletionStatus: DeletionStatus_Active,
		Slug:           "staff@testing",
	}
	hackSparrow := &User{
		Username:                "Hack Sparrow",
		OAuthSubject:            "Hack Sparrow",
		OrganizationMemberships: []*OrganizationMember{{Organization: staffOrg}},
		TeamMemberships:         []*TeamMember{{Team: staffTeamGlobal}, {Team: staffTeamTesting}},
		DeletionStatus:          DeletionStatus_Active,
		Slug:                    "hack-sparrow",
	}
	err = tx.Set("gorm:association_autoupdate", true).Create(hackSparrow).Error
	if err != nil {
		return GormToErrcode(err)
	}
	return nil
}
