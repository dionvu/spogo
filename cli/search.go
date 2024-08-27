package cli

import (
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/spotify/search"
	"github.com/joomcode/errorx"
	"github.com/manifoldco/promptui"
	urfave "github.com/urfave/cli/v2"
)

func searchComamnd(cli *Cli, config *config.Config) *urfave.Command {
	return &urfave.Command{
		Name:    "search",
		Aliases: []string{"s"},
		Usage:   "searches for `query` string",
		Args:    true,
		Action: func(ctx *urfave.Context) error {
			cli.CurrCommand = "search"

			searchType := []string{"album", "artist", "track", "playlist", "show", "episode"}

			query := ctx.Args().First()

			if query == "" {
				searchPrompt := promptui.Prompt{
					Label: "Enter search query",
				}

				query, err := searchPrompt.Run()
				if err != nil {
					return nil
				}
				if query == "" {
					return errors.Input.New("no search query provided")
				}
			}

			res, err := search.Search(query, searchType, cli.Session)
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
					names = append(names, album.String(config))
				}

				albumPrompt := promptui.Select{
					Label: "Select an album",
					Items: names,
				}

				i, _, err := albumPrompt.Run()
				if err != nil {
					return nil
				}

				err = cli.Player.Play(res.Albums.Items[i].Uri, "", cli.Session)

				if errorx.GetTypeName(err) == errors.NoDevice.String() {
					errors.Print(err)
					printNoDeviceHelp()
					return nil
				}

				errors.Catch(err)

			case "track":

				names := []string{}

				for _, track := range res.Tracks.Items {
					names = append(names, track.String(config))
				}

				albumPrompt := promptui.Select{
					Label: "Select a track",
					Items: names,
				}

				i, _, err := albumPrompt.Run()
				if err != nil {
					return nil
				}

				err = cli.Player.Play("", res.Tracks.Items[i].Uri, cli.Session)

				if errorx.GetTypeName(err) == errors.NoDevice.String() {
					errors.Print(err)
					printNoDeviceHelp()
					return nil
				}

				errors.Catch(err)

			case "playlist":
				names := []string{}

				for _, playlist := range res.Playlists.Items {
					names = append(names, playlist.String())
				}

				albumPrompt := promptui.Select{
					Label: "Select a playlist",
					Items: names,
				}

				i, _, err := albumPrompt.Run()
				if err != nil {
					return nil
				}

				err = cli.Player.Play("", res.Tracks.Items[i].Uri, cli.Session)

				if errorx.GetTypeName(err) == errors.NoDevice.String() {
					errors.Print(err)
					printNoDeviceHelp()
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

				err = cli.Player.Play("", res.Tracks.Items[i].Uri, cli.Session)

				if errorx.GetTypeName(err) == errors.NoDevice.String() {
					errors.Print(err)
					printNoDeviceHelp()
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

				err = cli.Player.Play("", res.Tracks.Items[i].Uri, cli.Session)

				if errorx.GetTypeName(err) == errors.NoDevice.String() {
					errors.Print(err)
					printNoDeviceHelp()
					return nil
				}
			}

			return nil
		},

		OnUsageError: func(ctx *urfave.Context, err error, isSubcommand bool) error {
			handleBadUsage(ctx, err)
			return nil
		},
	}
}
