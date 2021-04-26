package pwdb

import (
	"errors"

	"gorm.io/gorm"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func IsRecordNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
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
