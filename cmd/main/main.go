package main

import (
	"fmt"
	"path/filepath"

	"github.com/dionv/spogo/internal/app"
	"github.com/dionv/spogo/internal/app/config"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app, err := app.New()
	c := app.Config
	user := app.User()

	// configExists, e := app.Config.Exists()
	// if e != nil {
	// 	log.Fatalf("%+v", e)
	// }
	//
	// if !configExists {
	// 	// Log
	// 	log.Println("new cfg")
	// 	c.Create()
	// }

	pathAT := filepath.Join(c.Path(), config.ACCESSTOKENFILE)
	pathRT := filepath.Join(c.Path(), config.REQUESTTOKENFILE)

	err = c.Load()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	err = user.Authenticate(c)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	accessToken := user.AccessToken
	err = accessToken.Load(pathAT)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Println("access token")
	fmt.Println(user.AccessToken.String())

	refreshTok := app.User().RefreshToken
	err = refreshTok.Load(pathRT)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Println("refresh token")
	fmt.Println(user.RefreshToken.String())

	// err = accessToken.Refresh(&refreshTok, c.Spotify.ClientID(), c.Spotify.ClientSecret())
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//
	// fmt.Println(accessToken.String())
	//
	// path, _ = c.Root()
	// path += "/access-token.json"
	//
	// accessToken.Update(path)
}
