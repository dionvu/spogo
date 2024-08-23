package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joomcode/errorx"
	"github.com/manifoldco/promptui"

	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/device"
	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/icons"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/session"
	"github.com/dionvu/spogo/spotify/search"
	"github.com/urfave/cli/v2"
)

func main() {
	var commandName string

	c, err := config.New()
	errors.Catch(err)
	errors.Catch(c.Load())

	session, err := session.New(c)
	errors.Catch(err)

	player, err := player.New(c)
	errors.Catch(err)

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "search",
				Aliases: []string{"s"},
				Usage:   "searches for `query` string",
				Args:    true,
				Action: func(ctx *cli.Context) error {
					commandName = "search"

					searchType := []string{"album", "artist", "track", "playlist", "show", "episode"}

					query := ctx.Args().First()

					if query == "" {
						searchPrompt := promptui.Prompt{
							Label: "Enter search query",
						}

						query, err = searchPrompt.Run()
						if err != nil {
							return nil
						}
						if query == "" {
							return errors.Input.New("no search query provided")
						}
					}

					res, err := search.Search(query, searchType, session)
					if err != nil {
						return err
					}

					categoryPrompt := promptui.Select{
						Label: "Select a category",
						Items: searchType,
					}

					_, category, err := categoryPrompt.Run()

					switch category {
					case "album":
						names := []string{}

						for _, album := range res.Albums.Items {
							names = append(names, album.String())
						}

						albumPrompt := promptui.Select{
							Label: "Select an album",
							Items: names,
						}

						i, _, err := albumPrompt.Run()
						if err != nil {
							return nil
						}

						err = player.Play(res.Albums.Items[i].Uri, "", session)

						if errorx.GetTypeName(err) == errors.NoDevice.String() {
							errors.Print(err)
							PrintHelpCommand(ctx.Command)
							return nil
						}

						errors.Catch(err)

					case "track":

						names := []string{}

						for _, track := range res.Tracks.Items {
							names = append(names, track.String())
						}

						albumPrompt := promptui.Select{
							Label: "Select a track",
							Items: names,
						}

						i, _, err := albumPrompt.Run()
						if err != nil {
							return nil
						}

						err = player.Play("", res.Tracks.Items[i].Uri, session)

						if errorx.GetTypeName(err) == errors.NoDevice.String() {
							errors.Print(err)
							PrintHelpCommand(ctx.Command)
							return nil
						}

						errors.Catch(err)

					case "playlist":
						names := []string{}

						for _, playlist := range res.Playlists.Items {
							names = append(names, playlist.Name+" | "+playlist.Owner.DisplayName)
						}

						albumPrompt := promptui.Select{
							Label: "Select a playlist",
							Items: names,
						}

						i, _, err := albumPrompt.Run()
						if err != nil {
							return nil
						}

						err = player.Play("", res.Tracks.Items[i].Uri, session)

						if errorx.GetTypeName(err) == errors.NoDevice.String() {
							errors.Print(err)
							PrintHelpCommand(ctx.Command)
							return nil
						}
					case "show":
						names := []string{}

						for _, show := range res.Shows.Items {
							names = append(names, show.String())
						}

						albumPrompt := promptui.Select{
							Label: "Select a playlist",
							Items: names,
						}

						i, _, err := albumPrompt.Run()
						if err != nil {
							return nil
						}

						err = player.Play("", res.Tracks.Items[i].Uri, session)

						if errorx.GetTypeName(err) == errors.NoDevice.String() {
							errors.Print(err)
							PrintHelpCommand(ctx.Command)
							return nil
						}
					case "episode":
						names := []string{}

						for _, ep := range res.Episodes.Items {
							names = append(names, ep.String())
						}

						episodePrompt := promptui.Select{
							Label: "Select an episode",
							Items: names,
						}

						i, _, err := episodePrompt.Run()
						if err != nil {
							return nil
						}

						err = player.Play("", res.Tracks.Items[i].Uri, session)

						if errorx.GetTypeName(err) == errors.NoDevice.String() {
							errors.Print(err)
							PrintHelpCommand(ctx.Command)
							return nil
						}
					}

					return nil
				},

				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					HandleBadUsage(ctx, err)
					return nil
				},
			},
			{
				Name:    "devices",
				Aliases: []string{"d"},
				Usage:   "Select a playback device",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "detailed",
						Aliases: []string{"d"},
						Usage:   "addittional detailed information about each device",
					},
				},
				Action: func(ctx *cli.Context) error {
					var d *device.Device
					var err error

					if ctx.IsSet("detailed") {
						d, err = player.UserSelectDevice(session, c, true)
					} else {
						d, err = player.UserSelectDevice(session, c, false)
					}

					if errorx.GetTypeName(err) == errors.NoDevice.String() {
						errors.Print(err)
						PrintHelpCommand(ctx.Command)
						return nil
					}

					errors.Catch(player.SetDevice(d, c))

					return nil
				},
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					HandleBadUsage(ctx, err)
					return nil
				},
			},
			{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "prints info about the current track",
				Action: func(ctx *cli.Context) error {
					state, err := player.State(session)
					if errorx.GetTypeName(err) == errors.NoDevice.String() {
						errors.Print(err)
						PrintNoDevice()
						os.Exit(0)
					}

					errors.Catch(err)

					if state.Track == nil {
					} else {
						fmt.Println(state.Track.StringPlaying(state.ProgressMs))
					}

					return nil
				},
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					HandleBadUsage(ctx, err)
					os.Exit(0)
					return nil
				},
			},
			{
				Name:    "play/pause",
				Aliases: []string{"p"},
				Usage:   "toggles playback",
				Action: func(ctx *cli.Context) error {
					commandName = ctx.Command.Name

					state, err := player.State(session)
					if errorx.GetTypeName(err) == errors.NoDevice.String() {
						errors.Print(err)
						PrintHelpCommand(ctx.Command)
						os.Exit(0)
					}

					errors.Catch(err)

					if state.IsPlaying {
						err = player.Pause(session)
					} else {
						err = player.Resume(session)
					}

					errors.Catch(err)
					return nil
				},
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					HandleBadUsage(ctx, err)
					os.Exit(0)
					return nil
				},
			},
			{
				Name:    "next",
				Aliases: []string{"n"},
				Usage:   "skips playback to next track in queue",
				Action: func(ctx *cli.Context) error {
					commandName = ctx.Command.Name

					err := player.SkipNext(session)
					if errorx.GetTypeName(err) == errors.NoDevice.String() {
						errors.Print(err)
						PrintHelpCommand(ctx.Command)
						return nil
					}

					errors.Catch(err)

					return nil
				},
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					HandleBadUsage(ctx, err)
					os.Exit(0)
					return nil
				},
			},

			{
				Name:    "previous",
				Aliases: []string{"prev"},
				Usage:   "skips playback to the previous track",
				Action: func(ctx *cli.Context) error {
					commandName = ctx.Command.Name

					err := player.SkipPrev(session)
					if errorx.GetTypeName(err) == errors.NoDevice.String() {
						errors.Print(err)
						PrintHelpCommand(ctx.Command)
						return nil
					}
					errors.Catch(err)

					return nil
				},
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					HandleBadUsage(ctx, err)
					os.Exit(0)
					return nil
				},
			},
			{
				Name:    "volume",
				Aliases: []string{"v", "vol"},
				Usage:   "change playback device volume [0-100]",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "set",
						Aliases: []string{"s"},
						Usage:   "sets volume to `num`",
						Value:   0,
					},
					&cli.IntFlag{
						Name:    "up",
						Aliases: []string{"u"},
						Usage:   "raises volume by `num`",
					},
					&cli.IntFlag{
						Name:    "down",
						Aliases: []string{"d"},
						Usage:   "lowers volume by `num`",
					},
				},

				Action: func(ctx *cli.Context) error {
					if ctx.IsSet("set") {
						// Ensures volume is in the range [0, 100], and sets volume.
						vol := int(max(0, min(100, ctx.Int("set"))))

						err := player.SetVolume(session, vol)
						if errorx.GetTypeName(err) == errors.NoDevice.String() {
							errors.Print(err)
							PrintHelpCommand(ctx.Command)
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

					HandleNoFlag(ctx)

					return nil
				},

				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					HandleBadUsage(ctx, err)
					os.Exit(0)
					return nil
				},
			},

			{
				Name:    "forward",
				Aliases: []string{"f"},
				Usage:   "skips current track forward 15 seconds",
				Args:    true,
				Action: func(ctx *cli.Context) error {
					commandName = ctx.Command.Name

					state, err := player.State(session)
					if errorx.GetTypeName(err) == errors.NoDevice.String() {
						errors.Print(err)
						PrintHelpCommand(ctx.Command)
						return nil
					}

					errors.Catch(err)

					errors.Catch(player.SeekToPosition(session, state.ProgressMs+15000))

					return nil
				},
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					HandleBadUsage(ctx, err)
					os.Exit(0)
					return nil
				},
			},
			{
				Name:    "backward",
				Aliases: []string{"b"},
				Usage:   "skips current track backward 15 seconds",
				Action: func(ctx *cli.Context) error {
					commandName = ctx.Command.Name

					state, err := player.State(session)
					if errorx.GetTypeName(err) == errors.NoDevice.String() {
						errors.Print(err)
						PrintHelpCommand(ctx.Command)
						return nil
					}

					errors.Catch(err)

					errors.Catch(player.SeekToPosition(session, state.ProgressMs-15000))

					return nil
				},
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					HandleBadUsage(ctx, err)
					return nil
				},
			},
			{
				Name:  "shuffle",
				Usage: "toggle shuffling on current playlist/album",
				Action: func(ctx *cli.Context) error {
					commandName = ctx.Command.Name

					state, err := player.State(session)
					if errorx.GetTypeName(err) == errors.NoDevice.String() {
						errors.Print(err)
						PrintHelpCommand(ctx.Command)
						return nil
					}

					err = player.Shuffle(!state.ShuffleState, session)

					errors.Catch(err)

					return nil
				},
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					HandleBadUsage(ctx, err)
					return nil
				},
			},
		},

		Name:  "spogo",
		Usage: "control spotify directly in your terminal!",
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			return err
		},

		Action: func(ctx *cli.Context) error {
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

	// Runs the cli command and catches any error.
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), err.(*errorx.Error).Message())
		// fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), err)
		if commandName != "" {
			fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+config.APPNAME+" help "+commandName))
			return
		}

		fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+config.APPNAME+" help "+os.Args[1]))
	}
}

// Prints the command to print help information
// corresponding to the command that the user messed up on.
func HandleBadUsage(ctx *cli.Context, err error) {
	fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), err)
	fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+ctx.App.Name+" help "+ctx.Command.Name))
}

// Prints the error message corresponding to
// command requiring flags but none were provided.
func HandleNoFlag(ctx *cli.Context) {
	fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), "no flags provided")
	fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+config.APPNAME+" help "+ctx.Command.Name))
}

// Prints the help command corresponding to given command.
func PrintHelpCommand(c *cli.Command) {
	fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+config.APPNAME+" "+c.Name))
}

func PrintNoDevice() {
	fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+config.APPNAME+" device"))
}
