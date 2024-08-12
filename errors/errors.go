package errors

import "github.com/joomcode/errorx"

var (
	AppNamespace = errorx.NewNamespace("app")

	HTTPError        = AppNamespace.NewType("http_error")
	HTTPRequestError = AppNamespace.NewType("http_request_error")

	FileError = AppNamespace.NewType("file_error")
	JSONError = AppNamespace.NewType("json_error")

	YAMLError = AppNamespace.NewType("yaml_error")

	// InvalidTokenError     = AppNamespace.NewType("invalid_token_error")
	ReauthenticationError = AppNamespace.NewType("reauthentication_error")

	ApiError = AppNamespace.NewType("api_error")

	PromptTuiError = AppNamespace.NewType("prompt_tui_error")

	PlayBack = errorx.NewNamespace("playback")

	NoDeviceError = PlayBack.NewType("no_device_error")

	/////////
	Cli = errorx.NewNamespace("cli")

	NoFlagProvidedError = Cli.NewType("no_flag_provided_error")
)
