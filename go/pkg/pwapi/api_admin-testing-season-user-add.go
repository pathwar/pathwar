package pwapi

import (
	"context"

	"github.com/jinzhu/gorm"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminTestingSeasonUserAdd(ctx context.Context, in *AdminTestingSeasonUserAdd_Input) (*AdminTestingSeasonUserAdd_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	if in == nil || in.UserID == "" {
		return nil, errcode.ErrMissingInput
	}

	userID, err := pwdb.GetIDBySlugAndKind(svc.db, in.UserID, "user")
	if err != nil {
		return nil, err
	}
	var user pwdb.User
	err = svc.db.Where(&pwdb.User{ID: userID}).First(&user).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	testingSeasonID, err := pwdb.GetIDBySlugAndKind(svc.db, "testing", "season")
	if err != nil {
		return nil, err
	}
	var testingSeason pwdb.Season
	err = svc.db.Where(&pwdb.Season{ID: testingSeasonID}).First(&testingSeason).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	// check if user already has a team in this season
	var seasonMemberShipCount int
	err = svc.db.
		Model(pwdb.TeamMember{}).
		Joins("JOIN team on team.id = team_member.team_id").
		Where(pwdb.TeamMember{UserID: userID}).
		Where(&pwdb.Team{
			SeasonID:       testingSeasonID,
			DeletionStatus: pwdb.DeletionStatus_Active}).
		Count(&seasonMemberShipCount).
		Error
	if err != nil || seasonMemberShipCount != 0 {
		return nil, errcode.ErrAlreadyHasTeamForSeason.Wrap(err)
	}

	// retrieve user solo global organization
	var globalOrganizationMember pwdb.OrganizationMember
	err = svc.db.
		Model(pwdb.OrganizationMember{}).
		Joins("JOIN organization on organization.id = organization_member.organization_id").
		Where(&pwdb.OrganizationMember{UserID: userID}).
		Where(&pwdb.Organization{
			GlobalSeason: true,
		}).
		First(&globalOrganizationMember).
		Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	var teamMember pwdb.TeamMember
	err = svc.db.Transaction(func(tx *gorm.DB) error {
		// create team and team member
		team := pwdb.Team{
			Season:         &testingSeason,
			IsGlobal:       true,
			OrganizationID: globalOrganizationMember.OrganizationID,
			DeletionStatus: pwdb.DeletionStatus_Active,
		}
		err = tx.Create(&team).Error
		if err != nil {
			return err
		}
		teamMember = pwdb.TeamMember{
			UserID: userID,
			Team:   &team,
			Role:   pwdb.TeamMember_Owner,
		}
		return tx.Create(&teamMember).Error
	})
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	return &AdminTestingSeasonUserAdd_Output{TeamMember: &teamMember}, nil
}
