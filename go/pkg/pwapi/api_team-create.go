package pwapi

import (
	"context"
	"strings"

	"github.com/jinzhu/gorm"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) TeamCreate(ctx context.Context, in *TeamCreate_Input) (*TeamCreate_Output, error) {
	if in == nil || in.SeasonID == "" || (in.OrganizationID == "" && in.Name == "") {
		return nil, errcode.ErrMissingInput
	}
	if in.OrganizationID != "" && in.Name != "" {
		return nil, errcode.ErrInvalidInput
	}
	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated.Wrap(err)
	}

	var organizationID int64
	if in.OrganizationID != "" {
		organizationID, err = pwdb.GetIDBySlugAndKind(svc.db, in.OrganizationID, "organization")
		if err != nil {
			return nil, err
		}
	}

	seasonID, err := pwdb.GetIDBySlugAndKind(svc.db, in.SeasonID, "season")
	if err != nil {
		return nil, err
	}

	// fetch season
	var season pwdb.Season
	err = svc.db.First(&season, seasonID).Error
	if err != nil {
		return nil, errcode.ErrGetSeason.Wrap(err)
	}

	// check if season is available for this user
	if season.Status != pwdb.Season_Started {
		return nil, errcode.ErrSeasonDenied
	}
	if season.Visibility == pwdb.Season_Private {
		return nil, errcode.ErrSeasonDenied
	}
	if season.Subscription == pwdb.Season_Close {
		return nil, errcode.ErrSeasonDenied
	}

	// check if user already has a team in this season
	var seasonMemberShipCount int
	err = svc.db.
		Model(pwdb.TeamMember{}).
		Joins("JOIN team on team.id = team_member.team_id").
		Where(pwdb.TeamMember{UserID: userID}).
		Where(&pwdb.Team{
			SeasonID:       seasonID,
			DeletionStatus: pwdb.DeletionStatus_Active,
		}).
		Count(&seasonMemberShipCount).
		Error
	if err != nil || seasonMemberShipCount != 0 {
		return nil, errcode.ErrAlreadyHasTeamForSeason.Wrap(err)
	}

	if in.OrganizationID == "" && in.Name != "" {
		in.Name = normalizeName(in.Name)
		if isReservedName(in.Name) {
			return nil, errcode.ErrReservedName
		}

		// check for existing organization with that name
		var count int
		err = svc.db.Model(pwdb.Organization{}).Where(pwdb.Organization{Name: in.Name}).Count(&count).Error
		if err != nil || count != 0 {
			return nil, errcode.ErrCheckOrganizationUniqueName.Wrap(err)
		}

		// create new organization
		organization := pwdb.Organization{
			Name: in.Name,
			Members: []*pwdb.OrganizationMember{
				{
					UserID: userID,
					Role:   pwdb.OrganizationMember_Owner,
				},
			},
			DeletionStatus: pwdb.DeletionStatus_Active,
			// GravatarURL
			// Locale
		}
		err = svc.db.Create(&organization).Error
		if err != nil {
			return nil, errcode.ErrCreateOrganization.Wrap(err)
		}

		organizationID = organization.ID
	}

	// check that user is member of the organization
	var memberCount int
	err = svc.db.
		Model(pwdb.OrganizationMember{}).
		Where(pwdb.OrganizationMember{
			UserID:         userID,
			OrganizationID: organizationID,
		}).
		Count(&memberCount).
		Error
	if err != nil || memberCount == 0 {
		return nil, errcode.ErrUserNotInOrganization.Wrap(err)
	}

	// check if there is already a team for this organization and season couple
	var count int
	existingTeam := pwdb.Team{
		SeasonID:       seasonID,
		OrganizationID: organizationID,
		DeletionStatus: pwdb.DeletionStatus_Active,
	}
	err = svc.db.Model(pwdb.Team{}).Where(existingTeam).Count(&count).Error
	if err != nil || count != 0 {
		return nil, errcode.ErrOrganizationAlreadyHasTeamForSeason.Wrap(err)
	}

	// load organization
	var organization pwdb.Organization
	err = svc.db.Preload("Members").First(&organization, organizationID).Error
	if err != nil {
		return nil, errcode.ErrGetOrganization.Wrap(err)
	}

	// check if organization is in global season
	if organization.GlobalSeason {
		return nil, errcode.ErrCannotCreateTeamForGlobalOrganization
	}

	// check if user belongs to the organization
	found := false
	for _, member := range organization.Members {
		if member.UserID == userID {
			found = true
			break
		}
	}
	if !found {
		return nil, errcode.ErrUserNotInOrganization
	}

	// construct new team object
	team := pwdb.Team{
		SeasonID:       seasonID,
		OrganizationID: organizationID,
		DeletionStatus: pwdb.DeletionStatus_Active,
		Score:          0,
		GoldMedals:     0,
		SilverMedals:   0,
		BronzeMedals:   0,
		NbAchievements: 0,
		Members: []*pwdb.TeamMember{
			{
				UserID: userID,
				Role:   pwdb.TeamMember_Owner,
			},
		},
	}

	// save new team object in DB
	err = svc.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&team).Error
		if err != nil {
			return errcode.ErrCreateTeam.Wrap(err)
		}
		activity := pwdb.Activity{
			Kind:           pwdb.Activity_TeamCreation,
			AuthorID:       userID,
			TeamID:         team.ID,
			OrganizationID: team.OrganizationID,
			SeasonID:       seasonID,
		}
		return tx.Create(&activity).Error
	})
	if err != nil {
		return nil, err
	}

	// reload with associations
	var preloadedTeam pwdb.Team
	err = svc.db.
		Preload("Members").
		Preload("Organization").
		Preload("Season").
		First(&preloadedTeam, team.ID).
		Error
	if err != nil {
		return nil, errcode.ErrGetTeam.Wrap(err)
	}

	ret := TeamCreate_Output{Team: &preloadedTeam}
	return &ret, nil
}

func normalizeName(name string) string {
	return strings.TrimSpace(name)
}

func isReservedName(name string) bool {
	name = strings.ToLower(name)

	reserved := []string{
		"admin",
		"pathwar",
		"root",
		"staff",
		// FIXME: more complete list of blacklist
	}

	for _, try := range reserved {
		if try == name {
			return true
		}
	}

	return false
}
