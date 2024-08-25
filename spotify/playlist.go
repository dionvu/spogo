package spotify

import (
	"fmt"

	"github.com/fatih/color"
)

type Playlist struct {
	Images      []Image `json:"images"`
	Description string  `json:"description"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Total       int     `json:"total"`
	Public      bool    `json:"public"`
	Tracks      Tracks  `json:"tracks"`
	URI         string  `json:"uri"`
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

func (p *Playlist) String() string {
	name := p.Name

	if len(name) > 50 {
		name = name[:50] + "..."
	}

	return color.HiYellowString(name) + " " + color.HiBlueString(p.Owner.DisplayName) + fmt.Sprint(" Tracks: ", p.Tracks.Total)
}
