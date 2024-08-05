package responses

// Response from the "https://api.spotify.com/v1/me" endpoint.
type UserResponse struct {
	Country     string `json:"country"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`

	ExplicitContent struct {
		FilterEnabled bool `json:"filter_enabled"`
		FilterLocked  bool `json:"filter_locked"`
	} `json:"explicit_content"`

	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`

	Followers struct {
		Href  string `json:"href"`
		Total int    `json:"total"`
	} `json:"followers"`

	Href string `json:"href"`
	ID   string `json:"id"`

	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`

	Product string `json:"product"`
	Type    string `json:"type"`
	URI     string `json:"uri"`
}
