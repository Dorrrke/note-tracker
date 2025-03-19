package main

import (
	"github.com/Dorrrke/note-tracker/internal/app"
	"github.com/Dorrrke/note-tracker/internal/config"
	"github.com/Dorrrke/note-tracker/internal/repository/memstorage"
	"github.com/Dorrrke/note-tracker/internal/server"
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

	repo := memstorage.New()
	server := server.New(*cfg, repo)

	app := app.NewApp(*cfg, server, repo)

	if err := app.StartApp(); err != nil {
		panic(err)
	}
}
