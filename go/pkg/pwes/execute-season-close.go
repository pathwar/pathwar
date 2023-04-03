package pwes

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/pwdb"

	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
)

func (e *EventSeasonClose) execute(ctx context.Context, apiClient *pwapi.HTTPClient, logger *zap.Logger) error {
	if apiClient == nil {
		logger.Debug("missing apiClient in execute event method")
		return errcode.ErrMissingInput
	}

	if e.Season == nil {
		logger.Debug("missing season input in execute EventSeasonClose method")
		return errcode.ErrMissingInput
	}

	if e.Season.Subscription == pwdb.Season_Close {
		logger.Debug("season is already close")
		return nil
	}

	e.Season.Subscription = pwdb.Season_Close

	return nil
}
