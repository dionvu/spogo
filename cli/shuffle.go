package cli

import (
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/errors"
	"github.com/joomcode/errorx"
	urfave "github.com/urfave/cli/v2"
)

func shuffleToggleCommand(cli *Cli, config *config.Config) *urfave.Command {
	return &urfave.Command{
		Name:  "shuffle",
		Usage: "toggle shuffling on current playlist/album",
		Action: func(ctx *urfave.Context) error {
			cli.CurrCommand = ctx.Command.Name

			state, err := cli.Player.State(cli.Session)
			if errorx.GetTypeName(err) == errors.NoDevice.String() {
				errors.Print(err)
				printNoDeviceHelp()
				return nil
			}

			err = cli.Player.Shuffle(!state.ShuffleState, cli.Session)

			errors.Catch(err)

			return nil
		},
		OnUsageError: func(ctx *urfave.Context, err error, isSubcommand bool) error {
			handleBadUsage(ctx, err)
			return nil
		},
	}
}
