package pwapi

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
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

	err = svc.db.Transaction(func(tx *gorm.DB) error {
		// update user
		now := time.Now()
		updates := pwdb.User{
			OAuthSubject:   fmt.Sprintf("deleted_%s_%d", user.OAuthSubject, now.Unix()),
			DeletionReason: in.Reason,
			DeletionStatus: pwdb.DeletionStatus_Requested,
			DeletedAt:      &now,
		}
		err = tx.Model(&user).Updates(updates).Error
		if err != nil {
			return errcode.ErrUpdateUser.Wrap(err)
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
				err = tx.Model(&teamMembership.Team).Updates(updates).Error
				if err != nil {
					return errcode.ErrUpdateTeam.Wrap(err)
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
				err = tx.Model(&organizationMembership.Organization).Updates(updates).Error
				if err != nil {
					return errcode.ErrUpdateOrganization
				}
			}
		}

		// FIXME: invalide current JWT token
		// FIXME: add another task that pseudonymize the data

		// create activity
		activity := pwdb.Activity{
			Kind:     pwdb.Activity_UserDeleteAccount,
			AuthorID: user.ID,
			UserID:   user.ID,
		}
		return tx.Create(&activity).Error
	})

	ret := &UserDeleteAccount_Output{}
	return ret, err
}
