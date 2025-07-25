// Package errors предоставляет стандартные ошибки, используемые в приложении.
package errors

import "errors"

// ErrEmptyTasksList возвращается, когда список задач пуст.
var ErrEmptyTasksList = errors.New("empty tasks list")

// ErrTaskNotFound возвращается, когда запрашиваемая задача не найдена.
var ErrTaskNotFound = errors.New("task not found")

// ErrTaskAlreadyExists возвращается, когда попытка создать задачу с уже существующим идентификатором.
var ErrTaskAlreadyExists = errors.New("task already exists")

// ErrUserNotFound возвращается, когда запрашиваемый пользователь не найден.
var ErrUserNotFound = errors.New("user not found")

// ErrUserAlreadyExists возвращается, когда попытка создать пользователя с уже существующим логином или email.
var ErrUserAlreadyExists = errors.New("user already exists")

// ErrEmptyUsersList возвращается, когда список пользователей пуст.
var ErrEmptyUsersList = errors.New("empty users list")
