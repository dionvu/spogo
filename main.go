package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"github.com/dionv/spogo/internal/config"
	"github.com/dionv/spogo/internal/device"
	"github.com/dionv/spogo/internal/player"
	"github.com/dionv/spogo/internal/session"
	"github.com/dionv/spogo/pkg/utils"
	"github.com/dionv/spogo/public/icons"
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

	////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////

	app := &cli.App{
		EnableBashCompletion: true,

		Name:  "spogo",
		Usage: "spotify + go = spogo!",
		Action: func(ctx *cli.Context) error {
			// fmt.Println(cli.HelpFlag)

			return nil
		},

		Commands: []*cli.Command{
			{
				Name:    "resume",
				Aliases: []string{"r"},
				Usage:   "resume playback",
				Action: func(ctx *cli.Context) error {
					if err := p.Resume(s); err != nil {
						fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), "No active playback device detected")
						os.Exit(0)
					}

					return nil
				},
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), err)
					fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+ctx.App.Name+" help "+ctx.Command.Name))
					return nil
				},
			},
			{
				Name:    "pause",
				Aliases: []string{"p"},
				Usage:   "resume playback",
				Action: func(ctx *cli.Context) error {
					if err := p.Pause(s); err != nil {
						fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), "No active playback device detected")
						os.Exit(0)
					}

					return nil
				},
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), err)
					fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+ctx.App.Name+" help "+ctx.Command.Name))
					return nil
				},
			},
			{
				Name:    "volume",
				Aliases: []string{"v", "vol"},
				Usage:   "change playback device volume",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "set",
						Usage: "sets volume to given `NUMBER`",
					},
				},

				Action: func(ctx *cli.Context) error {
					if ctx.Args().First() == "" {
						fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), "No interger provided")
						fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+ctx.App.Name+" help "+ctx.Command.Name))
						return errors.New("No provided arguements")
					}

					num, err := strconv.ParseInt(ctx.Args().First(), 10, 8)
					if err != nil {
						fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), "Invalid integer")
						fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+ctx.App.Name+" help "+ctx.Command.Name))
						os.Exit(0)
					}

					v := max(0, min(100, num))

					if err := p.SetVolume(s, int(v)); err != nil {
						fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), "No active playback device detected")
						os.Exit(0)
					}

					return nil
				},
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), err)
					fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+ctx.App.Name+" help "+ctx.Command.Name))
					return nil
				},
			},
		},
	}

	// defd(p, s, c)

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), err)
		fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+config.APPNAME+" help <command>"))
	}
}

func defd(p *player.Player, s *session.Session, c *config.Config) {
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
}
