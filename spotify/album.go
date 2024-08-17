package spotify

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

func (a *Album) String() string {
	if len(a.Artists) == 0 {
		return a.Name
	}

	return a.Name + " | " + a.Artists[0].Name
}

type Track struct {
	Name       string   `json:"name"`
	Uri        string   `json:"uri"`
	Album      Album    `json:"album"`
	Artists    []Artist `json:"artists"`
	DurationMs int      `json:"duration_ms"`
	ID         string   `json:"id"`
}

func (t *Track) String() string {
	if len(t.Artists) == 0 {
		return t.Name
	}

	return t.Name + " | " + t.Artists[0].Name
}

type Artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Uri  string `json:"uri"`
}
