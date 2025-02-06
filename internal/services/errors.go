package services

import (
	"errors"
	"fmt"
)

var (
	ErrPermissionDenied = errors.New("permission denied")
)

var ErrInvalidArgument = errors.New("invalid argument")

func NewValidationError(param string) error {
	return fmt.Errorf("%w: %s cannot be empty", ErrInvalidArgument, param)
}
