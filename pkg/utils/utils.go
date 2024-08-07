package utils

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
)

func CatchErr(err error) {
	if err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error"), err)
	}
}

// Opens a url depending on user's system?
// Maybe lol, haven't tested with windows or mac.
func OpenURL(url string) {
	var cmd *exec.Cmd

	os := runtime.GOOS

	switch {

	// I haven't tested this lmao
	case os == "windows":
		cmd = exec.Command("start", url)

	// This either
	case os == "darwin":
		cmd = exec.Command("open", url)

	// Linux
	default:
		cmd = exec.Command("xdg-open", url)
	}

	err := cmd.Start()
	if err != nil {
		log.Fatal("Failed to open browser: ", err)
	}

	fmt.Println("Opening browser to -> ", url)
}

func PrintResponseBody(r io.ReadCloser) {
	b, _ := io.ReadAll(r)
	defer r.Close()
	fmt.Println(string(b))
}
