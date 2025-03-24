package service

import (
	"github.com/Dorrrke/note-tracker/internal/domain/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo  Repository
	valid *validator.Validate
}

func NewUserService(repo Repository) *UserService {
	valid := validator.New()
	return &UserService{repo: repo, valid: valid}
}

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
