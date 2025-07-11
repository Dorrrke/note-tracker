package pgstorage

import (
	"fmt"

	"github.com/Dorrrke/note-tracker/internal/domain/models"
	"github.com/Dorrrke/note-tracker/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresStorage struct {
	db *gorm.DB
}

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

func (p *PostgresStorage) GetTasks() ([]models.Task, error) {
	var tasks []models.Task
	result := p.db.Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

func (p *PostgresStorage) GetTask(id string) (models.Task, error) {
	var task models.Task
	result := p.db.First(&task, "tid = ?", id)
	if result.Error != nil {
		return models.Task{}, result.Error
	}
	return task, nil
}

func (p *PostgresStorage) SaveTask(task models.Task) error {
	result := p.db.Create(&task)
	return result.Error
}

func (p *PostgresStorage) UpdateTask(task models.Task) error {
	result := p.db.Save(&task)
	return result.Error
}

func (p *PostgresStorage) DeleteTask(id string) error {
	result := p.db.Delete(&models.Task{}, "tid = ?", id)
	return result.Error
}

func (p *PostgresStorage) RegisterUser(user models.User) (string, error) {
	result := p.db.Create(&user)
	if result.Error != nil {
		return "", result.Error
	}
	return user.UID, nil
}

func (p *PostgresStorage) LoginUser(userRequest models.UserRequest) (models.User, error) {
	var user models.User
	result := p.db.Where("login = ?", userRequest.Login).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (p *PostgresStorage) GetUsers() ([]models.User, error) {
	var users []models.User
	result := p.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (p *PostgresStorage) GetUser(id string) (models.User, error) {
	var user models.User
	result := p.db.First(&user, "uid = ?", id)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (p *PostgresStorage) DeleteUser(id string) error {
	result := p.db.Delete(&models.User{}, "uid = ?", id)
	return result.Error
}

func (p *PostgresStorage) UpdateUser(user models.User) error {
	result := p.db.Save(&user)
	return result.Error
}
