package components

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	LIST_HEIGHT_SMALL   = 7
	LIST_HEIGHT_NORMAL  = 10
	DEFAULT_WIDTH       = 20
	DEFAULT_LIST_HEIGHT = 20 // 7 Selections per page
	SMALL_LIST_HEIGHT   = 7  // 5 Selections per page
)

type UniqueItem struct {
	Name string
	Id   string
}

type UniqueItemDelegate struct{}

func (d UniqueItemDelegate) Height() int                             { return 1 }
func (d UniqueItemDelegate) Spacing() int                            { return 0 }
func (d UniqueItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d UniqueItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(UniqueItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.Name)

	fn := PlaylistViewStyle.ListItem.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return PlaylistViewStyle.ItemSelected.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (i UniqueItem) FilterValue() string {
	return ""
}

type ListItem string

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                             { return 1 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(ListItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i)

	fn := PlaylistViewStyle.ListItem.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return PlaylistViewStyle.ItemSelected.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (i ListItem) FilterValue() string {
	return ""
}

var CommonStyle = struct {
	MainControls struct {
		Selected lg.Style
		Normal   lg.Style
	}
}{
	MainControls: struct {
		Selected lg.Style
		Normal   lg.Style
	}{
		Normal:   lg.NewStyle().Faint(true),
		Selected: lg.NewStyle(),
	},
}

var PlayerViewStyle = struct {
	StatusBar struct {
		NowPlaying lg.Style
		Paused     lg.Style
		NoPlayer   lg.Style
	}
}{
	StatusBar: struct {
		NowPlaying lg.Style
		Paused     lg.Style
		NoPlayer   lg.Style
	}{
		NowPlaying: lg.NewStyle().
			Bold(true).
			Foreground(lg.Color("#282828")).
			Background(lg.Color("#98971a")).
			PaddingLeft(1).
			PaddingRight(1),

		Paused: lg.NewStyle().
			Bold(true).
			Foreground(lg.Color("#282828")).
			Background(lg.Color("#d79921")).
			PaddingLeft(1).
			PaddingRight(1),

		NoPlayer: lg.NewStyle().
			Bold(true).
			Foreground(lg.Color("#282828")).
			Background(lg.Color("#cc241d")).
			PaddingLeft(1).
			PaddingRight(1),
	},
}

var PlaylistViewStyle = struct {
	Title        lg.Style
	ItemSelected lg.Style
	ListItem     lg.Style
}{
	Title:        lg.NewStyle().Bold(true).Background(lg.Color("#458588")).Foreground(lg.Color("#282828")).PaddingLeft(1).PaddingRight(1),
	ItemSelected: lg.NewStyle().PaddingLeft(2),

	ListItem: lg.NewStyle().PaddingLeft(4).Faint(true),
}

var DeviceViewStyle = struct {
	Title        lg.Style
	ItemSelected lg.Style
	ListItem     lg.Style
}{
	Title:        lg.NewStyle().Bold(true).Background(lg.Color("#a89984")).Foreground(lg.Color("#282828")).PaddingLeft(1).PaddingRight(1),
	ItemSelected: lg.NewStyle().PaddingLeft(2),

	ListItem: lg.NewStyle().PaddingLeft(4).Faint(true),
}

func NewDefaultList(items []list.Item, title string) list.Model {
	l := list.New(items, ItemDelegate{}, DEFAULT_WIDTH, LIST_HEIGHT_NORMAL)
	l.Styles.Title = lg.NewStyle().MarginLeft(0)
	l.Title = title
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)

	return l
}

func NewCustomList(items []list.Item, title string, height int, width int) list.Model {
	l := list.New(items, ItemDelegate{}, height, width)
	l.Styles.Title = lg.NewStyle().MarginLeft(0)
	l.Title = title
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)

	return l
}

func NewDefaultUniqueItemList(items []list.Item, title string) list.Model {
	l := list.New(items, UniqueItemDelegate{}, DEFAULT_WIDTH, LIST_HEIGHT_NORMAL)
	l.Styles.Title = lg.NewStyle().MarginLeft(0)
	l.Title = title
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)

	return l
}

func NewCustomUniqueItemList(items []list.Item, title string, height int, width int) list.Model {
	l := list.New(items, UniqueItemDelegate{}, height, width)
	l.Styles.Title = lg.NewStyle().MarginLeft(0)
	l.Title = title
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)

	return l
}

// Returns a new table.Writer that has visble
// seperators and borders completely disabled.
func NewDefaultTable() table.Writer {
	l := table.NewWriter()
	l.Style().Options.DrawBorder = false
	l.Style().Options.SeparateColumns = false
	l.Style().Options.SeparateRows = false

	return l
}
