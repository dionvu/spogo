package cli

import (
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/errors"
	"github.com/joomcode/errorx"
	urfave "github.com/urfave/cli/v2"
)

func prevTrackCommand(cli *Cli, config *config.Config) *urfave.Command {
	return &urfave.Command{
		Name:    "previous",
		Aliases: []string{"prev"},
		Usage:   "skips playback to the previous track",
		Action: func(ctx *urfave.Context) error {
			cli.CurrCommand = ctx.Command.Name

			err := cli.Player.SkipPrev(cli.Session)
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
