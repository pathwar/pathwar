package pwes

import (
	"context"
	"fmt"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
)

func Compute(ctx context.Context, apiClient *pwapi.HTTPClient) error {
	activities, err := apiClient.AdminListActivities(ctx, &pwapi.AdminListActivities_Input{FilteringPreset: "validations"})
	if err != nil {
		return err
	}
	fmt.Println(activities)
	return nil
}
