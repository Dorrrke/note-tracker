package app

import (
	"github.com/Dorrrke/note-tracker/internal/config"
	"github.com/Dorrrke/note-tracker/internal/server"
	"github.com/Dorrrke/note-tracker/internal/service"
)

type App struct {
	cfg       config.Config
	ServerAPI *server.API
	repo      service.Repository
}

func NewApp(cfg config.Config, server *server.API, repo service.Repository) *App {
	return &App{
		cfg:       cfg,
		ServerAPI: server,
		repo:      repo,
	}
}

func (app *App) StartApp() error {
	if err := app.ServerAPI.Start(); err != nil {
		return err
	}
	return nil
}
