package main

import (
	"fmt"
	"os"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/config"
	"github.com/dionv/spogo/internal/device"
	"github.com/dionv/spogo/internal/session"
	"github.com/fatih/color"
	"github.com/joomcode/errorx"

	// "github.com/dionv/spogo/internal/device"
	"github.com/dionv/spogo/internal/player"
	"github.com/dionv/spogo/pkg/utils"
	// "github.com/fatih/color"
)

func main() {
	c, err := config.New()
	utils.CatchErr(err)

	err = c.Load()
	utils.CatchErr(err)

	s, err := session.New(c)
	utils.CatchErr(err)

	p, err := player.New(c)
	utils.CatchErr(err)

	devices, err := device.GetDevices(s)
	utils.CatchErr(err)

	if len(*devices) == 0 {
		fmt.Println(color.RedString("Error"), "No playback devices were found.")
		os.Exit(0)
	}

	err = p.SetDevice(&(*devices)[0], c)
	if err != nil {
		fmt.Println(err)
	}

	// handleErrorControls(e)
	err = p.Resume(s)

	// p.Pause(s)

	// p.SkipNext(s)
}

func handleErrorControls(err error) {
	if errorx.GetTypeName(err) == errors.ReauthenticationError.String() {
		fmt.Println("Please reauth")
		os.Exit(0)
	}

	utils.CatchErr(err)
}
