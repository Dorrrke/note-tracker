package app

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/Dorrrke/note-tracker/internal/config"
	"github.com/Dorrrke/note-tracker/internal/server"
	"github.com/Dorrrke/note-tracker/internal/service"
	"github.com/Dorrrke/note-tracker/pkg/logger"
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

func (app *App) StartApp(ctx context.Context) error {
	log := logger.Get()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Debug().Msg("ctx: Done; stop app")
		if err := app.StopApp(ctx); err != nil {
			log.Error().Err(err).Msg("failed to stop app")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.ServerApi.Start(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			log.Error().Err(err).Msg("failed to start server")
		}
	}()

	wg.Wait()
	return nil
}

func (app *App) StopApp(ctx context.Context) error {
	log := logger.Get()
	err := app.ServerApi.Stop(ctx)
	if err != nil {
		return err
	}
	log.Debug().Msg("server was stopped")
	err = app.repo.Stop(ctx)
	if err != nil {
		return err
	}
	log.Debug().Msg("repo was stopped")
	return nil
}
