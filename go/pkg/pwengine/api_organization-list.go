package pwengine

import (
	"context"
	"fmt"
)

func (e *engine) OrganizationList(context.Context, *OrganizationList_Input) (*OrganizationList_Output, error) {
	var organizations OrganizationList_Output
	err := e.db.
		Set("gorm:auto_preload", true). // FIXME: explicit preloading
		Find(&organizations.Items).Error
	if err != nil {
		return nil, fmt.Errorf("query organizations: %w", err)
	}

	return &organizations, nil
}
