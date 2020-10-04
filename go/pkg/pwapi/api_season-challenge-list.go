package pwapi

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) SeasonChallengeList(ctx context.Context, in *SeasonChallengeList_Input) (*SeasonChallengeList_Output, error) {
	if in == nil || in.SeasonID == 0 {
		return nil, errcode.ErrMissingInput
	}

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated.Wrap(err)
	}

	exists, err := seasonIDExists(svc.db, in.SeasonID)
	if err != nil || !exists {
		return nil, errcode.ErrInvalidSeasonID.Wrap(err)
	}

	team, err := userTeamForSeason(svc.db, userID, in.SeasonID)
	if err != nil {
		return nil, errcode.ErrUserHasNoTeamForSeason.Wrap(err)
	}

	var seasonChallenges []*pwdb.SeasonChallenge
	err = svc.db.
		// Preload("Season").
		Preload("Flavor").
		Joins("LEFT JOIN challenge_flavor ON season_challenge.flavor_id = challenge_flavor.id").
		Preload("Flavor.Challenge").
		Preload("Flavor.Instances").
		Preload("Flavor.Instances.Agent"). // FIXME: where status==active
		Preload("Subscriptions", "team_id = ?", team.ID).
		Preload("Subscriptions.Team").
		Preload("Subscriptions.Team.Season").
		Where(pwdb.SeasonChallenge{SeasonID: in.SeasonID}).
		Order("challenge_flavor.purchase_price asc, challenge_flavor.validation_reward asc").
		Find(&seasonChallenges).
		Error
	if err != nil {
		return nil, errcode.ErrGetSeasonChallenges.Wrap(err)
	}

	// prepare & cleanup
	for _, sc := range seasonChallenges {
		// FIXME: hide challenges without flavors?
		// fmt.Println(sc.ID, godev.PrettyJSON(sc.Flavor.Instances))
		if sc.Flavor.TagList != "" {
			sc.Flavor.Tags = strings.Split(sc.Flavor.TagList, ",")
			sc.Flavor.TagList = ""
		}
		if sc.Flavor.RedumpPolicyConfig != "" {
			err := json.Unmarshal([]byte(sc.Flavor.RedumpPolicyConfig), &sc.Flavor.RedumpPolicy)
			if err == nil {
				sc.Flavor.RedumpPolicyConfig = ""
			}
		}
		for _, instance := range sc.Flavor.Instances {
			// FIXME: hide instances without nginx-url?
			instance.InstanceConfig = nil
			if instance.Agent != nil {
				if len(sc.Subscriptions) > 0 {
					hash, err := pwdb.ChallengeInstancePrefixHash(fmt.Sprintf("%d", instance.ID), userID, instance.Agent.AuthSalt)
					if err != nil {
						return nil, errcode.ErrGeneratePrefixHash.Wrap(err)
					}
					instance.NginxURL = fmt.Sprintf("http://%s.%s", hash, instance.Agent.DomainSuffix)
				}
				instance.AgentID = 0
				instance.Agent = nil
			}
		}
	}

	ret := SeasonChallengeList_Output{
		Items: seasonChallenges,
	}
	return &ret, nil
}
