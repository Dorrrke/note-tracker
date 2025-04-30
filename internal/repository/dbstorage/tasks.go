package dbstorage

import (
	"context"

	"github.com/Dorrrke/note-tracker/internal/domain/models"
	"github.com/Dorrrke/note-tracker/pkg/logger"
)

func (d *DBStorage) GetTasks() ([]models.Task, error) {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	rows, err := d.db.Query(ctx, "SELECT * FROM tasks")
	if err != nil {
		log.Error().Err(err).Msg("failed to get tasks from db")
		return nil, err
	}
	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err = rows.Scan(&task.TID, &task.Title, &task.Description, &task.Stsatus, &task.CreatedAt, &task.UpdatedAt, &task.DoneAt); err != nil {
			log.Error().Err(err).Msg("failed to parse tasks from db")
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (d *DBStorage) GetTask(id string) (models.Task, error) {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	var task models.Task
	row := d.db.QueryRow(ctx, "SELECT * FROM tasks WHERE tid = $1", id)
	err := row.Scan(
		&task.TID,
		&task.Title,
		&task.Description,
		&task.Stsatus,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.DoneAt,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to get task from db")
		return models.Task{}, err
	}

	return task, nil
}

func (d *DBStorage) SaveTask(task models.Task) error {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	_, err := d.db.Exec(ctx, "INSERT INTO tasks (tid, title, description, status) VALUES ($1, $2, $3, $4)",
		task.TID, task.Title, task.Description, task.Stsatus)
	if err != nil {
		log.Error().Err(err).Msg("failed to save task to db")
		return err
	}

	return nil
}

func (d *DBStorage) SaveTasks(tasks []models.Task) error {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	tx, err := d.db.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to start transaction")
		return err
	}
	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			log.Debug().Err(err).Msg("failed to rollback transaction")
		}
	}()

	_, err = tx.Prepare(ctx, "save_task", "INSERT INTO tasks (tid, title, description, status) VALUES ($1, $2, $3, $4)")
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare statement")
		return err
	}

	for _, task := range tasks {
		_, err = tx.Exec(ctx, "save_task", task.TID, task.Title, task.Description, task.Stsatus)
		if err != nil {
			log.Error().Err(err).Msg("failed to save task to db")
			return err
		}
	}

	return tx.Commit(ctx)
}

func (d *DBStorage) UpdateTask(_ models.Task) error {
	panic("unimplemented")
}

func (d *DBStorage) DeleteTask(_ string) error {
	panic("unimplemented")
}
