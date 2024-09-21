package utils

import (
	"fmt"
	"io"
)

// Use to debug wacky api json confuzzling.
func PrintResponseBody(r io.ReadCloser) {
	b, _ := io.ReadAll(r)
	defer r.Close()
	fmt.Println(string(b))
}
