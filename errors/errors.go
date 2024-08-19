package errors

import (
	"fmt"
	"os"

	"github.com/dionv/spogo/icons"
	"github.com/fatih/color"
	"github.com/joomcode/errorx"
)

var (
	App = errorx.NewNamespace("app")

	HTTP          = App.NewType("http")
	HTTPRequest   = App.NewType("http-request")
	File          = App.NewType("file")
	FileOpen      = App.NewType("file-open")
	FileCreate    = App.NewType("file-create")
	FileRead      = App.NewType("file-read")
	FileWrite     = App.NewType("file-write")
	JSON          = App.NewType("json")
	JSONUnmarshal = App.NewType("json-unmarshal")
	JSONMarshal   = App.NewType("json-marshal")
	JSONEncode    = App.NewType("json-encode")
	JSONDecode    = App.NewType("json-decode")
	YAML          = App.NewType("yaml")
)

var (
	Dependency = errorx.NewNamespace("dependency")

	PromptTui = Dependency.NewType("prompt-tui")
)

var (
	User = errorx.NewNamespace("user")

	Reauthentication = User.NewType("reauthentication")
	NoFlagProvided   = User.NewType("no-flag")
	Input            = User.NewType("input")
)

var (
	PlayBack = errorx.NewNamespace("playback")

	NoDevice = PlayBack.NewType("no-device")
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
