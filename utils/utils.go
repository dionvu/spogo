package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Use to debug wacky api json confuzzling.
func PrintResponseBody(r io.ReadCloser) {
	b, _ := io.ReadAll(r)
	defer r.Close()
	fmt.Println(string(b))
}

// Use to debug wacky api json confuzzling.
func ResponseBody(r io.ReadCloser) string {
	b, _ := io.ReadAll(r)
	defer r.Close()
	return string(b)
}

// Caches the given image.
// Deprecated: Use Image.Cache instead.
func CacheImage(url string, path string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	return nil
}

// Converts the number of milliseconds into two string values
// of minutes and addittional seconds.
func MsToMinutesAndSeconds(ms int) (minutes string, seconds string) {
	m := ms / 60000
	s := (ms % 60000) / 1000

	minutes = fmt.Sprint(m)
	seconds = fmt.Sprint(s)

	return minutes, seconds
}
