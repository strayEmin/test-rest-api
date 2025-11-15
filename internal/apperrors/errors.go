package apperrors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInvironment          = errors.New("invalid environment")
	ErrHashingPasswordFailed       = errors.New("hashing password failed")
	ErrUnsuccesffulTokenGeneration = errors.New("unsuccesfful token generation")
	ErrEntityNotFound              = errors.New("entity not found")
	ErrUserAlreadyExists           = errors.New("user already exists")
	ErrUserNotFound                = errors.New("user not found")
	ErrWrongPassword               = errors.New("wrong password")
)

func Wrap(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", op, err)
}
