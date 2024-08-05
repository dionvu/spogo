package utils

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

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
