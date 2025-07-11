package app

import (
	"github.com/Dorrrke/note-tracker/internal/config"
	"github.com/Dorrrke/note-tracker/internal/server"
	"github.com/Dorrrke/note-tracker/internal/service"
)

type App struct {
	cfg  config.Config
	API  *server.API
	repo service.Repository
}

func NewApp(cfg config.Config, server *server.API, repo service.Repository) *App {
	return &App{
		cfg:  cfg,
		API:  server,
		repo: repo,
	}
}

func (app *App) StartApp() error {
	if err := app.API.Start(); err != nil {
		return err
	}
	return nil
}
