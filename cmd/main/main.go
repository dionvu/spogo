package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dionv/spogo/internal/config"
	"github.com/dionv/spogo/internal/user"
	"github.com/fatih/color"
)

func main() {
	c, err := config.New()
	if err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error"), err)
	}

	u := user.New()

	// app := app.New(c, u)
	if err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error"), err)
	}

	err = c.Load()
	if err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error"), err)
	}

	u.AccessToken.Load(filepath.Join(c.Path(), config.TOKENSDIRECTORY, config.ACCESSTOKENFILE))
	u.RefreshToken.Load(filepath.Join(c.Path(), config.TOKENSDIRECTORY, config.REQUESTTOKENFILE))

	// Reauthentication every 60 minutes (lifecycle of an access token)
	if time.Since(u.AccessToken.TimeCreated) > 59*time.Minute {
		valid, _ := c.ValidSpotifyCredentials()
		if !valid {
			fmt.Printf("%v Invalid spotify client ID & client secret: %v\n",
				color.RedString("Error"),
				color.GreenString(c.FilePath()))

			os.Exit(0)
		}

		err := u.Authenticate(c)
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
	}
}
