package pwapi

import (
	"context"
	"strings"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) TeamCreate(ctx context.Context, in *TeamCreate_Input) (*TeamCreate_Output, error) {
	if in == nil || in.SeasonID == 0 || (in.OrganizationID == 0 && in.Name == "") {
		return nil, errcode.ErrMissingInput
	}
	if in.OrganizationID != 0 && in.Name != "" {
		return nil, errcode.ErrInvalidInput
	}
	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated.Wrap(err)
	}

	// fetch season
	var season pwdb.Season
	err = svc.db.First(&season, in.SeasonID).Error
	if err != nil {
		return nil, errcode.ErrGetSeason.Wrap(err)
	}

	// check if season is available for this user
	if season.Status != pwdb.Season_Started {
		return nil, errcode.ErrSeasonDenied
	}

	// check if user already has a team in this season
	var existingSeasonMembership pwdb.TeamMember
	err = svc.db.
		Model(pwdb.TeamMember{}).
		Joins("JOIN team on team.id = team_member.team_id AND team.season_id = ? AND team.deletion_status = ?", in.SeasonID, pwdb.DeletionStatus_Active).
		Preload("Team").
		Where(pwdb.TeamMember{UserID: userID}).
		First(&existingSeasonMembership).
		Error
	if !pwdb.IsRecordNotFoundError(err) {
		return nil, errcode.ErrAlreadyHasTeamForSeason.Wrap(err)
	}

	if in.OrganizationID == 0 && in.Name != "" {
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
				{UserID: userID},
			},
			DeletionStatus: pwdb.DeletionStatus_Active,
			// GravatarURL
			// Locale
		}
		err = svc.db.Create(&organization).Error
		if err != nil {
			return nil, errcode.ErrCreateOrganization.Wrap(err)
		}

		in.OrganizationID = organization.ID
	}

	// check if there is already a team for this organization and season couple
	var count int
	existingTeam := pwdb.Team{
		SeasonID:       in.SeasonID,
		OrganizationID: in.OrganizationID,
		DeletionStatus: pwdb.DeletionStatus_Active,
	}
	err = svc.db.Model(pwdb.Team{}).Where(existingTeam).Count(&count).Error
	if err != nil || count != 0 {
		return nil, errcode.ErrOrganizationAlreadyHasTeamForSeason.Wrap(err)
	}

	// load organization
	var organization pwdb.Organization
	err = svc.db.Preload("Members").First(&organization, in.OrganizationID).Error
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
		SeasonID:       in.SeasonID,
		OrganizationID: in.OrganizationID,
		DeletionStatus: pwdb.DeletionStatus_Active,
		Members: []*pwdb.TeamMember{
			{UserID: userID},
		},
	}

	// save new team object in DB
	err = svc.db.Create(&team).Error
	if err != nil {
		return nil, errcode.ErrCreateTeam.Wrap(err)
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
