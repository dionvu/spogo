package main

import (
	"fmt"
	"os"

	"github.com/dionv/spogo/internal/config"
	"github.com/dionv/spogo/internal/device"
	"github.com/dionv/spogo/internal/player"
	"github.com/dionv/spogo/internal/session"
	"github.com/dionv/spogo/pkg/utils"
	"github.com/fatih/color"
)

func main() {
	c, err := config.New()
	utils.CatchErr(err)

	err = c.Load()
	utils.CatchErr(err)

	s, err := session.New(c)
	utils.CatchErr(err)

	devices, err := device.GetDevices(s)
	utils.CatchErr(err)

	if len(*devices) == 0 {
		fmt.Println(color.RedString("Error"), "No playback devices were found.")
		os.Exit(0)
	}

	p := player.New(&(*devices)[0])

	p.Pause(s)
	// p.Resume(s)
}
