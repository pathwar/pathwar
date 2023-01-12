package pwapi

import (
	"context"
	"fmt"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminListChallengeSubscriptions(ctx context.Context, in *AdminListChallengeSubscriptions_Input) (*AdminListChallengeSubscriptions_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	fmt.Println("WESH", in, "WESH")
	var challengeSubscriptions []*pwdb.ChallengeSubscription
	req := svc.db.
		Preload("Team").
		Preload("Team.Organization").
		Preload("Closer").
		Preload("Validations").
		Preload("SeasonChallenge").
		Preload("SeasonChallenge.Season").
		Preload("SeasonChallenge.Flavor").
		Preload("SeasonChallenge.Flavor.Challenge").
		Preload("Buyer")
	if in.SeasonChallengeID != "" {
		req = req.Where("season_challenge_id = ?", in.SeasonChallengeID)
	}
	switch in.FilteringPreset {
	case "default", "":
		// noop
	case "closed":
		req = req.Where(&pwdb.ChallengeSubscription{Status: pwdb.ChallengeSubscription_Closed})
	case "open":
		req = req.Where(&pwdb.ChallengeSubscription{Status: pwdb.ChallengeSubscription_Active})
	default:
		return nil, errcode.TODO
	}
	err := req.Find(&challengeSubscriptions).Error
	if err != nil {
		return nil, errcode.ErrListChallengeSubscriptions.Wrap(err)
	}

	out := AdminListChallengeSubscriptions_Output{Subscriptions: challengeSubscriptions}
	return &out, nil
}
