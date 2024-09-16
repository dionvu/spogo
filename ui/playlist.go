package ui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	"github.com/charmbracelet/bubbles/list"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/session"
	"github.com/dionvu/spogo/spotify"
)

type PlaylistView struct {
	Session *session.Session
	Player  *player.Player
	Config  *config.Config
	State   *player.PlayerState

	UserPlaylists *[]spotify.Playlist
	PlaylistsMap  map[string]*spotify.Playlist

	List *ListModel
}

func NewPlaylistView(s *session.Session) *PlaylistView {
	pv := &PlaylistView{
		Session: s,
	}

	pv.UserPlaylists, _ = spotify.UserPlaylists(pv.Session)

	items := []list.Item{}
	pv.PlaylistsMap = map[string]*spotify.Playlist{}

	for _, playlist := range *pv.UserPlaylists {
		items = append(items, Item(playlist.Name+"\n"))
		pv.PlaylistsMap[playlist.Name] = &playlist
	}

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = ""
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	pv.List = &ListModel{list: l}

	return pv
}

func (pv *PlaylistView) View(playerView *PlayerView) string {
	const DEFAULT_IMAGE_URL string = "https://cdn.pixabay.com/photo/2016/10/22/00/15/spotify-1759471_1280.jpg"
	mainControls := MAIN_CONTROLS_STYLE.Render("[ F1 Player | ") + MAIN_CONTROLS_SELECTED_STYLE.Render("F2 Playlists") + MAIN_CONTROLS_STYLE.Render(" | F3 Search | F4 Devices ]")
	title := TITLE_PLAYLIST_STYLE.Render("Playlists")

	var ascii string

	pv.List.list.Title = title

	if pv.List.choice == "" {
		pv.List.choice = (*pv.UserPlaylists)[0].Name
	}

	if len(pv.PlaylistsMap[strings.TrimSpace(pv.List.choice)].Images) > 1 {
		res, _ := http.Get(pv.PlaylistsMap[strings.TrimSpace(pv.List.choice)].Images[0].Url)

		cd, _ := os.UserCacheDir()
		filepath := filepath.Join(cd, config.APPNAME, "image.jpeg")

		file, _ := os.Create(filepath)

		io.Copy(file, res.Body)

		flags := aic_package.DefaultFlags()
		flags.Colored = true
		flags.Dimensions = []int{40, 20}
		flags.Braille = true

		ascii, _ = aic_package.Convert(filepath, flags)
	} else {
		res, _ := http.Get(DEFAULT_IMAGE_URL)

		cd, _ := os.UserCacheDir()
		filepath := filepath.Join(cd, config.APPNAME, "image.jpeg")

		file, _ := os.Create(filepath)

		io.Copy(file, res.Body)

		flags := aic_package.DefaultFlags()
		flags.Colored = true
		flags.Dimensions = []int{40, 20}
		flags.Braille = true

		ascii, _ = aic_package.Convert(filepath, flags)
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s", mainControls, ascii, pv.List.list.View())
}
