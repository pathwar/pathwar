package pwapi

import (
	"context"
	"fmt"
	"time"

	"pathwar.land/go/v2/pkg/errcode"
	"pathwar.land/go/v2/pkg/pwdb"
)

func (svc *service) UserDeleteAccount(ctx context.Context, in *UserDeleteAccount_Input) (*UserDeleteAccount_Output, error) {
	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated.Wrap(err)
	}

	// get user
	var user pwdb.User
	err = svc.db.
		Preload("TeamMemberships").
		Preload("TeamMemberships.Team.Members.User").
		Preload("OrganizationMemberships").
		Preload("OrganizationMemberships.Organization.Members.User").
		First(&user, userID).
		Error
	if err != nil {
		return nil, errcode.ErrGetUser.Wrap(err)
	}

	//fmt.Println(godev.PrettyJSON(user))

	// update user
	now := time.Now()
	updates := pwdb.User{
		OAuthSubject:   fmt.Sprintf("deleted_%s_%d", user.OAuthSubject, now.Unix()),
		DeletionReason: in.Reason,
		DeletionStatus: pwdb.DeletionStatus_Requested,
		DeletedAt:      &now,
	}
	// FIXME: use transaction
	err = svc.db.Model(&user).Updates(updates).Error
	if err != nil {
		return nil, errcode.ErrUpdateUser.Wrap(err)
	}

	// update teams
	for _, teamMembership := range user.TeamMemberships {
		haveAnotherActiveMember := false
		for _, member := range teamMembership.Team.Members {
			if member.User.ID == user.ID {
				continue
			}
			if member.User.DeletionStatus == pwdb.DeletionStatus_Active {
				haveAnotherActiveMember = true
				break
			}
		}
		if !haveAnotherActiveMember {
			updates := pwdb.Team{
				DeletionStatus: pwdb.DeletionStatus_Requested,
				DeletedAt:      &now,
			}
			err = svc.db.Model(&teamMembership.Team).Updates(updates).Error
			if err != nil {
				return nil, errcode.ErrUpdateTeam.Wrap(err)
			}
		}
	}

	// update organizations
	for _, organizationMembership := range user.OrganizationMemberships {
		haveAnotherActiveMember := false
		for _, member := range organizationMembership.Organization.Members {
			if member.User.ID == user.ID {
				continue
			}
			if member.User.DeletionStatus == pwdb.DeletionStatus_Active {
				haveAnotherActiveMember = true
				break
			}
		}
		if !haveAnotherActiveMember {
			updates := pwdb.Organization{
				Name:           fmt.Sprintf("deleted_%s_%d", organizationMembership.Organization.Name, now.Unix()),
				DeletionStatus: pwdb.DeletionStatus_Requested,
				DeletedAt:      &now,
			}
			err = svc.db.Model(&organizationMembership.Organization).Updates(updates).Error
			if err != nil {
				return nil, errcode.ErrUpdateOrganization
			}
		}
	}

	// FIXME: invalide current JWT token
	// FIXME: add another task that pseudonymize the data

	ret := &UserDeleteAccount_Output{}
	return ret, nil
}
