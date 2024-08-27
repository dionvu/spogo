package spotify

import (
	"fmt"
	"strconv"

	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/utils"
)

type Track struct {
	Name       string   `json:"name"`
	Uri        string   `json:"uri"`
	Album      Album    `json:"album"`
	Artists    []Artist `json:"artists"`
	DurationMs int      `json:"duration_ms"`
	ID         string   `json:"id"`
}

func (t *Track) PlayingInfo(c *config.Config, progressMs int) string {
	artists := ""

	for i := 0; i < len(t.Artists); i++ {
		artists += t.Artists[i].Name

		if i != len(t.Artists)-1 {
			artists += ", "
		}
	}

	info := "Track: " + t.Name + "\n"

	if len(t.Artists) == 1 {
		info += "Artist: "
	} else {
		info += "Artists: "
	}

	info += artists

	progressSeconds := strconv.Itoa(((progressMs / 1000) % 60))
	progressMinutes := strconv.Itoa((progressMs / 1000) / 60)

	durationSeconds := strconv.Itoa((t.DurationMs / 1000) % 60)
	durationMinutes := strconv.Itoa((t.DurationMs / 1000) / 60)

	// Honestly idk what im doing here, but its coolish ig.
	for _, time := range []*string{&progressSeconds, &durationSeconds} {
		if len(*time) == 1 {
			*time = "0" + *time
		}
	}

	return info + fmt.Sprintf("\n%s %vm:%vs / %vm:%vs", "Progress:",
		progressMinutes, progressSeconds, durationMinutes, durationSeconds)
}

func (t *Track) String(c *config.Config) string {
	artists := t.Artists[0].Name

	if len(t.Artists) > 1 {
		artists += "& more"
	}

	info := utils.Color(t.Name, c.Color.Track.Name) + " " + utils.Color(artists, c.Color.Track.Artist) + " "

	durationSeconds := strconv.Itoa((t.DurationMs / 1000) % 60)
	durationMinutes := strconv.Itoa((t.DurationMs / 1000) / 60)

	duration := utils.Color(fmt.Sprintf("%s %vm:%vs", "Duration:",
		durationMinutes, durationSeconds), c.Color.Track.Other)

	return info + duration
}
