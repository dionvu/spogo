package app

import (
	"github.com/dionv/spogo/internal/app/config"
	"github.com/dionv/spogo/internal/app/user"
	"github.com/dionv/spogo/internal/app/user/tokens"
)

type App struct {
	Config *config.Config
	user   *user.User
}

func (a *App) User() *user.User {
	return a.user
}

func New() (*App, error) {
	c, err := config.New()
	if err != nil {
		return nil, err
	}

	return &App{
		Config: c,
		user: &user.User{
			AccessToken:  &tokens.AccessToken{},
			RefreshToken: &tokens.RefreshToken{},
		},
	}, nil
}
