package repository

import (
	"errors"
	"fmt"

	"github.com/Dorrrke/note-tracker/internal/auth-service/domain/models"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	gormDB *gorm.DB
}

func New(connStr string) (*Repository, error) {
	conn, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  connStr,
			PreferSimpleProtocol: true,
		}))
	if err != nil {
		return nil, err
	}

	return &Repository{gormDB: conn}, nil
}

func (r *Repository) LoginUser(userCreds models.UserCreds) (models.User, error) {
	var user models.User
	if err := r.gormDB.Where("login = ?", userCreds.Login).First(&user).Error; err != nil {
		return models.User{}, fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userCreds.Password)); err != nil {
		return models.User{}, fmt.Errorf("invalid password")
	}

	return user, nil
}

func (r *Repository) RegisterUser(user models.User) (string, error) {
	uid := uuid.New().String()
	user.UID = uid
	if err := r.gormDB.Create(&user).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Debug().Any("pgErr", pgErr).Msg("error form pgx.PgError")
			if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				return "", errors.New("user already exists")
			}
		}
		return "", err
	}

	return uid, nil
}
