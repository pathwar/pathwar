package pwapi

import (
	"context"

	"github.com/jinzhu/gorm"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) TeamAcceptInvite(ctx context.Context, in *TeamAcceptInvite_Input) (*TeamAcceptInvite_Output, error) {
	if in == nil || in.TeamInviteID == "" {
		return nil, errcode.ErrMissingInput
	}

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated.Wrap(err)
	}

	teamInviteID, err := pwdb.GetIDBySlugAndKind(svc.db, in.TeamInviteID, "team-invite")
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	var teamInvite pwdb.TeamInvite
	err = svc.db.
		Where(&pwdb.TeamInvite{
			ID:     teamInviteID,
			UserID: userID,
		}).
		Preload("Team").
		First(&teamInvite).
		Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	// check if user already has a team in this season
	var seasonMemberShipCount int
	err = svc.db.
		Model(pwdb.TeamMember{}).
		Preload("Team").
		Joins("JOIN team on team.id = team_member.team_id").
		Where(pwdb.TeamMember{UserID: userID}).
		Where(&pwdb.Team{
			SeasonID:       teamInvite.Team.SeasonID,
			DeletionStatus: pwdb.DeletionStatus_Active}).
		Count(&seasonMemberShipCount).
		Error
	if err != nil || seasonMemberShipCount != 0 {
		return nil, errcode.ErrAlreadyHasTeamForSeason.Wrap(err)
	}

	teamMember := &pwdb.TeamMember{
		UserID: userID,
		TeamID: teamInvite.TeamID,
	}
	orgaMember := &pwdb.OrganizationMember{
		UserID:         userID,
		OrganizationID: teamInvite.Team.OrganizationID,
	}
	err = svc.db.Transaction(func(tx *gorm.DB) error {
		// create team member
		err = tx.Create(&teamMember).Error
		if err != nil {
			return pwdb.GormToErrcode(err)
		}
		// create orga member
		err = tx.Create(&orgaMember).Error
		if err != nil {
			return pwdb.GormToErrcode(err)
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	err = svc.db.Transaction(func(tx *gorm.DB) error {
		// remove invite
		err = tx.Delete(&teamInvite).Error
		if err != nil {
			return pwdb.GormToErrcode(err)
		}

		activity := pwdb.Activity{
			Kind:           pwdb.Activity_TeamInviteAccept,
			AuthorID:       userID,
			UserID:         teamInvite.UserID,
			TeamID:         teamInvite.TeamID,
			TeamMemberID:   teamMember.ID,
			OrganizationID: teamInvite.Team.OrganizationID,
			SeasonID:       teamInvite.Team.SeasonID,
		}
		return tx.Create(&activity).Error
	})
	if err != nil {
		return nil, err
	}

	ret := TeamAcceptInvite_Output{
		TeamMember: teamMember,
	}
	return &ret, nil
}
