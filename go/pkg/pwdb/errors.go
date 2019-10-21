package pwdb

import (
	"errors"

	"github.com/jinzhu/gorm"
)

func IsRecordNotFoundError(err error) bool {
	return gorm.IsRecordNotFoundError(err) ||
		gorm.IsRecordNotFoundError(errors.Unwrap(err))
}
