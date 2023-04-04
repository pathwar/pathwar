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
		return nil, err
	}

	var teamInvite pwdb.TeamInvite
	err = svc.db.
		Where(&pwdb.TeamInvite{
			ID:     teamInviteID,
			UserID: userID,
		}).
		Preload("Team").
		Preload("Team.Season").
		First(&teamInvite).
		Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	if !in.Accept {
		err = svc.db.Transaction(func(tx *gorm.DB) error {
			err = tx.Delete(&teamInvite).Error
			if err != nil {
				return err
			}
			activity := pwdb.Activity{
				Kind:         pwdb.Activity_TeamInviteDecline,
				AuthorID:     userID,
				UserID:       teamInvite.UserID,
				TeamInviteID: teamInvite.ID,
			}

			err = tx.Create(&activity).Error
			if err != nil {
				return pwdb.GormToErrcode(err)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return &TeamAcceptInvite_Output{}, nil
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
			DeletionStatus: pwdb.DeletionStatus_Active,
		}).
		Count(&seasonMemberShipCount).
		Error
	if err != nil || seasonMemberShipCount != 0 {
		return nil, errcode.ErrAlreadyHasTeamForSeason.Wrap(err)
	}

	// check if season rules are respected
	seasonRules := NewSeasonRules()
	err = seasonRules.ParseSeasonRulesString([]byte(teamInvite.Team.Season.RulesBundle))
	if err != nil {
		return nil, errcode.ErrParseSeasonRule.Wrap(err)
	}

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
		Where(pwdb.TeamMember{TeamID: teamInvite.Team.ID}).
		Count(&teamMemberCount).
		Error
	if err != nil || seasonRules.IsLimitPlayersPerTeamReached(teamMemberCount) {
		return nil, errcode.ErrSeasonTeamLimitIsFull.Wrap(err)
	}

	var user pwdb.User
	err = svc.db.Model(pwdb.User{}).Select("email").First(&user, userID).Error
	if err != nil {
		return nil, errcode.ErrGetUser.Wrap(err)
	}

	if !seasonRules.IsEmailDomainAllowed(user.Email) {
		return nil, errcode.ErrSeasonEmailDomainNotAllowed
	}

	teamMember := &pwdb.TeamMember{
		UserID: userID,
		TeamID: teamInvite.TeamID,
	}
	var organizationMembership int
	err = svc.db.
		Model(&pwdb.OrganizationMember{}).
		Where(pwdb.OrganizationMember{
			UserID:         userID,
			OrganizationID: teamInvite.Team.OrganizationID,
		}).
		Count(&organizationMembership).
		Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	var orgaMember *pwdb.OrganizationMember
	if organizationMembership == 0 {
		orgaMember = &pwdb.OrganizationMember{
			UserID:         userID,
			OrganizationID: teamInvite.Team.OrganizationID,
		}
	}
	err = svc.db.Transaction(func(tx *gorm.DB) error {
		// create team member
		err = tx.Create(&teamMember).Error
		if err != nil {
			return pwdb.GormToErrcode(err)
		}

		// create orga member if needed
		if orgaMember != nil {
			err = tx.Create(&orgaMember).Error
			if err != nil {
				return pwdb.GormToErrcode(err)
			}
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
