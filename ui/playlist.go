package ui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/session"
	"github.com/dionvu/spogo/spotify"
)

type PlaylistView struct {
	Session *session.Session
	Config  *config.Config

	UserPlaylists *[]spotify.Playlist
	PlaylistsMap  map[string]*spotify.Playlist

	ListModel *ListModel
	ItemsMap  map[list.Item]string
}

func NewPlaylistView(s *session.Session, c *config.Config) *PlaylistView {
	items := []list.Item{}

	pv := &PlaylistView{
		Session:  s,
		Config:   c,
		ItemsMap: map[list.Item]string{},
	}

	pv.PlaylistsMap = map[string]*spotify.Playlist{}

	pv.UserPlaylists, _ = spotify.UserPlaylists(pv.Session)

	for _, playlist := range *pv.UserPlaylists {
		items = append(items, Item(playlist.Name))
		pv.ItemsMap[Item(playlist.Name)] = playlist.Name

		pv.PlaylistsMap[playlist.Name] = &playlist
	}

	pv.ListModel = NewListModel(items, TITLE_PLAYLIST_STYLE.Render("Playlists"))

	if len(pv.ListModel.list.Items()) > 0 {
		pv.ListModel.choice = pv.ItemsMap[items[0]]
	}

	return pv
}

func (pv *PlaylistView) View(playerView *PlayerView, terminalSize int) string {
	const DEFAULT_IMAGE_URL string = "https://cdn.pixabay.com/photo/2016/10/22/00/15/spotify-1759471_1280.jpg"
	var res *http.Response

	imagePath := filepath.Join(pv.Config.CachePath(), "image.jpeg")
	imageFile, _ := os.Create(imagePath)

	if len(pv.PlaylistsMap[pv.ItemsMap[pv.ListModel.list.SelectedItem()]].Images) > 0 {
		res, _ = http.Get(pv.PlaylistsMap[pv.ItemsMap[pv.ListModel.list.SelectedItem()]].Images[0].Url)
	} else {
		res, _ = http.Get(DEFAULT_IMAGE_URL)
	}

	io.Copy(imageFile, res.Body)

	if terminalSize <= TERMINALSIZE.Small {
		pv.ListModel.list.SetHeight(SMALL_LIST_HEIGHT)

		if len((*pv.UserPlaylists)) > 1 {
			return fmt.Sprintf("\n\n%s\n\n%s",
				AsciiView(imagePath, ASCII_FLAGS_SMALL),
				pv.ListModel.View())
		}

		return fmt.Sprintf("\n\n%s\n\n%s",
			AsciiView(imagePath, ASCII_FLAGS_SMALL),
			padLines("No playlists :(", TAB_WIDTH))
	}

	pv.ListModel.list.SetHeight(DEFAULT_LIST_HEIGHT)

	if len((*pv.UserPlaylists)) > 1 {
		return fmt.Sprintf("%s\n\n%s\n\n%s", MainControlsView(PLAYLIST_VIEW),
			AsciiView(imagePath, ASCII_FLAGS_NORMAL),
			pv.ListModel.View())
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s", MainControlsView(PLAYLIST_VIEW),
		AsciiView(imagePath, ASCII_FLAGS_NORMAL),
		padLines("No playlists :(", TAB_WIDTH))
}

// if pv.ListModel.choice != "" && len((*pv.UserPlaylists)) > 0 &&
// 	len(pv.PlaylistsMap[pv.ListModel.choice].Images) > 0 &&
// 	len(pv.PlaylistsMap[pv.ListModel.choice].Images) > 0 {
// 	res, _ = http.Get(pv.PlaylistsMap[(pv.ListModel.choice)].Images[0].Url)
// } else if len((*pv.UserPlaylists)) > 0 && len((*pv.UserPlaylists)[0].Images) > 0 {
// 	res, _ = http.Get((*pv.UserPlaylists)[0].Images[0].Url)
// } else {
// }
