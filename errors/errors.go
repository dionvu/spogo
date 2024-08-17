package errors

import (
	"fmt"
	"os"

	"github.com/dionv/spogo/icons"
	"github.com/fatih/color"
	"github.com/joomcode/errorx"
)

var (
	AppNamespace = errorx.NewNamespace("app")

	HTTPError        = AppNamespace.NewType("http_error")
	HTTPRequestError = AppNamespace.NewType("http_request_error")

	FileError = AppNamespace.NewType("file_error")
	JSONError = AppNamespace.NewType("json_error")

	YAMLError = AppNamespace.NewType("yaml_error")

	ReauthenticationError = AppNamespace.NewType("reauthentication_error")

	ApiError = AppNamespace.NewType("api_error")

	PromptTuiError = AppNamespace.NewType("prompt_tui_error")

	PlayBack = errorx.NewNamespace("playback")

	DeviceError = PlayBack.NewType("no_device_error")

	Cli = errorx.NewNamespace("cli")

	NoFlagProvidedError = Cli.NewType("no_flag_provided_error")

	User = errorx.NewNamespace("user")

	InputError = User.NewType("input_error")
)

// If the error is not nil: prints the error, and exits the program.
func Catch(err error) {
	if err != nil {
		fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), err.(*errorx.Error).Message())
		os.Exit(0)
	}
}

// Prints the error even if it's nil.
func Print(err error) {
	fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), err.(*errorx.Error).Message())
}
