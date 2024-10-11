package components

import "github.com/jedib0t/go-pretty/v6/table"

// Returns a new table.Writer that has visble
// seperators and borders completely disabled.
func NewDefaultTable() table.Writer {
	l := table.NewWriter()
	l.Style().Options.DrawBorder = false
	l.Style().Options.SeparateColumns = false
	l.Style().Options.SeparateRows = false

	return l
}
