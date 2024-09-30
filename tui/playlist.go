package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/spotify"
	"github.com/jedib0t/go-pretty/v6/table"
)

// The view that handles everything related to
// the user's playlist library.
type PlaylistView struct {
	// The content to be displayed when the view
	// is being displayed.
	Content Content

	// The playlists of the user.
	UserPlaylists *[]spotify.Playlist

	// The list selection for user to select
	// playlists through hovering and selecting.
	PlaylistList PlaylistList

	Images []Image

	// The detailed information about the selected
	// playlist.
	PlaylistInfo *PlaylistInfo

	// Displays the alternative main views, with the
	// current view (player view) highlighted.
	ViewStatus *ViewStatus

	// Used to access playlists from a selected item or
	// selected playlist name.
	playlistsMap map[string]*spotify.Playlist
	itemsMap     map[list.Item]string

	Session *auth.Session

	ImageMap map[list.Item]*Image
}

func NewPlaylistView(s *auth.Session) *PlaylistView {
	items := []list.Item{}
	cd, _ := os.UserCacheDir()
	path := filepath.Join(cd, config.APPNAME, "playlist")

	pv := &PlaylistView{
		Session:      s,
		itemsMap:     map[list.Item]string{},
		playlistsMap: map[string]*spotify.Playlist{},

		PlaylistInfo: &PlaylistInfo{},

		ViewStatus: &ViewStatus{},

		Images:   []Image{},
		ImageMap: map[list.Item]*Image{},
	}

	pv.UserPlaylists, _ = spotify.UserPlaylists(pv.Session)

	for i, playlist := range *pv.UserPlaylists {
		pv.Images = append(pv.Images, Image{FilePath: path + fmt.Sprint(i) + ".jpeg"})
		pv.Images[i].Update(playlist.Images[0].Url)
		pv.ImageMap[Item(playlist.Name)] = &pv.Images[i]

		items = append(items, Item(playlist.Name))
		pv.playlistsMap[playlist.Name] = &playlist
		pv.itemsMap[items[len(items)-1]] = playlist.Name
	}

	pv.PlaylistList = NewPlaylistListModel(items, PlaylistViewStyle.Title.Render("Playlists"))

	if len(pv.PlaylistList.list.Items()) > 0 {
		pv.PlaylistList.choice = (*pv.UserPlaylists)[0].Name
	}

	return pv
}

func (pv *PlaylistView) UpdateContent(term Terminal) {
	pv.PlaylistInfo.Update(pv.GetSelectedPlaylist())
	pv.ViewStatus.Update(PLAYLIST_VIEW)

	container := table.NewWriter()
	container.Style().Options.DrawBorder = false
	container.Style().Options.SeparateColumns = false
	// container.Style().Box = table.StyleBoxRounded

	pv.Content = func() Content {
		// if term.IsSizeSmall() {
		pv.PlaylistList.list.SetHeight(LIST_HEIGHT)

		container.AppendRow(table.Row{
			Join([]Content{
				pv.PlaylistList.Content().Prepend('\n', 2).Append('\n', 1),
				pv.PlaylistInfo.Content(term).PadLinesLeft(2),
				Content("").Append(' ', 45),
			}, "\n"),
		})

		t := table.NewWriter()
		t.Style().Options.DrawBorder = false
		t.Style().Options.SeparateColumns = false
		t.AppendRow(table.Row{
			Content(container.Render()).CenterVertical(term),
			pv.ImageMap[pv.PlaylistList.list.SelectedItem()].AsciiNormalBW().Content().Prepend('\n', 1).Append('\n', 1).PadLinesLeft(2).CenterVertical(term),
		})

		return Content(t.Render()).CenterHorizontal(term)
	}()
}

func (pv *PlaylistView) View(playerView *PlayerView, term Terminal) string {
	pv.UpdateContent(term)

	return pv.Content.String()
}

// Gets the playlist object with the same name as what the
// user is selecting.
func (pv *PlaylistView) GetSelectedPlaylist() *spotify.Playlist {
	name := pv.itemsMap[pv.PlaylistList.list.SelectedItem()]
	return pv.playlistsMap[name]
}
