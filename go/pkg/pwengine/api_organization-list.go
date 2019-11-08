package pwengine

import (
	"context"
	"fmt"
)

func (e *engine) OrganizationList(context.Context, *OrganizationListInput) (*OrganizationListOutput, error) {
	var organizations OrganizationListOutput
	err := e.db.
		Set("gorm:auto_preload", true). // FIXME: explicit preloading
		Find(&organizations.Items).Error
	if err != nil {
		return nil, fmt.Errorf("query organizations: %w", err)
	}

	return &organizations, nil
}
