package main

import (
	"log"

	"github.com/dionv/spogo/internal/auth"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	err := auth.AuthenticateUser()
	if err != nil {
		log.Fatal(err)
	}
}
