package cli

import (
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/errors"
	"github.com/joomcode/errorx"
	urfave "github.com/urfave/cli/v2"
)

func volumeCommand(cli *Cli, config *config.Config) *urfave.Command {
	return &urfave.Command{
		Name:    "volume",
		Aliases: []string{"v", "vol"},
		Usage:   "change playback device volume [0-100]",
		Flags: []urfave.Flag{
			&urfave.IntFlag{
				Name:    "set",
				Aliases: []string{"s"},
				Usage:   "sets volume to `num`",
				Value:   0,
			},
			&urfave.IntFlag{
				Name:    "up",
				Aliases: []string{"u"},
				Usage:   "raises volume by `num`",
			},
			&urfave.IntFlag{
				Name:    "down",
				Aliases: []string{"d"},
				Usage:   "lowers volume by `num`",
			},
		},

		Action: func(ctx *urfave.Context) error {
			if ctx.IsSet("set") {
				// Ensures volume is in the range [0, 100], and sets volume.
				vol := int(max(0, min(100, ctx.Int("set"))))

				err := cli.Player.SetVolume(cli.Session, vol)
				if errorx.GetTypeName(err) == errors.NoDevice.String() {
					errors.Print(err)
					printNoDeviceHelp()
					return nil
				}

				errors.Catch(err)

				return nil
			}

			if ctx.IsSet("up") {
				// TODO
			}

			if ctx.IsSet("down") {
				// TODO
			}

			handleNoFlag(ctx)

			return nil
		},

		OnUsageError: func(ctx *urfave.Context, err error, isSubcommand bool) error {
			handleBadUsage(ctx, err)
			return nil
		},
	}
}
