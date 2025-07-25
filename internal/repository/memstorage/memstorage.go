// Package memstorage предоставляет реализацию хранилища данных в памяти.
// Используется для хранения задач и пользователей во время выполнения приложения.
package memstorage

import (
	"github.com/Dorrrke/note-tracker/internal/domain/errors"
	"github.com/Dorrrke/note-tracker/internal/domain/models"
)

// MemStorage реализует интерфейс хранилища данных, храня информацию в памяти.
// Подходит для тестирования и локальной разработки.
type MemStorage struct {
	tasks map[string]models.Task
	users map[string]models.User
}

// New создает и возвращает новый экземпляр MemStorage.
// Инициализирует внутренние карты для хранения задач и пользователей.
func New() *MemStorage {
	return &MemStorage{
		tasks: make(map[string]models.Task),
		users: make(map[string]models.User),
	}
}

// GetTasks возвращает список всех задач из хранилища.
// Если задач нет, возвращает ошибку ErrEmptyTasksList.
func (m *MemStorage) GetTasks() ([]models.Task, error) {
	var tasks []models.Task
	if len(m.tasks) == 0 {
		return nil, errors.ErrEmptyTasksList
	}
	for id, task := range m.tasks {
		task.TID = id
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetTask возвращает задачу по её идентификатору.
// Если задача не найдена, возвращает ошибку ErrTaskNotFound.
func (m *MemStorage) GetTask(id string) (models.Task, error) {
	task, ok := m.tasks[id]
	if !ok {
		return models.Task{}, errors.ErrTaskNotFound
	}
	return task, nil
}

// SaveTask сохраняет новую задачу в хранилище.
// Проверяет, что задача с таким заголовком ещё не существует.
// В случае дубликата возвращает ошибку ErrTaskAlreadyExists.
func (m *MemStorage) SaveTask(task models.Task) error {
	for _, t := range m.tasks {
		if t.Title == task.Title {
			return errors.ErrTaskAlreadyExists
		}
	}
	m.tasks[task.TID] = task
	return nil
}

// UpdateTask обновляет существующую задачу в хранилище.
// Если задача с указанным идентификатором не найдена, возвращает ошибку ErrTaskNotFound.
func (m *MemStorage) UpdateTask(task models.Task) error {
	_, ok := m.tasks[task.TID]
	if !ok {
		return errors.ErrTaskNotFound
	}
	m.tasks[task.TID] = task
	return nil
}

// DeleteTask удаляет задачу из хранилища по её идентификатору.
// Если задача не найдена, возвращает ошибку ErrTaskNotFound.
func (m *MemStorage) DeleteTask(id string) error {
	_, ok := m.tasks[id]
	if !ok {
		return errors.ErrTaskNotFound
	}
	delete(m.tasks, id)
	return nil
}

// LoginUser выполняет аутентификацию пользователя по логину и паролю.
// Возвращает данные пользователя, если аутентификация успешна.
// Если пользователь не найден, возвращает ошибку ErrUserNotFound.
func (m *MemStorage) LoginUser(user models.UserRequest) (models.User, error) {
	for _, us := range m.users {
		if us.Login == user.Login {
			return us, nil
		}
	}
	return models.User{}, errors.ErrUserNotFound
}

// RegisterUser регистрирует нового пользователя в системе.
// Проверяет, что пользователь с таким логином ещё не существует.
// В случае дубликата возвращает ошибку ErrUserAlreadyExists.
// Возвращает уникальный идентификатор зарегистрированного пользователя.
func (m *MemStorage) RegisterUser(user models.User) (string, error) {
	for _, usr := range m.users {
		if usr.Login == user.Login {
			return "", errors.ErrUserAlreadyExists
		}
	}
	m.users[user.UID] = user
	return user.UID, nil
}

// GetUsers возвращает список всех пользователей из хранилища.
// Если пользователей нет, возвращает ошибку ErrEmptyUsersList.
func (m *MemStorage) GetUsers() ([]models.User, error) {
	var users []models.User
	if len(m.users) == 0 {
		return nil, errors.ErrEmptyUsersList
	}
	for id, user := range m.users {
		user.UID = id
		users = append(users, user)
	}
	return users, nil
}

// GetUser возвращает пользователя по его идентификатору.
// Если пользователь не найден, возвращает ошибку ErrUserNotFound.
func (m *MemStorage) GetUser(id string) (models.User, error) {
	task, ok := m.users[id]
	if !ok {
		return models.User{}, errors.ErrUserNotFound
	}
	return task, nil
}

// DeleteUser удаляет пользователя из хранилища по его идентификатору.
// Если пользователь не найден, возвращает ошибку ErrUserNotFound.
func (m *MemStorage) DeleteUser(id string) error {
	_, ok := m.users[id]
	if !ok {
		return errors.ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

// UpdateUser обновляет данные существующего пользователя в хранилище.
// Если пользователь с указанным идентификатором не найден, возвращает ошибку ErrUserNotFound.
func (m *MemStorage) UpdateUser(user models.User) error {
	_, ok := m.users[user.UID]
	if !ok {
		return errors.ErrUserNotFound
	}
	m.users[user.UID] = user
	return nil
}
