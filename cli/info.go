package cli

import (
	"fmt"
	"os"

	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/errors"
	"github.com/joomcode/errorx"
	urfave "github.com/urfave/cli/v2"
)

func infoCommand(cli *Cli, config *config.Config) *urfave.Command {
	return &urfave.Command{
		Name:    "info",
		Aliases: []string{"i"},
		Usage:   "prints info about the current track",
		Action: func(ctx *urfave.Context) error {
			state, err := cli.Player.State(cli.Session)
			if errorx.GetTypeName(err) == errors.NoDevice.String() {
				errors.Print(err)
				printNoDeviceHelp()
				os.Exit(0)
			}

			errors.Catch(err)

			if state.Track == nil {
			} else {
				fmt.Println(state.Track.PlayingInfo(config, state.ProgressMs))
			}

			return nil
		},
		OnUsageError: func(ctx *urfave.Context, err error, isSubcommand bool) error {
			handleBadUsage(ctx, err)
			os.Exit(0)
			return nil
		},
	}
}
