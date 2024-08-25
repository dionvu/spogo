package spotify

import "github.com/fatih/color"

type Artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

func (a *Artist) String() string {
	return color.HiYellowString(a.Name)
}
