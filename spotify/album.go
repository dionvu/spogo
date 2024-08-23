package spotify

import (
	"fmt"
	"strconv"
)

type Album struct {
	AlbumType   string   `json:"album_type"`
	TotalTracks int      `json:"total_tracks"`
	ID          string   `json:"id"`
	Images      []Image  `json:"images"`
	Name        string   `json:"name"`
	ReleaseDate string   `json:"release_date"`
	Type        string   `json:"type"`
	Uri         string   `json:"uri"`
	Artists     []Artist `json:"artists"`
}

func (t *Album) String() string {
	return ""
}

type Track struct {
	Name       string   `json:"name"`
	Uri        string   `json:"uri"`
	Album      Album    `json:"album"`
	Artists    []Artist `json:"artists"`
	DurationMs int      `json:"duration_ms"`
	ID         string   `json:"id"`
}

type Artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

func (t *Track) String() string {
	// if len(s.Track.Artists) == 0 {
	// 	return s.Track.Name
	// }

	names := ""

	for i := 0; i < len(t.Artists); i++ {
		names += t.Artists[i].Name

		if i != len(t.Artists)-1 {
			names += ", "
		}
	}

	info := t.Name

	if len(t.Artists) == 1 {
		info += " Artist: "
	} else {
		info += " Artists: "
	}
	info += names

	durationSeconds := strconv.Itoa((t.DurationMs / 1000) % 60)
	durationMinutes := strconv.Itoa((t.DurationMs / 1000) / 60)

	return info + fmt.Sprintf(" %s %vm:%vs", "Duration:",
		durationMinutes, durationSeconds)
}

func (t *Track) StringPlaying(progressMs int) string {
	// if len(s.Track.Artists) == 0 {
	// 	return s.Track.Name
	// }

	names := ""

	for i := 0; i < len(t.Artists); i++ {
		names += t.Artists[i].Name

		if i != len(t.Artists)-1 {
			names += ", "
		}
	}

	info := "Track: " + t.Name + "\n"

	if len(t.Artists) == 1 {
		info += "Artist: "
	} else {
		info += "Artists: "
	}
	info += names

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
