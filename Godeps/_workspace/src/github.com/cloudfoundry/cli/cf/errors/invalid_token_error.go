package errors

import (
	. "github.com/emirozer/cf-fastpush-plugin/Godeps/_workspace/src/github.com/cloudfoundry/cli/cf/i18n"
)

type InvalidTokenError struct {
	description string
}

func NewInvalidTokenError(description string) error {
	return &InvalidTokenError{description: description}
}

func (err *InvalidTokenError) Error() string {
	return T("Invalid auth token: ") + err.description
}
