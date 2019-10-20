package pwdb

import (
	"errors"

	"github.com/jinzhu/gorm"
)

func IsRecordNotFoundError(err error) bool {
	if unwrapped := errors.Unwrap(err); unwrapped != nil {
		err = unwrapped
	}
	return gorm.IsRecordNotFoundError(err)
}
