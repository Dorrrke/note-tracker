package main

import (
	"context"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/Dorrrke/note-tracker/internal/app"
	"github.com/Dorrrke/note-tracker/internal/config"
	"github.com/Dorrrke/note-tracker/internal/repository/dbstorage"
	"github.com/Dorrrke/note-tracker/internal/repository/memstorage"
	"github.com/Dorrrke/note-tracker/internal/server"
	"github.com/Dorrrke/note-tracker/internal/service"
	"github.com/Dorrrke/note-tracker/pkg/logger"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(fmt.Errorf("failed read config: %w", err))
	}

	log := logger.Get(cfg.Debug)
	log.Debug().Msg("logger was initialized")
	log.Debug().Str("host", cfg.Host).Int("port", cfg.Port).Send()

	var repo service.Repository
	repo, err = dbstorage.New(context.Background(), cfg.DBDsn)
	if err != nil {
		log.Warn().Err(err).Msg("failed to connect to db, using in-memory storage instead")
		repo = memstorage.New()
	} else {
		if err = dbstorage.Migrations(cfg.DBDsn, cfg.MigratePath); err != nil {
			log.Warn().Err(err).Msg("failed to connect to db, using in-memory storage instead")
			repo = memstorage.New()
		}
	}

	userService := service.NewUserService(repo)
	taskService := service.NewTaskService(repo)
	server := server.New(*cfg, userService, taskService)
	app := app.NewApp(*cfg, server, repo)

	if err = app.StartApp(); err != nil {
		panic(err)
	}
}
