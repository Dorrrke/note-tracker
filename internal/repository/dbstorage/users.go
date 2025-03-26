package dbstorage

import "github.com/Dorrrke/note-tracker/internal/domain/models"

func (d *DBStorage) LoginUser(user models.UserRequest) (models.User, error) {
	panic("unimplemented")
}

func (d *DBStorage) RegisterUser(user models.User) (string, error) {
	panic("unimplemented")
}
