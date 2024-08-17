package spotify

type Playlist struct {
	Description string    `json:"description"`
	Followers   Followers `json:"followers"`
	ID          string    `json:"id"`
	Images      []Image   `json:"images"`
	Name        string    `json:"name"`
	Public      bool      `json:"public"`
	Tracks      Tracks    `json:"tracks"`
	URI         string    `json:"uri"`
	Owner       struct {
		Followers   Followers `json:"followers"`
		ID          string    `json:"id"`
		Type        string    `json:"type"`
		URI         string    `json:"uri"`
		DisplayName string    `json:"display_name"`
	} `json:"owner"`
}

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

type Track struct {
	Name       string   `json:"name"`
	Uri        string   `json:"uri"`
	Album      Album    `json:"album"`
	Artists    []Artist `json:"artists"`
	DurationMs int      `json:"duration_ms"`
	ID         string   `json:"id"`
}

// Required for playlist, why not just an []Track?
// No idea, spotify moment.
type Tracks struct {
	Limit    int     `json:"limit"`
	Next     string  `json:"next"`
	Offset   int     `json:"offset"`
	Previous string  `json:"previous"`
	Total    int     `json:"total"`
	Items    []Track `json:"items"`
}

type Artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

type Image struct {
	Url    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type Show struct {
	Description   string  `json:"description"`
	ID            string  `json:"id"`
	Images        []Image `json:"images"`
	Name          string  `json:"name"`
	Uri           string  `json:"uri"`
	TotalEpisodes int     `json:"total_episodes"`
}

type Episode struct {
	Description string      `json:"description"`
	DurationMs  int         `json:"duration_ms"`
	ID          string      `json:"id"`
	Images      []Image     `json:"images"`
	Name        string      `json:"name"`
	ReleaseDate string      `json:"release_date"`
	ResumePoint ResumePoint `json:"resume_point"`
	Uri         string      `json:"uri"`
}

type ResumePoint struct {
	FullyPlayed      bool `json:"fully_played"`
	ResumePositionMs int  `json:"resume_position_ms"`
}

type Followers struct {
	Total int `json:"total"`
}
