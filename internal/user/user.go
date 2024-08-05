package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dionv/spogo/internal/session"
	"github.com/fatih/color"
)

type User struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	ID          string `json:"id"`

	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`

	Followers struct {
		// Always null atm (due to spotify api)
		Href string `json:"href"`

		Total int `json:"total"`
	} `json:"followers"`

	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
}

func New(s *session.Session) (*User, error) {
	ep := "https://api.spotify.com/v1/me"
	req, _ := http.NewRequest(http.MethodGet, ep, nil)

	req.Header.Add("Authorization", "Bearer "+s.AccessToken.String())

	res, _ := http.DefaultClient.Do(req)

	u := &User{}

	_ = json.NewDecoder(res.Body).Decode(u)

	return u, nil
}

func (u *User) Print() {
	a := strconv.Itoa(u.Followers.Total)

	fmt.Printf("%v | %v | %v\n", color.RedString(u.DisplayName), color.BlueString(a+" Followers"), color.YellowString(u.Email))
}
