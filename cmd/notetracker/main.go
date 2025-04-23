package main

import (
	"context"

	"github.com/Dorrrke/note-tracker/gen/auth"
	"github.com/Dorrrke/note-tracker/internal/note-tracker/app"
	"github.com/Dorrrke/note-tracker/internal/note-tracker/config"
	"github.com/Dorrrke/note-tracker/internal/note-tracker/repository/dbstorage"
	"github.com/Dorrrke/note-tracker/internal/note-tracker/repository/memstorage"
	"github.com/Dorrrke/note-tracker/internal/note-tracker/server"
	"github.com/Dorrrke/note-tracker/internal/note-tracker/service"
	"github.com/Dorrrke/note-tracker/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	repo, err = dbstorage.New(context.Background(), cfg.DbDsn)
	if err != nil {
		log.Warn().Err(err).Msg("failed to connect to db, using in-memory storage instead")
		repo = memstorage.New()
	} else {
		repo = memstorage.New()
	}

	conn, err := grpc.NewClient(cfg.AuthHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create grpc client")
	}

	client := auth.NewAuthServiceClient(conn)

	userService := service.NewUserService(repo)
	taskService := service.NewTaskService(repo)
	server := server.New(*cfg, userService, taskService, client)
	app := app.NewApp(*cfg, server, repo)

	if err := app.StartApp(); err != nil {
		panic(err)
	}
}
