package service

import (
	"github.com/Dorrrke/note-tracker/internal/domain/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService предоставляет бизнес-логику для управления пользователями.
//
// Отвечает за валидацию входных данных, хеширование паролей,
// а также взаимодействие с репозиторием для сохранения и получения данных.
type UserService struct {
	repo  Repository
	valid *validator.Validate
}

// NewUserService создаёт новый экземпляр сервиса пользователей.
//
// Принимает реализацию интерфейса Repository и инициализирует валидатор.
// Возвращает указатель на UserService.
//
// Параметры:
//   - repo: реализация хранилища данных
//
// Возвращает:
//   - *UserService: инициализированный сервис пользователей
func NewUserService(repo Repository) *UserService {
	valid := validator.New()
	return &UserService{repo: repo, valid: valid}
}

// LoginUser выполняет аутентификацию пользователя по учётным данным.
//
// Проверяет корректность структуры запроса, находит пользователя в хранилище
// и сравнивает хеш пароля. Возвращает UID при успешной аутентификации.
//
// Параметры:
//   - user: структура с email и паролем
//
// Возвращает:
//   - string: уникальный идентификатор пользователя
//   - error: ошибка валидации, несуществующий пользователь или неверный пароль
func (us *UserService) LoginUser(user models.UserRequest) (string, error) {
	if err := us.valid.Struct(user); err != nil {
		return "", err
	}

	dbUser, err := us.repo.LoginUser(user)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return "", err
	}
	return dbUser.UID, nil
}

// RegisterUser регистрирует нового пользователя в системе.
//
// Валидирует входные данные, хеширует пароль, генерирует уникальный UID
// и сохраняет пользователя через репозиторий.
//
// Параметры:
//   - user: полная структура пользователя (User) с паролем в открытом виде
//
// Возвращает:
//   - string: UID созданного пользователя
//   - error: ошибка валидации, хеширования или сохранения
func (us *UserService) RegisterUser(user models.User) (string, error) {
	if err := us.valid.Struct(user); err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hash)

	uuid := uuid.New().String()
	user.UID = uuid

	userID, err := us.repo.RegisterUser(user)
	if err != nil {
		return "", err
	}

	return userID, nil
}

// GetUsers возвращает список всех зарегистрированных пользователей.
//
// Делегирует запрос репозиторию.
//
// Возвращает:
//   - []models.User: массив пользователей
//   - error: ошибка при получении данных из хранилища
func (us *UserService) GetUsers() ([]models.User, error) {
	users, err := us.repo.GetUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUser возвращает данные пользователя по его уникальному идентификатору.
//
// Параметры:
//   - id: строковый UID пользователя
//
// Возвращает:
//   - models.User: найденный пользователь
//   - error: ошибка, если пользователь не найден
func (us *UserService) GetUser(id string) (models.User, error) {
	user, err := us.repo.GetUser(id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// CreateUser создаёт нового пользователя (аналог RegisterUser).
//
// Используется для совместимости с API-маршрутами. Повторяет логику RegisterUser.
//
// Параметры:
//   - user: структура пользователя с открытым паролем
//
// Возвращает:
//   - string: UID созданного пользователя
//   - error: ошибка валидации, хеширования или сохранения
func (us *UserService) CreateUser(user models.User) (string, error) {
	if err := us.valid.Struct(user); err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hash)

	uuid := uuid.New().String()
	user.UID = uuid

	userID, err := us.repo.RegisterUser(user)
	if err != nil {
		return "", err
	}

	return userID, nil
}

// DeleteUser удаляет пользователя по его уникальному идентификатору.
//
// Параметры:
//   - id: UID пользователя
//
// Возвращает:
//   - error: ошибка при удалении (например, пользователь не существует)
func (us *UserService) DeleteUser(id string) error {
	err := us.repo.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUser обновляет данные существующего пользователя.
//
// Хеширует новый пароль и сохраняет обновлённую запись.
//
// Параметры:
//   - user: обновлённая структура пользователя
//
// Возвращает:
//   - error: ошибка хеширования или обновления в хранилище
func (us *UserService) UpdateUser(user models.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	err = us.repo.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}
