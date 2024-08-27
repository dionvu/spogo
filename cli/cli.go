package cli

import (
	"fmt"
	"os"

	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/icons"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/session"
	"github.com/fatih/color"
	"github.com/joomcode/errorx"
	urfave "github.com/urfave/cli/v2"
)

type Cli struct {
	Session *session.Session
	Player  *player.Player
	app     *urfave.App

	CurrCommand string
}

func New(session *session.Session, player *player.Player) *Cli {
	return &Cli{
		Session: session,
		Player:  player,
	}
}

func (c *Cli) SetUp(config *config.Config) error {
	c.app = &urfave.App{
		Commands: []*urfave.Command{
			searchComamnd(c, config),
			deviceCommand(c, config),
			infoCommand(c, config),
			playPauseCommand(c, config),
			nextTrackCommand(c, config),
			prevTrackCommand(c, config),
			volumeCommand(c, config),
			forwardCommand(c, config),
			backwardCommand(c, config),
			shuffleToggleCommand(c, config),
		},

		Name:  "spogo",
		Usage: "control spotify directly in your terminal!",
		OnUsageError: func(cCtx *urfave.Context, err error, isSubcommand bool) error {
			return err
		},

		Action: func(ctx *urfave.Context) error {
			fmt.Printf("%v", ""+
				" ___  ___  ___  ___  ___\n"+
				"|_ -|| . || . || . || . |\n"+
				"|___||  _||___||_  ||___|\n"+
				"     |_|       |___|\n\n")

			fmt.Println(color.HiGreenString(icons.NoteBox + "Spotify " + icons.Multiply + "Go " + icons.Equals + "Spogo!"))
			fmt.Println(color.YellowString(icons.Question + " Help: --help, -h"))

			return nil
		},
	}

	return nil
}

// Runs the urfave after commands have been set up, printing any errors to stdout.
func (c *Cli) Run() {
	if err := c.app.Run(os.Args); err != nil {
		fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), err.(*errorx.Error).Message())

		// fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), err)

		if c.CurrCommand != "" {
			fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+config.APPNAME+" help "+c.CurrCommand))
			return
		}

		fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+config.APPNAME+" help "+os.Args[1]))
	}
}
