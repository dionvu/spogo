package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joomcode/errorx"
	"github.com/urfave/cli/v2"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/config"
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

	session, err := session.New(c)
	utils.CatchErr(err)

	p, err := player.New(c)
	utils.CatchErr(err)

	app := &cli.App{
		EnableBashCompletion: true,

		Name:  "spogo",
		Usage: "spotify + go = spogo!",
		Action: func(ctx *cli.Context) error {
			printWelcome()
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "devices",
				Aliases: []string{"d"},
				Usage:   "Select a playback device",
				Action: func(ctx *cli.Context) error {
					d, err := p.UserSelectDevice(session, c)
					if errorx.GetTypeName(err) == errors.NoDeviceError.String() {
						HandleNoDevice()
						os.Exit(0)
					}

					if err = p.SetDevice(d, c); err != nil {
						return err
					}

					return nil
				},
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					PrintHelpCommand(ctx, err)
					os.Exit(0)
					return nil
				},
			},
			{
				Name:    "resume",
				Aliases: []string{"r"},
				Usage:   "resume playback",
				Action: func(ctx *cli.Context) error {
					if err := p.Resume(session); err != nil {
						HandleNoDevice()
					}

					return nil
				},
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					PrintHelpCommand(ctx, err)
					os.Exit(0)
					return nil
				},
			},
			{
				Name:    "pause",
				Aliases: []string{"p"},
				Usage:   "pause playback",
				Action: func(ctx *cli.Context) error {
					if err := p.Pause(session); err != nil {
						HandleNoDevice()
						os.Exit(0)
					}

					return nil
				},
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					PrintHelpCommand(ctx, err)
					os.Exit(0)
					return nil
				},
			},
			{
				Name:    "next",
				Aliases: []string{"n"},
				Usage:   "skips playback to next track in queue",
				Action: func(ctx *cli.Context) error {
					if err := p.SkipNext(session); err != nil {
						HandleNoDevice()
					}

					return nil
				},
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					PrintHelpCommand(ctx, err)
					os.Exit(0)
					return nil
				},
			},

			{
				Name:    "previous",
				Aliases: []string{"prev", "back", "b"},
				Usage:   "skips playback to the previous track",
				Action: func(ctx *cli.Context) error {
					if err := p.SkipPrev(session); err != nil {
						HandleNoDevice()
					}
					return nil
				},
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					PrintHelpCommand(ctx, err)
					os.Exit(0)
					return nil
				},
			},

			{
				Name:    "volume",
				Aliases: []string{"v", "vol"},
				Usage:   "change playback device volume",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "set",
						Aliases: []string{"s"},
						Usage:   "sets volume to `INT`",
						Value:   0,
					},
					&cli.IntFlag{
						Name:    "up",
						Aliases: []string{"u"},
						Usage:   "raises volume by `INT`",
					},
					&cli.IntFlag{
						Name:    "down",
						Aliases: []string{"d"},
						Usage:   "lowers volume by `INT`",
					},
				},

				Action: func(ctx *cli.Context) error {
					if ctx.IsSet("set") {
						// Ensures volume is in the range [0, 100], and sets volume.
						vol := int(max(0, min(100, ctx.Int("set"))))

						if err := p.SetVolume(session, vol); err != nil {
							HandleNoDevice()
							os.Exit(0)
						}
						return nil
					}

					if ctx.IsSet("up") {
					}

					if ctx.IsSet("down") {
					}

					HandleNoFlag(ctx)

					return nil
				},

				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					PrintHelpCommand(ctx, err)
					os.Exit(0)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), err)
		fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+config.APPNAME+" help <command>"))
	}
}

// Prints the command to pritn help information
// corresponding to the command that the user messed up on.
func PrintHelpCommand(ctx *cli.Context, err error) {
	fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), err)
	fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+ctx.App.Name+" help "+ctx.Command.Name))
}

func HandleNoFlag(ctx *cli.Context) {
	fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), "no flags provided")
	fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+config.APPNAME+" help "+ctx.Command.Name))
}

// Prints the error message corresponding to no
// active or selected playback device.
func HandleNoDevice() {
	fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), "no active playback device")
	fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+config.APPNAME+" devices"))
}

func printWelcome() {
	fmt.Println(cli.HelpFlag)
}
