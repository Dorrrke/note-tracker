package dbstorage

import (
	"context"
	"errors"
	"time"

	repoErros "github.com/Dorrrke/note-tracker/internal/note-tracker/domain/errors"
	"github.com/Dorrrke/note-tracker/internal/note-tracker/domain/models"
	"github.com/Dorrrke/note-tracker/pkg/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (d *DBStorage) LoginUser(user models.UserRequest) (models.User, error) {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dbUser models.User
	row := d.db.QueryRow(ctx, "SELECT * FROM users WHERE login = $1", user.Login)
	if err := row.Scan(&dbUser.UID, &dbUser.Name, &dbUser.Login, &dbUser.Password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, repoErros.ErrUserNotFound
		}
		log.Error().Err(err).Msg("failed to get user from db")
		return models.User{}, err
	}
	return dbUser, nil
}

func (d *DBStorage) RegisterUser(user models.User) error {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := d.db.Exec(ctx, "INSERT INTO users (uid, name, login, password) VALUES ($1, $2, $3, $4)", user.UID, user.Name, user.Login, user.Password)
	if err != nil {
		log.Debug().Any("err", err).Str("error", err.Error()).Msg("error form pgx.Exec")
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Debug().Any("pgErr", pgErr).Msg("error form pgx.PgError")
			if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				return repoErros.ErrUserAlreadyExists
			}
		}
		log.Error().Err(err).Msg("failed to save user to db")
		return err
	}

	return nil
}
