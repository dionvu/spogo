package errors

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
	"github.com/joomcode/errorx"
)

var (
	errorLogger *log.Logger
	apiLogger   *log.Logger
)

// Initiates both the error logger and api call logger.
func Init() {
	cacheDir, err := os.UserCacheDir()
	Catch(err)

	os.Mkdir(filepath.Join(cacheDir, "spogo"), 0777)
	os.Create(filepath.Join(cacheDir, "spogo", "spogo.log"))
	logFileErr, err := os.OpenFile(filepath.Join(cacheDir, "spogo", "errors.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open error log file: %v", err)
	}

	logFileApi, err := os.OpenFile(filepath.Join(cacheDir, "spogo", "api.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open error log file: %v", err)
	}

	errorLogger = log.New(logFileErr, "ERROR: ", log.Ldate|log.Ltime)
	apiLogger = log.New(logFileApi, "API: ", log.Ldate|log.Ltime)
}

func LogApiCall(endpoint string, statusCode int) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		apiLogger.Printf("%s:%v EP: %s, Status: %v\n", file, line, endpoint, statusCode)
	}
}

func Log(err error) {
	errorLogger.Println(err)
}

// If the error is not nil: prints the error, and exits the program.
func Catch(err error) {
	if err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err.(*errorx.Error).Message())
		os.Exit(0)
	}
}

// Prints the error even if it's nil.
func Print(err error) {
	fmt.Printf("%v %v\n", color.RedString("Error:"), err.(*errorx.Error).Message())
}

var (
	App           = errorx.NewNamespace("app")
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

	User             = errorx.NewNamespace("user")
	Reauthentication = User.NewType("reauthentication")
	NoFlagProvided   = User.NewType("no-flag")
	Input            = User.NewType("input")
	PlayBack         = errorx.NewNamespace("playback")
	NoDevice         = PlayBack.NewType("no-device")
	Jpeg             = PlayBack.NewType("jpeg")

	PlayerView             = errorx.NewNamespace("player-view")
	PlayerViewInvalidState = PlayerView.NewType("invalid-state")
	PlayerViewImageCache   = PlayerView.NewType("caching-image")
)
