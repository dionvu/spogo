package main

import (
	"log"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/ui"
)

func main() {
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

	tp := tea.NewProgram(ui.New(auth, player, c), tea.WithAltScreen())
	if _, err := tp.Run(); err != nil {
		log.Fatal(err)
	}
}
