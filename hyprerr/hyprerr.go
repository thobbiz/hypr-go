package hyprerr

import (
	"errors"
	"fmt"
)

type ErrType string

const (
	ErrIPC          ErrType = "ipc_error"
	ErrNotFound     ErrType = "not_found"
	ErrIO           ErrType = "io_error"
	ErrParse        ErrType = "parse_error"
	ErrNoMatch      ErrType = "no_match"
	ErrInvalidInput ErrType = "invalid_input"
)

type HyprError struct {
	Op   string
	Kind ErrType
	Err  error
}

func (e *HyprError) Error() string {
	return fmt.Sprintf("%s [%s]: %v", e.Op, e.Kind, e.Err)
}

func (e *HyprError) Unwrap() error { return e.Err }

func New(op string, kind ErrType, err error) *HyprError {
	return &HyprError{Op: op, Kind: kind, Err: err}
}

func IsKind(err error, kind ErrType) bool {
	var e *HyprError
	return errors.As(err, &e) && e.Kind == kind
}
