package main

import (
	"fmt"
	"log"

	"github.com/dionv/spogo/internal/app"
)

// "log"
//
// "github.com/dionv/spogo/internal/auth"
// "github.com/joho/godotenv"

func main() {
	// godotenv.Load()
	//
	// err := auth.EnsureAuthenticated()
	// if err != nil {
	// 	log.Fatal(err)
	// 	// log.Fatalf("%+v", err)
	// }

	appConfig := app.New()

	configExists, e := appConfig.Exists()

	if e != nil {
		log.Fatalf("%+v", e)
	}

	if !configExists {
		log.Println("new cfg")
		appConfig.Create()
	}

	err := appConfig.Load()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(appConfig.Spotify.ClientID)
	fmt.Println(appConfig.Spotify.ClientSecret)
}
