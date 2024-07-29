package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	errx "github.com/dionv/spogo/errors"
)

func ParseJsonResponse(res *http.Response) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errx.FileError.Wrap(err, "Failed to read response body")
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errx.JSONError.Wrap(err, "Failed to unmarshal response body")
	}

	return data, nil
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
