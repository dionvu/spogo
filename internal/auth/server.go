package auth

import (
	"net/http"
	"os"

	"github.com/dionv/spogo/utils"
)

func startServer() {
	port := os.Getenv("PORT")

	go func() {
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			utils.LogError("Failed to starting server", err)
		}
	}()
}
