package app

import (
	"github.com/Dorrrke/note-tracker/internal/note-tracker/config"
	"github.com/Dorrrke/note-tracker/internal/note-tracker/server"
	"github.com/Dorrrke/note-tracker/internal/note-tracker/service"
)

type App struct {
	cfg       config.Config
	ServerApi *server.ServerApi
	repo      service.Repository
}

func NewApp(cfg config.Config, server *server.ServerApi, repo service.Repository) *App {
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
