package errors

import "github.com/joomcode/errorx"

var (
	AppNamespace = errorx.NewNamespace("app")

	HTTPError        = AppNamespace.NewType("http_error")
	HTTPRequestError = AppNamespace.NewType("http_request_error")

	FileError = AppNamespace.NewType("file_error")
	JSONError = AppNamespace.NewType("json_error")

	YAMLError = AppNamespace.NewType("yaml_error")

	InvalidTokenError     = AppNamespace.NewType("invalid_token_error")
	ReauthenticationError = AppNamespace.NewType("reauthentication_error")

	ApiError = AppNamespace.NewType("api_error")

	EncryptionError = AppNamespace.NewType("encryption_error")
	DecryptionError = AppNamespace.NewType("decryption_error")

	PLAYBACKERROR = AppNamespace.NewType("playback_error")
)
