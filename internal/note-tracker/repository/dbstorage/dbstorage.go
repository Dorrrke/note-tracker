package dbstorage

import (
	"context"

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
