package spotify

import (
	"strconv"

	"github.com/dionvu/spogo/config"
)

type TracksResponse struct {
	Href     string  `json:"href"`
	Limit    int     `json:"limit"`
	Next     string  `json:"next"`
	Offset   int     `json:"offset"`
	Previous string  `json:"previous"`
	Total    int     `json:"total"`
	Items    []Track `json:"items"`
}

type Track struct {
	Name       string   `json:"name"`
	Uri        string   `json:"uri"`
	Album      Album    `json:"album"`
	Artists    []Artist `json:"artists"`
	DurationMs int      `json:"duration_ms"`
	ID         string   `json:"id"`
}

func (t *Track) InfoString(c *config.Config, progressMs int) (
	track string, artist string,
	progressMinutes string, progressSeconds string,
	durationMinutes string, durationSeconds string,
) {
	for i := 0; i < len(t.Artists); i++ {
		artist += t.Artists[i].Name

		if i != len(t.Artists)-1 {
			artist += ", "
		}
	}

	track = t.Name

	progressSeconds = strconv.Itoa(((progressMs / 1000) % 60))
	progressMinutes = strconv.Itoa((progressMs / 1000) / 60)

	durationSeconds = strconv.Itoa((t.DurationMs / 1000) % 60)
	durationMinutes = strconv.Itoa((t.DurationMs / 1000) / 60)

	// Honestly idk what im doing here, but its coolish ig.
	for _, time := range []*string{&progressSeconds, &durationSeconds} {
		if len(*time) == 1 {
			*time = "0" + *time
		}
	}

	return track, artist, progressMinutes, progressSeconds, durationMinutes, durationSeconds
}
