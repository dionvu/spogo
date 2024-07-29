package auth

import (
	"log"
	"net/http"
	"os"
)

func startServer() {
	port := os.Getenv("PORT")

	go func() {
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatalf("Failed to start server on port: %v", port)
		}
	}()
}
