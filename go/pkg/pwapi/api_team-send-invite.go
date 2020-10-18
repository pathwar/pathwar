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
		return nil, pwdb.GormToErrcode(err)
	}

	invitedUserID, err := pwdb.GetIDBySlugAndKind(svc.db, in.UserID, "user")
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	// check team status
	var team pwdb.Team
	err = svc.db.
		Where(pwdb.Team{
			ID:             teamID,
			DeletionStatus: pwdb.DeletionStatus_Active,
		}).
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

	// check that the invited user is not already in this season
	var existCheck pwdb.TeamMember
	err = svc.db.
		Where(pwdb.TeamMember{
			UserID: invitedUserID,
		}).
		Where(&pwdb.Season{
			ID: team.SeasonID,
		}).
		Joins("JOIN team ON team_member.team_id = team.id").
		Joins("JOIN season ON team.season_id = season.id").
		First(&existCheck).
		Error
	if err == nil {
		return nil, errcode.ErrAlreadyTeamMember.Wrap(err)
	} else if err != gorm.ErrRecordNotFound {
		return nil, pwdb.GormToErrcode(err)
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
