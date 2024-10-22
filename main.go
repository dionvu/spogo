package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/err"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/spotify/auth"
	"github.com/dionvu/spogo/tui"
)

func main() {
	errors.Init()

	c, err := config.New()
	errors.Catch(err)
	errors.Catch(c.Load())

	auth, err := auth.New(c)
	errors.Catch(err)

	player, err := player.New(c)
	errors.Catch(err)

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	program := tui.New(auth, player, c)
	if err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
