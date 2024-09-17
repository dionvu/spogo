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
}

func NewPlaylistView(s *session.Session, c *config.Config) *PlaylistView {
	items := []list.Item{}

	pv := &PlaylistView{
		Session: s,
		Config:  c,
	}

	pv.PlaylistsMap = map[string]*spotify.Playlist{}

	pv.UserPlaylists, _ = spotify.UserPlaylists(pv.Session)

	for _, playlist := range *pv.UserPlaylists {
		items = append(items, Item(playlist.Name))

		pv.PlaylistsMap[playlist.Name] = &playlist
	}

	pv.ListModel = NewListModel(items, TITLE_PLAYLIST_STYLE.Render("Playlists"))

	if len((*pv.UserPlaylists)) > 1 {
		pv.ListModel.choice = (*pv.UserPlaylists)[0].Name
	}

	return pv
}

func (pv *PlaylistView) View(playerView *PlayerView, terminalSize int) string {
	const DEFAULT_IMAGE_URL string = "https://cdn.pixabay.com/photo/2016/10/22/00/15/spotify-1759471_1280.jpg"
	var res *http.Response

	imagePath := filepath.Join(pv.Config.CachePath(), "image.jpeg")
	imageFile, _ := os.Create(imagePath)

	if len((*pv.UserPlaylists)) > 0 && len(pv.PlaylistsMap[pv.ListModel.choice].Images) > 0 {
		res, _ = http.Get(pv.PlaylistsMap[(pv.ListModel.choice)].Images[0].Url)
	} else {
		res, _ = http.Get(DEFAULT_IMAGE_URL)
	}

	io.Copy(imageFile, res.Body)

	if len((*pv.UserPlaylists)) > 1 {
		return fmt.Sprintf("%s\n\n%s\n\n%s", MainControlsView(PLAYLIST_VIEW), AsciiView(imagePath, ASCII_FLAGS_NORMAL), pv.ListModel.View())
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s", MainControlsView(PLAYLIST_VIEW), AsciiView(imagePath, ASCII_FLAGS_NORMAL), padLines("No playlists :(", TAB_WIDTH))
}
