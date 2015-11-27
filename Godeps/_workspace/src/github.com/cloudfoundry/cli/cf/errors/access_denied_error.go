package errors

import . "github.com/emirozer/cf-fastpush-plugin/Godeps/_workspace/src/github.com/cloudfoundry/cli/cf/i18n"

type AccessDeniedError struct {
}

func NewAccessDeniedError() *AccessDeniedError {
	return &AccessDeniedError{}
}

func (err *AccessDeniedError) Error() string {
	return T("Server error, status code: 403: Access is denied.  You do not have privileges to execute this command.")
}
