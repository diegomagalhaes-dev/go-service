package auth

import (
	"errors"
	"fmt"
)

type authError struct {
	msg string
}

func NewAuthError(format string, args ...any) error {
	return &authError{
		msg: fmt.Sprintf(format, args...),
	}
}

func (ae *authError) Error() string {
	return ae.msg
}

func IsAuthError(err error) bool {
	var ae *authError
	return errors.As(err, &ae)
}
