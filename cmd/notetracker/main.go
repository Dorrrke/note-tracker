package main

import (
	"github.com/Dorrrke/note-tracker/internal/app"
	"github.com/Dorrrke/note-tracker/internal/config"
	"github.com/Dorrrke/note-tracker/internal/repository/memstorage"
	"github.com/Dorrrke/note-tracker/internal/repository/pgstorage"
	"github.com/Dorrrke/note-tracker/internal/server"
	"github.com/Dorrrke/note-tracker/internal/service"
	"github.com/Dorrrke/note-tracker/pkg/logger"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.Get(cfg.Debug)
	log.Debug().Msg("logger was initialized")
	log.Debug().Str("host", cfg.Host).Int("port", cfg.Port).Send()

	var repo service.Repository
	repo, err = pgstorage.NewPostgresStorage(cfg.DBDSN)
	if err != nil {
		log.Warn().Err(err).Msg("failed to connect to db, using in-memory storage instead")
		repo = memstorage.New()
	} else {
		if err = pgstorage.RunMigrations(cfg.DBDSN); err != nil {
			log.Warn().Err(err).Msg("failed to connect to db, using in-memory storage instead")
			repo = memstorage.New()
		}
	}

	userService := service.NewUserService(repo)
	taskService := service.NewTaskService(repo)
	server := server.New(*cfg, userService, taskService)
	app := app.NewApp(*cfg, server, repo)

	if err := app.StartApp(); err != nil {
		panic(err)
	}
}
