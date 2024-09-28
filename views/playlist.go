package ui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/spotify"
	"github.com/jedib0t/go-pretty/v6/table"
)

type PlaylistView struct {
	Session *auth.Session
	Config  *config.Config

	UserPlaylists *[]spotify.Playlist
	playlistsMap  map[string]*spotify.Playlist

	ItemsMap map[list.Item]string

	PlaylistListModel *PlaylistListModel
}

type PlaylistListModel struct {
	list     list.Model
	choice   string
	quitting bool
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

func (pv *PlaylistView) View(playerView *PlayerView, terminal Terminal) string {
	pv.PlaylistListModel.list.SetHeight(DEFAULT_LIST_HEIGHT)

	const DEFAULT_IMAGE_URL string = "https://cdn.pixabay.com/photo/2016/10/22/00/15/spotify-1759471_1280.jpg"
	var res *http.Response

	imageFile, _ := os.Create(pv.CachedImagePath())

	if len(pv.playlistsMap) > 0 && len(pv.playlistsMap[pv.PlaylistListModel.choice].Images) > 0 {
		res, _ = http.Get(pv.GetSelectedPlaylist().Images[0].Url)
	} else {
		res, _ = http.Get(DEFAULT_IMAGE_URL)
	}

	io.Copy(imageFile, res.Body)

	if terminal.IsSizeSmall() {
		return pv.viewSmall(pv.CachedImagePath(), terminal)
	}

	// ascii, err := AsciiRender(pv.CachedImagePath(), AsciiFlagsSmall())
	ascii := ""
	// if err != nil {
	// 	ascii = "Ascii image unavailable"
	// }

	if len((*pv.UserPlaylists)) < 1 {
		return fmt.Sprintf("\n\n%s\n\n%s\n\n%s", MainControlsRender(PLAYLIST_VIEW),
			padLines(ascii, TAB_WIDTH),
			padLines("No playlists :(", TAB_WIDTH))
	}

	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false

	t.AppendRow(table.Row{
		"\n\n" + padLines(ascii, TAB_WIDTH),
		"\n\n\n" + pv.PlaylistListModel.View(),
	})

	playlist := pv.GetSelectedPlaylist()

	t2 := table.NewWriter()
	t2.Style().Options.DrawBorder = false
	t2.Style().Options.SeparateColumns = false

	style := lipgloss.NewStyle().Bold(true)

	plName := playlist.Name
	plNameLen := len(strings.Split(plName, ""))
	if plNameLen >= terminal.Width-4 {
		plName = plName[:terminal.Width-10]
		plName += "..."
	}

	t2.AppendRow(table.Row{
		padLines(style.Render("\n\n"+plName), TAB_WIDTH) + "\n",
	})

	t2.AppendRow(table.Row{
		padLines("Tracks: "+fmt.Sprint(playlist.Tracks.Total), TAB_WIDTH) + "\n",
	})

	t2.AppendRow(table.Row{
		padLines("Owner: "+fmt.Sprint(playlist.Owner.DisplayName), TAB_WIDTH) + "\n",
	})

	return "\n\n" + MainControlsRender(PLAYLIST_VIEW) + "\n" + t.Render() + "\n" + t2.Render()
}

func (pv *PlaylistView) viewSmall(imagePath string, terminal Terminal) string {
	pv.PlaylistListModel.list.SetHeight(SMALL_LIST_HEIGHT)

	var ascii string

	// ascii, err := AsciiRender(imagePath, AsciiFlagsSmall())
	ascii = ""
	// if err != nil {
	// 	ascii = "Ascii image unavailable"
	// }

	if len((*pv.UserPlaylists)) < 1 {
		return fmt.Sprintf("\n\n%s\n\n%s",
			padLines(ascii, TAB_WIDTH),
			padLines("No playlists :(", TAB_WIDTH))
	}

	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false

	t.AppendRow(table.Row{
		"\n\n\n" + padLines(ascii, TAB_WIDTH),
		"\n\n\n\n" + pv.PlaylistListModel.View(),
	})

	playlist := pv.GetSelectedPlaylist()

	t2 := table.NewWriter()
	t2.Style().Options.DrawBorder = false
	t2.Style().Options.SeparateColumns = false

	style := lipgloss.NewStyle().Bold(true)

	plName := playlist.Name
	plNameLen := len(strings.Split(plName, ""))
	if plNameLen >= terminal.Width-4 {
		plName = plName[:terminal.Width-10]
		plName += "..."
	}

	t2.AppendRow(table.Row{
		padLines(style.Render("\n\n"+plName), TAB_WIDTH) + "\n",
	})

	t2.AppendRow(table.Row{
		padLines("Tracks: "+fmt.Sprint(playlist.Tracks.Total), TAB_WIDTH) + "\n",
	})

	t2.AppendRow(table.Row{
		padLines("Owner: "+fmt.Sprint(playlist.Owner.DisplayName), TAB_WIDTH) + "\n",
	})

	return t.Render() + "\n\n" + t2.Render()
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

func NewPlaylistListModel(items []list.Item, title string) *PlaylistListModel {
	l := list.New(items, itemDelegate{}, DEFAULT_WIDTH, LIST_HEIGHT)
	l.SetFilteringEnabled(false)
	l.Title = padLines(title, 2)
	l.Styles.Title = lipgloss.NewStyle().MarginLeft(0)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)

	lm := &PlaylistListModel{list: l}

	return lm
}

func (m PlaylistListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}

		return m, cmd
	}

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m PlaylistListModel) View() string {
	return m.list.View()
}

func (m PlaylistListModel) Init() tea.Cmd {
	return nil
}

func (pv *PlaylistView) CachedImagePath() string {
	cd, _ := os.UserCacheDir()
	path := filepath.Join(cd, config.APPNAME, "playlist_ascii.jpeg")
	return path
}
