package cli

import (
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/errors"
	"github.com/joomcode/errorx"
	urfave "github.com/urfave/cli/v2"
)

func playPauseCommand(cli *Cli, config *config.Config) *urfave.Command {
	return &urfave.Command{
		Name:    "play/pause",
		Aliases: []string{"p"},
		Usage:   "toggles playback",
		Action: func(ctx *urfave.Context) error {
			cli.CurrCommand = ctx.Command.Name

			state, err := cli.Player.State(cli.Session)

			// If state returns a no device, there is a chance
			// user opened a device and hasn't activated it yet.
			if errorx.GetTypeName(err) == errors.NoDevice.String() {

				err = cli.Player.Resume(cli.Session)

				// If attempt to resume failed, user has no device selected.
				if errorx.GetTypeName(err) == errors.NoDevice.String() {
					errors.Print(err)
					printNoDeviceHelp()

					return nil
				}

				// Successful started playback on newly device.
				return nil
			}

			errors.Catch(err)

			if state.IsPlaying {
				err = cli.Player.Pause(cli.Session)
			} else {
				err = cli.Player.Resume(cli.Session)
			}

			errors.Catch(err)
			return nil
		},

		OnUsageError: func(ctx *urfave.Context, err error, isSubcommand bool) error {
			handleBadUsage(ctx, err)

			return nil
		},
	}
}
