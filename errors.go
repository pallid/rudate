package rudate

import "errors"

var (
	ErrNotFound   = errors.New("rudate: date/time not found in text")
	ErrNoDuration = errors.New("rudate: duration not recognized")
)

type ParseError struct {
	Inner error
}

func (e *ParseError) Error() string {
	return "rudate: " + e.Inner.Error()
}

func (e *ParseError) Unwrap() error {
	return e.Inner
}

func (e *ParseError) Is(target error) bool {
	_, ok := target.(*ParseError)
	return ok
}
