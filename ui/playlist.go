package ui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/spotify"
)

type PlaylistView struct {
	Session *auth.Session
	Config  *config.Config

	UserPlaylists *[]spotify.Playlist
	playlistsMap  map[string]*spotify.Playlist

	ItemsMap map[list.Item]string

	PlaylistListModel *PlaylistListModel
}

func NewPlaylistView(s *auth.Session, c *config.Config) *PlaylistView {
	items := []list.Item{}

	pv := &PlaylistView{
		Session: s,
		Config:  c,
	}

	pv.ItemsMap = map[list.Item]string{}

	pv.playlistsMap = map[string]*spotify.Playlist{}

	pv.UserPlaylists, _ = spotify.UserPlaylists(pv.Session)

	for _, playlist := range *pv.UserPlaylists {
		item := Item(playlist.Name)

		items = append(items, item)

		pv.playlistsMap[playlist.Name] = &playlist

		pv.ItemsMap[items[len(items)-1]] = playlist.Name
	}

	pv.PlaylistListModel = NewPlaylistListModel(items, PlaylistViewStyle.Title.Render("Playlists"))

	if len(pv.PlaylistListModel.list.Items()) > 0 {
		pv.PlaylistListModel.choice = (*pv.UserPlaylists)[0].Name
	}

	return pv
}

func (pv *PlaylistView) View(playerView *PlayerView, terminalSize int) string {
	const DEFAULT_IMAGE_URL string = "https://cdn.pixabay.com/photo/2016/10/22/00/15/spotify-1759471_1280.jpg"
	var res *http.Response

	imagePath := filepath.Join(pv.Config.CachePath(), "image.jpeg")
	imageFile, _ := os.Create(imagePath)

	if len(pv.playlistsMap) > 0 && len(pv.playlistsMap[pv.PlaylistListModel.choice].Images) > 0 {
		res, _ = http.Get(pv.GetSelectedPlaylist().Images[0].Url)
	} else {
		res, _ = http.Get(DEFAULT_IMAGE_URL)
	}

	io.Copy(imageFile, res.Body)

	if terminalSize <= TERMINALSIZE.Small {
		pv.PlaylistListModel.list.SetHeight(SMALL_LIST_HEIGHT)

		if len((*pv.UserPlaylists)) > 0 {
			return fmt.Sprintf("\n\n%s\n\n%s",
				AsciiView(imagePath, ASCII_FLAGS_SMALL),
				pv.PlaylistListModel.View())
		}

		return fmt.Sprintf("\n\n%s\n\n%s",
			AsciiView(imagePath, ASCII_FLAGS_SMALL),
			padLines("No playlists :(", TAB_WIDTH))
	}

	pv.PlaylistListModel.list.SetHeight(DEFAULT_LIST_HEIGHT)

	if len((*pv.UserPlaylists)) > 0 {
		return fmt.Sprintf("\n\n%s\n\n%s\n\n%s", MainControlsView(PLAYLIST_VIEW),
			AsciiView(imagePath, ASCII_FLAGS_NORMAL),
			pv.PlaylistListModel.View())
	}

	return fmt.Sprintf("\n\n%s\n\n%s\n\n%s", MainControlsView(PLAYLIST_VIEW),
		AsciiView(imagePath, ASCII_FLAGS_NORMAL),
		padLines("No playlists :(", TAB_WIDTH))
}

func (pv *PlaylistView) GetPlaylistFromChoice(choice string) *spotify.Playlist {
	return pv.playlistsMap[choice]
}

func (pv *PlaylistView) GetSelectedName() string {
	return pv.ItemsMap[pv.PlaylistListModel.list.SelectedItem()]
}

func (pv *PlaylistView) GetSelectedPlaylist() *spotify.Playlist {
	return pv.GetPlaylistFromChoice(pv.GetSelectedName())
}
