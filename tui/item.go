package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	LIST_HEIGHT         = 10
	DEFAULT_WIDTH       = 20
	DEFAULT_LIST_HEIGHT = 20 // 7 Selections per page
	SMALL_LIST_HEIGHT   = 7  // 5 Selections per page
)

type Item string

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i)

	fn := PlaylistViewStyle.Item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return PlaylistViewStyle.ItemSelected.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (i Item) FilterValue() string {
	return string(i)
}
