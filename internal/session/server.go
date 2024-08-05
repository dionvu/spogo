package session

import (
	"log"
	"net/http"
)

// Starts the fucking server :DD.
func startServer() {
	go func() {
		err := http.ListenAndServe(":"+PORT, nil)
		if err != nil {
			log.Fatalf("Failed to start server on port: %v", PORT)
		}
	}()
}
