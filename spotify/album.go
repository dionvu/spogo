package spotify

import (
	"fmt"

	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/utils"
)

type AlbumsResponse struct {
	Href     string  `json:"href"`
	Limit    int     `json:"limit"`
	Next     string  `json:"next"`
	Offset   int     `json:"offset"`
	Previous string  `json:"previous"`
	Total    int     `json:"total"`
	Items    []Album `json:"items"`
}

type Album struct {
	Images      []Image  `json:"images"`
	Artists     []Artist `json:"artists"`
	TotalTracks int      `json:"total_tracks"`
	AlbumType   string   `json:"album_type"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	ReleaseDate string   `json:"release_date"`
	Type        string   `json:"type"`
	Uri         string   `json:"uri"`
}

func (a *Album) String(c *config.Config) string {
	artist := a.Artists[0].Name

	if len(a.Artists) > 1 {
		artist += " & more"
	}

	return utils.Color(a.Name, c.Color.Album.Name) + " " +
		utils.Color(artist, c.Color.Album.Artist) +
		utils.Color(" Release: "+a.ReleaseDate, c.Color.Album.Other) +
		utils.Color(fmt.Sprint(" Tracks: ", a.TotalTracks), c.Color.Album.Other)
}
