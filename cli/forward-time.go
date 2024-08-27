package cli

import (
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/errors"
	"github.com/joomcode/errorx"
	urfave "github.com/urfave/cli/v2"
)

func forwardCommand(cli *Cli, config *config.Config) *urfave.Command {
	return &urfave.Command{
		Name:    "forward",
		Aliases: []string{"f"},
		Usage:   "skips current track forward 15 seconds",
		Args:    true,
		Action: func(ctx *urfave.Context) error {
			cli.CurrCommand = ctx.Command.Name

			state, err := cli.Player.State(cli.Session)
			if errorx.GetTypeName(err) == errors.NoDevice.String() {
				errors.Print(err)
				printNoDeviceHelp()
				return nil
			}

			errors.Catch(err)

			errors.Catch(cli.Player.SeekToPosition(cli.Session, state.ProgressMs+15000))

			return nil
		},
		OnUsageError: func(ctx *urfave.Context, err error, isSubcommand bool) error {
			handleBadUsage(ctx, err)
			return nil
		},
	}
}
