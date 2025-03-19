package errors

import "errors"

var ErrEmptyTasksList = errors.New("empty tasks list")
var ErrTaskNotFound = errors.New("task not found")
var ErrTaskAlreadyExists = errors.New("task already exists")
