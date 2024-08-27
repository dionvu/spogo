package cli

import (
	"fmt"

	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/device"
	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/icons"
	"github.com/fatih/color"
	"github.com/joomcode/errorx"
	urfave "github.com/urfave/cli/v2"
)

func deviceCommand(cli *Cli, config *config.Config) *urfave.Command {
	return &urfave.Command{
		Name:    "devices",
		Aliases: []string{"d"},
		Usage:   "Select a playback device",
		Flags: []urfave.Flag{
			&urfave.BoolFlag{
				Name:    "detailed",
				Aliases: []string{"d"},
				Usage:   "addittional detailed information about each device",
			},
		},
		Action: func(ctx *urfave.Context) error {
			var d *device.Device
			var err error

			if ctx.IsSet("detailed") {
				d, err = cli.Player.UserSelectDevice(cli.Session, config, true)
			} else {
				d, err = cli.Player.UserSelectDevice(cli.Session, config, false)
			}

			if errorx.GetTypeName(err) == errors.NoDevice.String() {
				fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), "no playback devices open or detected")
				printNoDeviceHelp()

				return nil
			}

			errors.Catch(cli.Player.SetDevice(d, config))

			return nil
		},
		OnUsageError: func(ctx *urfave.Context, err error, isSubcommand bool) error {
			handleBadUsage(ctx, err)
			return nil
		},
	}
}
