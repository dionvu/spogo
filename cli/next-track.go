package cli

import (
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/errors"
	"github.com/joomcode/errorx"
	urfave "github.com/urfave/cli/v2"
)

func nextTrackCommand(cli *Cli, config *config.Config) *urfave.Command {
	return &urfave.Command{
		Name:    "next",
		Aliases: []string{"n"},
		Usage:   "skips playback to next track in queue",
		Action: func(ctx *urfave.Context) error {
			cli.CurrCommand = ctx.Command.Name

			err := cli.Player.SkipNext(cli.Session)
			if errorx.GetTypeName(err) == errors.NoDevice.String() {
				errors.Print(err)
				printNoDeviceHelp()
				return nil
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
