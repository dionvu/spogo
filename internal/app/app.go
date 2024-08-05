package app

import (
	"github.com/dionv/spogo/internal/config"
	"github.com/dionv/spogo/internal/session"
)

type App struct {
	Config  *config.Config
	session *session.Session
}

func (a *App) Session() *session.Session {
	return a.session
}

func New(c *config.Config, s *session.Session) *App {
	return &App{
		Config:  c,
		session: s,
	}
}
