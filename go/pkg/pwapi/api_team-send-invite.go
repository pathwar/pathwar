package pwapi

import (
	"context"

	"github.com/jinzhu/gorm"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) TeamSendInvite(ctx context.Context, in *TeamSendInvite_Input) (*TeamSendInvite_Output, error) {
	if in == nil || in.TeamID == "" || in.UserID == "" {
		return nil, errcode.ErrMissingInput
	}

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated.Wrap(err)
	}

	teamID, err := pwdb.GetIDBySlugAndKind(svc.db, in.TeamID, "team")
	if err != nil {
		return nil, err
	}

	invitedUserID, err := pwdb.GetIDBySlugAndKind(svc.db, in.UserID, "user")
	if err != nil {
		return nil, err
	}

	// check team status
	var team pwdb.Team
	err = svc.db.
		Where(pwdb.Team{
			ID:             teamID,
			DeletionStatus: pwdb.DeletionStatus_Active,
		}).
		Preload("Season").
		First(&team).
		Error
	if err != nil {
		return nil, errcode.ErrTeamDoesNotExist.Wrap(err)
	}

	// check that the user is owner of the team
	var teamOwner pwdb.TeamMember
	err = svc.db.
		Where(pwdb.TeamMember{
			UserID: userID,
			TeamID: teamID,
			Role:   pwdb.TeamMember_Owner,
		}).
		First(&teamOwner).
		Error
	if err != nil {
		return nil, errcode.ErrNotTeamOwner.Wrap(err)
	}

	// check if invited user already has a team in this season
	var seasonMemberShipCount int
	err = svc.db.
		Model(pwdb.TeamMember{}).
		Joins("JOIN team on team.id = team_member.team_id").
		Where(pwdb.TeamMember{UserID: invitedUserID}).
		Where(&pwdb.Team{
			SeasonID:       team.SeasonID,
			DeletionStatus: pwdb.DeletionStatus_Active,
		}).
		Count(&seasonMemberShipCount).
		Error
	if err != nil || seasonMemberShipCount != 0 {
		return nil, errcode.ErrAlreadyHasTeamForSeason.Wrap(err)
	}

	// check if season rules are respected
	seasonRules := NewSeasonRules()
	err = seasonRules.ParseSeasonRulesString([]byte(team.Season.RulesBundle))

	if !seasonRules.IsStarted() {
		return nil, errcode.ErrSeasonIsNotStarted
	}

	if seasonRules.IsEnded() {
		return nil, errcode.ErrSeasonIsEnded
	}

	// retrieve total number of team members
	var teamMemberCount int32
	err = svc.db.
		Model(pwdb.TeamMember{}).
		Where(pwdb.TeamMember{TeamID: teamID}).
		Count(&teamMemberCount).
		Error
	if err != nil || seasonRules.IsLimitPlayersPerTeamReached(teamMemberCount) {
		return nil, errcode.ErrSeasonTeamLimitIsFull.Wrap(err)
	}

	// don't create new invite if user was already invited
	var teamInvite pwdb.TeamInvite
	err = svc.db.
		Where(pwdb.TeamInvite{
			UserID: invitedUserID,
			TeamID: teamID,
		}).
		First(&teamInvite).
		Error
	if err == nil {
		return nil, errcode.ErrAlreadyInvitedInTeam.Wrap(err)
	} else if err != gorm.ErrRecordNotFound {
		return nil, pwdb.GormToErrcode(err)
	}

	teamInvite = pwdb.TeamInvite{
		UserID: invitedUserID,
		TeamID: teamID,
	}
	// add invite to database
	err = svc.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&teamInvite).Error
		if err != nil {
			return pwdb.GormToErrcode(err)
		}
		activity := pwdb.Activity{
			Kind:           pwdb.Activity_TeamInviteSend,
			AuthorID:       userID,
			UserID:         invitedUserID,
			TeamID:         team.ID,
			TeamMemberID:   teamOwner.ID,
			OrganizationID: team.OrganizationID,
			SeasonID:       team.SeasonID,
			TeamInviteID:   teamInvite.ID,
		}
		return tx.Create(&activity).Error
	})
	if err != nil {
		return nil, err
	}

	// FIXME: Notify invited user

	ret := TeamSendInvite_Output{
		TeamInvite: &teamInvite,
	}
	return &ret, nil
}
