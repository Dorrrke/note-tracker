package dbstorage

import (
	"context"
	"time"

	"github.com/Dorrrke/note-tracker/internal/domain/models"
	"github.com/Dorrrke/note-tracker/pkg/logger"
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

func (d *DBStorage) GetTasks() ([]models.Task, error) {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := d.db.Query(ctx, "SELECT * FROM tasks")
	if err != nil {
		log.Error().Err(err).Msg("failed to get tasks from db")
		return nil, err
	}
	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.TID, &task.Title, &task.Description, &task.Stsatus, &task.CreatedAt, &task.UpdatedAt, &task.DoneAt); err != nil {
			log.Error().Err(err).Msg("failed to parse tasks from db")
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (d *DBStorage) GetTask(id string) (models.Task, error) {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var task models.Task
	row := d.db.QueryRow(ctx, "SELECT * FROM tasks WHERE tid = $1", id)
	err := row.Scan(&task.TID, &task.Title, &task.Description, &task.Stsatus, &task.CreatedAt, &task.UpdatedAt, &task.DoneAt)
	if err != nil {
		log.Error().Err(err).Msg("failed to get task from db")
		return models.Task{}, err
	}

	return task, nil
}

func (d *DBStorage) SaveTask(task models.Task) error {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := d.db.Exec(ctx, "INSERT INTO tasks (tid, title, description, status, created_at, updated_at, done_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		task.TID, task.Title, task.Description, task.Stsatus, task.CreatedAt, task.UpdatedAt, task.DoneAt)
	if err != nil {
		log.Error().Err(err).Msg("failed to save task to db")
		return err
	}

	return nil
}

func (d *DBStorage) UpdateTask(task models.Task) error {
	panic("unimplemented")
}

func (d *DBStorage) DeleteTask(id string) error {
	panic("unimplemented")
}
