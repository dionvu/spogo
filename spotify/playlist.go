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
type Tracks struct {
	Limit    int     `json:"limit"`
	Next     string  `json:"next"`
	Offset   int     `json:"offset"`
	Previous string  `json:"previous"`
	Total    int     `json:"total"`
	Items    []Track `json:"items"`
}

type Followers struct {
	Total int `json:"total"`
}
