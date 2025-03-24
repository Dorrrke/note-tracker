package service

import (
	"time"

	"github.com/Dorrrke/note-tracker/internal/domain/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Repository interface {
	GetTasks() ([]models.Task, error)
	GetTask(string) (models.Task, error)
	SaveTask(models.Task) error
	UpdateTask(models.Task) error
	DeleteTask(string) error

	LoginUser(models.UserRequest) (models.User, error)
	RegisterUser(models.User) (string, error)
}

type TaskService struct {
	repo  Repository
	valid *validator.Validate
}

func NewTaskService(repo Repository) *TaskService {
	valid := validator.New()
	return &TaskService{repo: repo, valid: valid}
}

func (t *TaskService) CreateTask(task models.Task) error {
	tID := uuid.New().String()
	task.TID = tID
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now
	err := t.repo.SaveTask(task)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskService) GetTasks() ([]models.Task, error) {
	tasks, err := t.repo.GetTasks()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
