package ui

import (
	"strings"

	lg "github.com/charmbracelet/lipgloss"
)

// Content is a string with methods that allow ease
// of manipulating the content dynamically with
// respect to given terminal's dimensions.
type Content string

// Joins subcontents into a since content, seperating
// them by a single sep.
func Join(contents []Content, sep string) Content {
	s := ""

	for i, c := range contents {
		s += string(c)
		if i == len(contents)-1 {
			break
		}

		s += sep
	}

	return Content(s)
}

// Splits into subcontents and seperateing
// them by a single sep.
func (v Content) Split(sep byte) []Content {
	// Split the string content using the specified separator
	parts := strings.Split(string(v), string(sep))

	// Convert the result back into a slice of Content
	var result []Content
	for _, part := range parts {
		result = append(result, Content(part))
	}

	return result
}

// Centers the content along the Y-axis given the terminal size.
func (v Content) CenterVertical(t Terminal, offset ...int) Content {
	s := string(v)

	center := t.Height/2 - lg.Height(s)/2

	if len(offset) > 0 {
		center += offset[0]
	}

	return v.PadLinesTop(center)
}

// Centers the content along the X-axis given the terminal size.
func (v Content) CenterHorizontal(t Terminal, offset ...int) Content {
	lines := v.Split('\n')

	for i, line := range lines {
		// center := int(math.Ceil(float64(t.Width)/2 - float64(lg.Width(string(line)))/2))
		center := t.Width/2 - lg.Width(string(line))/2

		if len(offset) > 0 {
			center += offset[0]
		}

		lines[i] = line.PadLinesLeft(center)
	}

	return Join(lines, "\n")
}

// The content as a string.
func (a Content) String() string {
	return string(a)
}

// Appends a '\n' to the front of the string for the
// and repeats count amount of times.
func (v Content) PadLinesTop(count int) Content {
	s := string(v)

	if count < 0 {
		return v
	}

	pad := strings.Repeat("\n", count)

	return Content(pad + s)
}

// Appends a ' ' to the front of each line of
// the string for the and repeats count amount
// of times.
func (v Content) PadLinesLeft(count int) Content {
	s := string(v)

	if count < 0 {
		return v
	}

	pad := strings.Repeat(" ", count)
	lines := strings.Split(s, "\n")

	for i, line := range lines {
		lines[i] = pad + line
	}

	return Content(strings.Join(lines, "\n"))
}
