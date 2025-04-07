package dbstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/Dorrrke/note-tracker/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
)

// type Repository interface {
// 	GetTasks() ([]models.Task, error)
// 	GetTask(string) (models.Task, error)
// 	SaveTask(models.Task) error
// 	UpdateTask(models.Task) error
// 	DeleteTask(string) error

// 	LoginUser(models.UserRequest) (models.User, error)
// 	RegisterUser(models.User) (string, error)
// }

type DBStorage struct {
	db *pgx.Conn
}

func New(ctx context.Context, addr string) (*DBStorage, error) {
	conn, err := pgx.Connect(ctx, addr)
	if err != nil {
		return nil, err
	}

	return &DBStorage{db: conn}, nil
}

func Migrations(dbDsn string, migratePath string) error {
	log := logger.Get()

	path := fmt.Sprintf("file://%s", migratePath)
	m, err := migrate.New(path, dbDsn)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("failed to create migrate instance")
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Debug().Msg("no migrations were applied")
			return nil
		}
		log.Error().Err(err).Msg("failed to migrate db")
		return err
	}

	log.Debug().Msg("all migrations were applied")
	return nil
}
