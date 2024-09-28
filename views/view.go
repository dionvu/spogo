package ui

import (
	"strings"

	lg "github.com/charmbracelet/lipgloss"
)

type View string

func (v View) CenterVertical(terminal Terminal) View {
	return View(centerVertical(string(v), terminal))
}

func (v View) CenterHorizontal(terminal Terminal) View {
	return View(centerHorizontal(string(v), terminal))
}

func (a View) String() string {
	return string(a)
}

func centerHorizontal(s string, t Terminal, offset ...int) string {
	center := t.Width/2 - lg.Width(s)/2 - 2

	if len(offset) > 0 {
		center += offset[0]
	}

	return padLines(s, center)
}

func centerVertical(s string, t Terminal, offset ...int) string {
	center := t.Height/2 - lg.Height(s)/2 - 2

	if len(offset) > 0 {
		center += offset[0]
	}

	return padLinesTop(s, center)
}

func padLinesTop(s string, padding int) string {
	if padding < 0 {
		return s
	}

	pad := strings.Repeat("\n", padding)

	return pad + s
}

func padLines(s string, padding int) string {
	if padding < 0 {
		return s
	}

	pad := strings.Repeat(" ", padding)
	lines := strings.Split(s, "\n")

	for i, line := range lines {
		lines[i] = pad + line
	}

	return strings.Join(lines, "\n")
}
