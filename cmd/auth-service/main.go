package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/Dorrrke/note-tracker/internal/auth-service/app"
	"github.com/Dorrrke/note-tracker/internal/auth-service/config"
	"github.com/Dorrrke/note-tracker/internal/auth-service/repository"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	if err := Migrations(cfg.DBDsn, cfg.MigratePath); err != nil {
		log.Fatal(err)
	}

	repository, err := repository.New(cfg.DBDsn)
	if err != nil {
		log.Fatal(err)
	}
	app := app.NewApp(cfg, repository)
	if err := app.StartApp(); err != nil {
		panic(err)
	}
}

func Migrations(dbDsn string, migratePath string) error {
	path := fmt.Sprintf("file://%s", migratePath)
	m, err := migrate.New(path, dbDsn)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	}

	return nil
}
