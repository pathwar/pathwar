package pwdb

import (
	"errors"

	"github.com/jinzhu/gorm"
	"pathwar.land/v2/go/pkg/errcode"
)

func IsRecordNotFoundError(err error) bool {
	return gorm.IsRecordNotFoundError(err) ||
		gorm.IsRecordNotFoundError(errors.Unwrap(err))
}

func GormToErrcode(err error) error {
	if IsRecordNotFoundError(err) {
		return errcode.ErrDBNotFound.Wrap(err)
	}

	if err != nil {
		return errcode.ErrDBInternal.Wrap(err)
	}

	return nil
}
