package main

import (
	"github.com/Dorrrke/note-tracker/internal/config"
	"github.com/Dorrrke/note-tracker/internal/repository/memstorage"
	"github.com/Dorrrke/note-tracker/internal/server"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	repo := memstorage.New()

	server := server.New(*cfg, repo)

	if err := server.Start(); err != nil {
		panic(err)
	}

}
