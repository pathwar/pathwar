package errcode

import (
	"fmt"
	"io"

	"golang.org/x/xerrors"
)

type WithCode interface {
	error
	Code() int32
	Format(fmt.State, rune)
}

// Code returns the code of the actual error without trying to unwrap it, or -1.
func Code(err error) int32 {
	typed, ok := err.(WithCode)
	if ok {
		return typed.Code()
	}
	return -1
}

// LastCode walks the passed error and returns the code of the latest ErrCode, or -1.
func LastCode(err error) int32 {
	if err == nil {
		return -1
	}

	if cause := genericCause(err); cause != nil {
		if ret := LastCode(cause); ret != -1 {
			return ret
		}
	}

	return Code(err)
}

// FirstCode walks the passed error and returns the code of the first ErrCode met, or -1.
func FirstCode(err error) int32 {
	if err == nil {
		return -1
	}

	if code := Code(err); code != -1 {
		return code
	}

	if cause := genericCause(err); cause != nil {
		return FirstCode(cause)
	}

	return -1
}

func genericCause(err error) error {
	type causer interface {
		Cause() error
	}
	type wrapper interface {
		Unwrap() error
	}

	if causer, ok := err.(causer); ok {
		return causer.Cause()
	}

	if wrapper, ok := err.(wrapper); ok {
		return wrapper.Unwrap()
	}

	return nil
}

//
// Error
//

func (e ErrCode) Error() string {
	name, ok := ErrCode_name[int32(e)]
	if ok {
		return fmt.Sprintf("%s(#%d)", name, int32(e))
	}
	return fmt.Sprintf("UNKNOWN_ERRCODE(#%d)", int32(e))
}

func (e ErrCode) Code() int32 {
	return int32(e)
}

func (e ErrCode) Wrap(inner error) WithCode {
	return wrappedError{
		code:  int32(e),
		inner: inner,
		frame: xerrors.Caller(1),
	}
}

func (e ErrCode) Format(f fmt.State, c rune) {
	xerrors.FormatError(e, f, c)
}

func (e ErrCode) FormatError(p xerrors.Printer) error {
	p.Print(e.Error())
	return nil
}

//
// ConfigurableError
//

type wrappedError struct {
	code  int32
	inner error
	frame xerrors.Frame
}

func (e wrappedError) Error() string {
	return fmt.Sprintf("%s: %v", ErrCode(e.code), e.inner)
}

func (e wrappedError) Code() int32 {
	return e.code
}

func (e wrappedError) Format(f fmt.State, c rune) {
	xerrors.FormatError(e, f, c)
	if f.Flag('+') {
		_, _ = io.WriteString(f, "\n")
		if sub := genericCause(e); sub != nil {
			if typed, ok := sub.(wrappedError); ok {
				sub = lightWrappedError{wrappedError: typed}
			}
			formatter, ok := sub.(fmt.Formatter)
			if ok {
				formatter.Format(f, c)
			}
		}
	}
}

func (e wrappedError) FormatError(p xerrors.Printer) error {
	p.Print(e.Error())
	if p.Detail() {
		e.frame.Format(p)
	}
	return nil
}

// Cause returns the inner error (github.com/pkg/errors)
func (e wrappedError) Cause() error {
	return e.inner
}

// Unwrap returns the inner error (go1.13)
func (e wrappedError) Unwrap() error {
	return e.inner
}

// light wrapped errors

type lightWrappedError struct {
	wrappedError
	deepness int
}

func (e lightWrappedError) Error() string {
	return ""
}

func (e lightWrappedError) Format(f fmt.State, c rune) {
	xerrors.FormatError(e, f, c)
	if f.Flag('+') {
		_, _ = io.WriteString(f, "\n")
		if sub := genericCause(e); sub != nil {
			if typed, ok := sub.(wrappedError); ok {
				sub = lightWrappedError{wrappedError: typed, deepness: e.deepness + 1}
			}
			formatter, ok := sub.(fmt.Formatter)
			if ok {
				formatter.Format(f, c)
			}
		}
	}
}

func (e lightWrappedError) FormatError(p xerrors.Printer) error {
	p.Printf("#%d", e.deepness+1)
	e.frame.Format(p)
	return nil
}
