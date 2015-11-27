package errors

import (
	. "github.com/emirozer/cf-fastpush-plugin/Godeps/_workspace/src/github.com/cloudfoundry/cli/cf/i18n"
)

type NotAuthorizedError struct {
}

func NewNotAuthorizedError() error {
	return &NotAuthorizedError{}
}

func (err *NotAuthorizedError) Error() string {
	return T("Server error, status code: 403, error code: 10003, message: You are not authorized to perform the requested action")
}
