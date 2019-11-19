package pwengine

import "errors"

var (
	ErrNotImplemented            = errors.New("not implemented")
	ErrMissingArgument           = errors.New("missing argument(s)")
	ErrInvalidArgument           = errors.New("invalid argument(s)")
	ErrDuplicate                 = errors.New("duplicate")
	ErrMissingRequiredValidation = errors.New("missing required validation")
	ErrInternalServerError       = errors.New("internal server error")
)
