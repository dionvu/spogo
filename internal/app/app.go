package app

import (
	"github.com/dionv/spogo/internal/config"
	"github.com/dionv/spogo/internal/user"
)

type App struct {
	Config *config.Config
	user   *user.User
}

func (a *App) User() *user.User {
	return a.user
}

func New(c *config.Config, u *user.User) *App {
	return &App{
		Config: c,
		user:   u,
	}
}
