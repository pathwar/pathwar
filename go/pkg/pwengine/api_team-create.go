package pwengine

import (
	"context"
	"fmt"
	"strings"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) TeamCreate(ctx context.Context, in *TeamCreate_Input) (*TeamCreate_Output, error) {
	// validation
	if in == nil {
		return nil, ErrMissingArgument
	}
	if in.SeasonID == 0 {
		return nil, ErrMissingArgument
	}
	if in.OrganizationID == 0 && in.Name == "" {
		return nil, ErrMissingArgument // requires existing organization OR a new name
	}
	if in.OrganizationID != 0 && in.Name != "" {
		return nil, ErrInvalidArgument
	}
	userID, err := userIDFromContext(ctx, e.db)
	if err != nil {
		return nil, fmt.Errorf("get userid from context: %w", err)
	}

	// fetch season
	var season pwdb.Season
	err = e.db.First(&season, in.SeasonID).Error
	if err != nil {
		return nil, ErrInvalidArgument
	}

	// check if season is available for this user
	if season.Status != pwdb.Season_Started {
		return nil, ErrInvalidArgument
	}

	// check if user already has a team in this season
	var existingSeasonMembership pwdb.TeamMember
	err = e.db.
		Model(pwdb.TeamMember{}).
		Joins("JOIN team on team.id = team_member.team_id AND team.season_id = ? AND team.deletion_status = ?", in.SeasonID, pwdb.DeletionStatus_Active).
		Preload("Team").
		Where(pwdb.TeamMember{UserID: userID}).
		First(&existingSeasonMembership).
		Error
	switch {
	case err == nil:
		return nil, ErrInvalidArgument // user already has a team for this season
	case pwdb.IsRecordNotFoundError(err):
		// everything is okay!
	default:
		return nil, err // 500
	}

	if in.OrganizationID == 0 && in.Name != "" {
		in.Name = normalizeName(in.Name)
		if isReservedName(in.Name) {
			return nil, ErrInvalidArgument
		}

		// check for existing organization with that name
		var count int
		err = e.db.Model(pwdb.Organization{}).Where(pwdb.Organization{Name: in.Name}).Count(&count).Error
		if err != nil {
			return nil, err
		}
		if count != 0 {
			return nil, ErrInvalidArgument
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

		err = e.db.Create(&organization).Error
		if err != nil {
			return nil, err
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
	err = e.db.Model(pwdb.Team{}).Where(existingTeam).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count != 0 {
		return nil, ErrInvalidArgument
	}

	// load organization
	var organization pwdb.Organization
	err = e.db.Preload("Members").First(&organization, in.OrganizationID).Error
	if err != nil {
		return nil, ErrInvalidArgument
	}
	found := false
	for _, member := range organization.Members {
		if member.UserID == userID {
			found = true
			break
		}
	}
	if !found {
		return nil, ErrInvalidArgument
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
	err = e.db.Create(&team).Error
	if err != nil {
		return nil, err
	}

	var preloadedTeam pwdb.Team
	err = e.db.
		Preload("Members").
		Preload("Organization").
		Preload("Season").
		First(&preloadedTeam, team.ID).
		Error
	if err != nil {
		return nil, err
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
