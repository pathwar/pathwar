package pwes

import (
	"context"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
)

func (e EventAgentRegister) execute(ctx context.Context, apiClient *pwapi.HTTPClient) error {
	if apiClient == nil {
		return errcode.ErrMissingInput
	}
	return nil
}
