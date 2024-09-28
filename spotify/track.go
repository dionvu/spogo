package spotify

type Track struct {
	Name       string   `json:"name"`
	Uri        string   `json:"uri"`
	Album      Album    `json:"album"`
	Artists    []Artist `json:"artists"`
	DurationMs int      `json:"duration_ms"`
	ID         string   `json:"id"`
}

// Simplified track struct that doesn't
// link back to the album.
type AlbumTrack struct {
	Name       string   `json:"name"`
	Uri        string   `json:"uri"`
	Artists    []Artist `json:"artists"`
	DurationMs int      `json:"duration_ms"`
	ID         string   `json:"id"`
}
