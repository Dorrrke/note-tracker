package service

import (
	"time"

	"github.com/Dorrrke/note-tracker/internal/domain/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Repository определяет контракт для работы с хранилищем данных.
//
// Включает методы для управления задачами и пользователями.
// Должен быть реализован конкретным адаптером.
type Repository interface {
	GetTasks() ([]models.Task, error)
	GetTask(string) (models.Task, error)
	SaveTask(models.Task) error
	UpdateTask(models.Task) error
	DeleteTask(string) error

	LoginUser(models.UserRequest) (models.User, error)
	RegisterUser(models.User) (string, error)
	GetUsers() ([]models.User, error)
	GetUser(string) (models.User, error)
	DeleteUser(string) error
	UpdateUser(models.User) error
}

// TaskService предоставляет бизнес-логику для управления задачами.
//
// Инкапсулирует валидацию, генерацию ID и работу с репозиторием.
type TaskService struct {
	repo  Repository
	valid *validator.Validate
}

// NewTaskService создаёт новый экземпляр сервиса задач.
//
// Принимает реализацию Repository и инициализирует валидатор.
// Возвращает указатель на TaskService.
//
// Параметры:
//   - repo: реализация интерфейса Repository
//
// Возвращает:
//   - *TaskService: инициализированный сервис задач
func NewTaskService(repo Repository) *TaskService {
	valid := validator.New()
	return &TaskService{repo: repo, valid: valid}
}

// CreateTask создаёт новую задачу с уникальным ID и временными метками.
//
// Генерирует UUID, устанавливает CreatedAt и UpdatedAt, затем сохраняет задачу.
//
// Параметры:
//   - task: исходная структура задачи без ID
//
// Возвращает:
//   - error: ошибка при сохранении в репозиторий
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

// GetTasks возвращает все задачи из хранилища.
//
// Делегирует вызов репозиторию.
//
// Возвращает:
//   - []models.Task: список задач
//   - error: ошибка при получении данных
func (t *TaskService) GetTasks() ([]models.Task, error) {
	tasks, err := t.repo.GetTasks()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTask возвращает задачу по её уникальному идентификатору.
//
// Параметры:
//   - id: строковый идентификатор задачи
//
// Возвращает:
//   - models.Task: найденная задача
//   - error: ошибка, если задача не найдена
func (t *TaskService) GetTask(id string) (models.Task, error) {
	task, err := t.repo.GetTask(id)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

// DeleteTask удаляет задачу по её ID.
//
// Параметры:
//   - id: идентификатор задачи
//
// Возвращает:
//   - error: ошибка при удалении
func (t *TaskService) DeleteTask(id string) error {
	err := t.repo.DeleteTask(id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateTask обновляет существующую задачу.
//
// Предполагается, что задача уже содержит корректный TID.
//
// Параметры:
//   - task: обновлённая структура задачи
//
// Возвращает:
//   - error: ошибка при обновлении в хранилище
func (t *TaskService) UpdateTask(task models.Task) error {
	err := t.repo.UpdateTask(task)
	if err != nil {
		return err
	}
	return nil
}
