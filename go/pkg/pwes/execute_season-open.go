package pwes

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/pwdb"

	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
)

func (e *EventSeasonOpen) execute(ctx context.Context, apiClient *pwapi.HTTPClient, logger *zap.Logger) error {
	if apiClient == nil {
		logger.Debug("missing apiClient in execute event method")
		return errcode.ErrMissingInput
	}

	if e.Season == nil {
		logger.Debug("missing season input in execute EventSeasonOpen method")
		return errcode.ErrMissingInput
	}

	if e.Season.Subscription == pwdb.Season_Open {
		logger.Debug("season is already open")
		return nil
	}

	e.Season.Subscription = pwdb.Season_Open
	_, err := apiClient.AdminUpdateSeasonMetadata(ctx, &pwapi.AdminUpdateSeasonMetadata_Input{Season: e.Season})
	if err != nil {
		return err
	}
	return nil
}
