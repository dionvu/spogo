package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/session"
	"github.com/dionvu/spogo/ui"
)

func main() {
	c, err := config.New()
	errors.Catch(err)
	errors.Catch(c.Load())

	session, err := session.New(c)
	errors.Catch(err)

	player, err := player.New(c)
	errors.Catch(err)

	if !player.GetDevice().IsActive {
		player.Resume(session, false)
		time.Sleep(time.Millisecond * 500)
	}

	// cli := cli.New(session, player)
	//
	// err = cli.SetUp(c)
	// errors.Catch(err)
	//
	// cli.Run()

	initialState, err := player.State(session)
	errors.Catch(err)

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	tp := tea.NewProgram(ui.New(session, player, c, initialState))
	if _, err := tp.Run(); err != nil {
		log.Fatal(err)
	}
}
