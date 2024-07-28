package errors

import "github.com/joomcode/errorx"

var (
	AppNamespace = errorx.NewNamespace("auth")

	FileError             = AppNamespace.NewType("file_error")
	JSONError             = AppNamespace.NewType("json_error")
	HTTPRequestError      = AppNamespace.NewType("http_request_error")
	TokenRefreshError     = AppNamespace.NewType("token_refresh_error")
	InvalidTokenError     = AppNamespace.NewType("invalid_token_error")
	ReauthenticationError = AppNamespace.NewType("reauthentication_error")
	InvalidStateError     = AppNamespace.NewType("invalid_state")
)
