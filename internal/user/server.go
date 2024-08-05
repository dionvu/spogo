package user

import (
	"log"
	"net/http"
)

func startServer() {
	go func() {
		err := http.ListenAndServe(":"+PORT, nil)
		if err != nil {
			log.Fatalf("Failed to start server on port: %v", PORT)
		}
	}()
}
