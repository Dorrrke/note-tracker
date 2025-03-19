package app

import (
	"github.com/Dorrrke/note-tracker/internal/config"
	"github.com/Dorrrke/note-tracker/internal/server"
)

type App struct {
	cfg       config.Config
	ServerApi *server.ServerApi
	repo      any
}

func NewApp(cfg config.Config, server *server.ServerApi, repo any) *App {
	return &App{
		cfg:       cfg,
		ServerApi: server,
		repo:      repo,
	}
}

func (app *App) StartApp() error {
	if err := app.ServerApi.Start(); err != nil {
		return err
	}
	return nil
}
