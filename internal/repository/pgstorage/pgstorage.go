package pgstorage

import (
	"fmt"

	"github.com/Dorrrke/note-tracker/internal/domain/models"
	"github.com/Dorrrke/note-tracker/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresStorage реализует интерфейс Repository для работы с PostgreSQL.
//
// Использует GORM в качестве ORM для взаимодействия с базой данных.
type PostgresStorage struct {
	db *gorm.DB
}

// NewPostgresStorage создаёт и инициализирует новое подключение к PostgreSQL.
//
// Выполняет автоматическую миграцию моделей Task и User при первом запуске.
//
// Параметры:
//   - dsn: строка подключения к базе данных (Data Source Name)
//
// Возвращает:
//   - *PostgresStorage: инициализированный адаптер хранилища
//   - error: ошибка подключения или миграции
func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("fail when trying to connect to database: %w", err)
	}

	err = db.AutoMigrate(&models.Task{}, &models.User{})
	if err != nil {
		return nil, fmt.Errorf("fail when trying to migrate database: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

// RunMigrations выполняет миграции базы данных только при отсутствии таблиц.
//
// Проверяет наличие таблиц `tasks` и `users`. Если хотя бы одной нет — запускает AutoMigrate.
// Используется для безопасного обновления схемы при запуске приложения.
//
// Параметры:
//   - dsn: строка подключения к PostgreSQL
//
// Возвращает:
//   - error: ошибка подключения или выполнения миграции
func RunMigrations(dsn string) error {
	log := logger.Get()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if !db.Migrator().HasTable(&models.Task{}) || !db.Migrator().HasTable(&models.User{}) {
		log.Info().Msg("Tables do not exist. Performing migration...")
		if err = db.AutoMigrate(&models.User{}, &models.Task{}); err != nil {
			return fmt.Errorf("failed to migrate database: %w", err)
		}
	} else {
		log.Info().Msg("Tables already exist. Skipping migration.")
	}
	return nil
}

// GetTasks возвращает все задачи из базы данных.
//
// Возвращает:
//   - []models.Task: список всех задач
//   - error: ошибка выполнения запроса
func (p *PostgresStorage) GetTasks() ([]models.Task, error) {
	var tasks []models.Task
	result := p.db.Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

// GetTask возвращает задачу по её уникальному идентификатору (TID).
//
// Ищет запись в таблице tasks по полю tid.
//
// Параметры:
//   - id: строковый идентификатор задачи
//
// Возвращает:
//   - models.Task: найденная задача
//   - error: ошибка или запись не найдена
func (p *PostgresStorage) GetTask(id string) (models.Task, error) {
	var task models.Task
	result := p.db.First(&task, "tid = ?", id)
	if result.Error != nil {
		return models.Task{}, result.Error
	}
	return task, nil
}

// SaveTask создаёт новую задачу в базе данных.
//
// Параметры:
//   - task: структура задачи для сохранения
//
// Возвращает:
//   - error: ошибка при вставке записи
func (p *PostgresStorage) SaveTask(task models.Task) error {
	result := p.db.Create(&task)
	return result.Error
}

// UpdateTask обновляет существующую задачу.
//
// Параметры:
//   - task: обновлённая структура задачи (должна содержать TID)
//
// Возвращает:
//   - error: ошибка при обновлении
func (p *PostgresStorage) UpdateTask(task models.Task) error {
	result := p.db.Save(&task)
	return result.Error
}

// DeleteTask удаляет задачу по её TID.
//
// Параметры:
//   - id: идентификатор задачи
//
// Возвращает:
//   - error: ошибка выполнения удаления
func (p *PostgresStorage) DeleteTask(id string) error {
	result := p.db.Delete(&models.Task{}, "tid = ?", id)
	return result.Error
}

// RegisterUser создаёт нового пользователя в базе данных.
//
// Параметры:
//   - user: структура пользователя
//
// Возвращает:
//   - string: UID созданного пользователя
//   - error: ошибка при вставке
func (p *PostgresStorage) RegisterUser(user models.User) (string, error) {
	result := p.db.Create(&user)
	if result.Error != nil {
		return "", result.Error
	}
	return user.UID, nil
}

// LoginUser находит пользователя по логину.
//
// Параметры:
//   - userRequest: структура с полем Login
//
// Возвращает:
//   - models.User: найденный пользователь
//   - error: пользователь не найден или ошибка запроса
func (p *PostgresStorage) LoginUser(userRequest models.UserRequest) (models.User, error) {
	var user models.User
	result := p.db.Where("login = ?", userRequest.Login).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

// GetUsers возвращает список всех пользователей.
//
// Возвращает:
//   - []models.User: массив пользователей
//   - error: ошибка выполнения запроса
func (p *PostgresStorage) GetUsers() ([]models.User, error) {
	var users []models.User
	result := p.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// GetUser возвращает пользователя по его уникальному идентификатору.
//
// Ищет запись в таблице users по полю uid.
//
// Параметры:
//   - id: строковый UID
//
// Возвращает:
//   - models.User: найденный пользователь
//   - error: запись не найдена или ошибка запроса
func (p *PostgresStorage) GetUser(id string) (models.User, error) {
	var user models.User
	result := p.db.First(&user, "uid = ?", id)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

// DeleteUser удаляет пользователя по его UID.
//
// Параметры:
//   - id: идентификатор пользователя
//
// Возвращает:
//   - error: ошибка выполнения удаления
func (p *PostgresStorage) DeleteUser(id string) error {
	result := p.db.Delete(&models.User{}, "uid = ?", id)
	return result.Error
}

// UpdateUser обновляет данные пользователя.
//
// Параметры:
//   - user: обновлённая структура пользователя
//
// Возвращает:
//   - error: ошибка при обновлении
func (p *PostgresStorage) UpdateUser(user models.User) error {
	result := p.db.Save(&user)
	return result.Error
}
