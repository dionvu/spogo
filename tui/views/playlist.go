package views

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/spotify"
	"github.com/dionvu/spogo/spotify/auth"
	comp "github.com/dionvu/spogo/tui/views/components"
	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	DEFAULT_PLAYLIST_IMAGE_URL = "https://i.pinimg.com/control/564x/84/29/d1/8429d1c27414bdf99dc5adf9b25a96b3.jpg"
	IMAGES_FOLDER_NAME         = "assets"
	MAX_PLAYLIST_WIDTH         = 35
	MAX_PLAYLIST_ITEM_WIDTH    = 30
	TOP_MARGIN_PLAYLIST        = 7
)

// The view that handles everything related to
// the user's playlist library.
type Playlist struct {
	// The content to be displayed when the view
	// is being displayed.
	Content comp.Content

	// The playlists of the user.
	UserPlaylists *[]spotify.Playlist

	// The list selection for user to select
	// playlists through hovering and selecting.
	PlaylistList PlaylistList

	// The images of all the user's playlists.
	Images   []comp.Image
	imageMap map[list.Item]*comp.Image

	// The detailed information about the selected
	// playlist.
	PlaylistInfo *PlaylistInfo

	// Displays the alternative main views, with the
	// current view (player view) highlighted.
	ViewStatus *ViewStatus

	// Used to access playlists from a selected playlistListItem or
	// selected playlist name.
	playlistsMap map[list.Item]*spotify.Playlist

	Session *auth.Session
	Config  *config.Config
}

// Creates the new playlist view by fetching the user's spotify playlists, determining
// their images to be displayed. Appending all playlists to a bubbletea list.
func NewPlaylistView(s *auth.Session, initialTerm comp.Terminal, cfg *config.Config) *Playlist {
	playlistListItems := []list.Item{}

	pv := &Playlist{
		Images:       []comp.Image{},
		PlaylistInfo: &PlaylistInfo{},
		ViewStatus:   &ViewStatus{},
		playlistsMap: map[list.Item]*spotify.Playlist{},
		imageMap:     map[list.Item]*comp.Image{},
		Session:      s,
		Config:       cfg,
	}

	pv.UserPlaylists, _ = spotify.UserPlaylists(pv.Session)
	if pv.UserPlaylists == nil {
		return pv
	}

	os.MkdirAll(filepath.Join(cfg.CachePath(), IMAGES_FOLDER_NAME), os.ModePerm)

	for i, playlist := range *pv.UserPlaylists {
		pv.Images = append(pv.Images, comp.Image{FilePath: filepath.Join(cfg.CachePath(), IMAGES_FOLDER_NAME, playlist.ID+comp.FILE_EXTENSION)})

		if len(playlist.Images) != 0 {
			pv.Images[i].Update(playlist.Images[0].Url)
		} else {
			pv.Images[i].Update(DEFAULT_PLAYLIST_IMAGE_URL)
		}

		playlistListItem := comp.ListItem(comp.Content(playlist.Name).AdjustFit(MAX_PLAYLIST_ITEM_WIDTH))
		playlistListItems = append(playlistListItems, playlistListItem)

		pv.imageMap[playlistListItem] = &pv.Images[i]
		pv.playlistsMap[playlistListItem] = &playlist
	}

	pv.PlaylistList = PlaylistList{list: comp.NewDefaultList(playlistListItems, "Playlists")}

	return pv
}

// Updates the content to be displayed based on the dimensions of the terminal.
func (pv *Playlist) UpdateContent(term comp.Terminal) {
	mainContainer := comp.NewDefaultTable()
	leftContainer := comp.NewDefaultTable()

	pv.PlaylistInfo.Update(pv.GetSelectedPlaylist())

	vs := ViewStatus{CurrentView: PLAYLIST_VIEW}

	pv.Content = func() comp.Content {
		leftContainer.AppendRow(table.Row{
			comp.Join([]comp.Content{
				pv.PlaylistList.Content().Prepend('\n', 1),
				comp.InvisibleBar(MAX_PLAYLIST_WIDTH),
			}),
		})

		mainContainer.AppendRow(table.Row{
			comp.Content(leftContainer.Render()),
			pv.SelectedImage().AsciiSmall(pv.Config).Content().Append('\n', 1).PadLinesLeft(2),
		})

		if term.WidthIsSmall() || term.HeightIsSmall() {
			return comp.Join([]comp.Content{
				comp.InvisibleBarV(TOP_MARGIN_PLAYLIST),
				comp.Content(mainContainer.Render()).Append('\n', 1).CenterHorizontal(term, -2),
				pv.PlaylistInfo.Content(term).Append('\n', 1).CenterHorizontal(term),
				comp.InvisibleBarV(1),
			}).CenterVertical(term)
		}

		return comp.Join([]comp.Content{
			comp.InvisibleBarV(TOP_MARGIN_PLAYLIST),
			comp.Content(mainContainer.Render()).Append('\n', 2),
			pv.PlaylistInfo.Content(term).Append('\n', 2),
			vs.Content(pv.Config),
		}).CenterVertical(term).CenterHorizontal(term)
	}()
}

// Gets the playlist struct corresponding to the playlist that the user
// is hovering or selected.
func (pv *Playlist) GetSelectedPlaylist() *spotify.Playlist {
	return pv.playlistsMap[pv.PlaylistList.list.SelectedItem()]
}

// The selected playlist's image.
func (pv *Playlist) SelectedImage() *comp.Image {
	return pv.imageMap[pv.PlaylistList.list.SelectedItem()]
}

// Updates the content and renders the view as a string.
func (pv *Playlist) View(playerView *Player, term comp.Terminal) string {
	pv.UpdateContent(term)
	return pv.Content.String()
}

// The detailed infomation about a playlist to
// be displayed when a playlist is selected
// or hovered.
type PlaylistInfo struct {
	Name        PlaylistName
	TotalTracks int
	Owner       string
}

func (pi *PlaylistInfo) Update(playlist *spotify.Playlist) {
	pi.Name = PlaylistName(playlist.Name)
	pi.TotalTracks = playlist.Tracks.Total
	pi.Owner = playlist.Owner.DisplayName
}

// Renders the playlistInfo as a content string.
func (pi PlaylistInfo) Content(term comp.Terminal) comp.Content {
	return comp.Join(
		[]string{
			pi.Name.AdjustFit(term).String(),
			"\nTracks: " + fmt.Sprint(pi.TotalTracks),
		})
}

type PlaylistName string

// Adjusts the playlist string to fit
// within the terminal if it is too big.
func (pn PlaylistName) AdjustFit(term comp.Terminal) PlaylistName {
	c := comp.Content(pn)

	return PlaylistName(c.AdjustFit(term.Width))
}

func (pn PlaylistName) String() string {
	return string(pn)
}

// The list displaying the playlist names.
type PlaylistList struct {
	list     list.Model
	choice   string
	quitting bool
}

func (pll PlaylistList) Update(msg tea.Msg) (PlaylistList, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			pll.quitting = true
			return pll, tea.Quit

		case "esc":
			return pll, nil
		}
	}

	pll.list, cmd = pll.list.Update(msg)

	return pll, cmd
}

// The list as a content string.
func (pl PlaylistList) Content() comp.Content {
	return comp.Content(pl.list.View())
}

func (pl PlaylistList) View() string {
	return pl.list.View()
}

func (_ PlaylistList) Init() tea.Cmd {
	return nil
}
