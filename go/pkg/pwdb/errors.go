package pwdb

import (
	"errors"

	"github.com/jinzhu/gorm"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func GormToErrcode(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errcode.ErrDBNotFound.Wrap(err)
	}

	if err != nil {
		return errcode.ErrDBInternal.Wrap(err)
	}

	return nil
}
