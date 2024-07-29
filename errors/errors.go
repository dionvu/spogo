package errors

import "github.com/joomcode/errorx"

var (
	AppNamespace = errorx.NewNamespace("auth")

	HTTPError        = AppNamespace.NewType("http_error")
	HTTPRequestError = AppNamespace.NewType("http_request_error")

	FileError = AppNamespace.NewType("file_error")
	JSONError = AppNamespace.NewType("json_error")

	YAMLError = AppNamespace.NewType("yaml_error")

	TokenRefreshError     = AppNamespace.NewType("token_refresh_error")
	InvalidTokenError     = AppNamespace.NewType("invalid_token_error")
	ReauthenticationError = AppNamespace.NewType("reauthentication_error")
	InvalidStateError     = AppNamespace.NewType("invalid_state_error")

	ApiError = AppNamespace.NewType("api_error")

	EncryptionError = AppNamespace.NewType("encryption_error")
	DecryptionError = AppNamespace.NewType("decryption_error")
)
