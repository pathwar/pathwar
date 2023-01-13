package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) OrganizationCreate(ctx context.Context, in *OrganizationCreate_Input) (*OrganizationCreate_Output, error) {
	if in == nil || in.Name == "" {
		return nil, errcode.ErrMissingInput
	}
	return &OrganizationCreate_Output{
		Organization: &pwdb.Organization{},
	}, nil
}
