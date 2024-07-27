package utils

import (
	"fmt"
	"log"
	"os"
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

	fmt.Println("Opened browser?")
}

// Custom log.Fatal w/o datatime & w/ file & line #.
func LogError(msg string, err error) {
	const ERR = "\033[31mErr: \033[0m"

	fileAndLine := true

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	if fileAndLine {
		if err == nil {
			fmt.Printf("%s%s (file: %s, line: %d)\n%s\n", ERR, msg, file, line, err)
		} else {
			fmt.Printf("%s%s (file: %s, line: %d)\n%s\n", ERR, msg, file, line, err)
		}
	} else {
		if err == nil {
			fmt.Printf("%s%s\n%s\n", ERR, msg, err)
		} else {
			fmt.Printf("%s%s\n%s\n", ERR, msg, err)
		}
	}

	os.Exit(0)
}
