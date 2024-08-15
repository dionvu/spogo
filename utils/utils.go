package utils

import (
	"fmt"
	"io"
	"os/exec"
	"runtime"

	"github.com/dionv/spogo/icons"
	"github.com/fatih/color"
)

// Opens a url depending on user's system.
func OpenURL(url string) error {
	var cmd *exec.Cmd

	os := runtime.GOOS

	switch {
	case os == "windows":
		cmd = exec.Command("start", url)
	case os == "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	fmt.Println(color.HiGreenString(icons.Attention+"Opening -> ", url))

	return nil
}

// Use to debug wacky api json confuzzling.
func PrintResponseBody(r io.ReadCloser) {
	b, _ := io.ReadAll(r)
	defer r.Close()
	fmt.Println(string(b))
}
