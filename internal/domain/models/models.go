// Package models содержит структуры данных (модели), используемые в приложении.
package models

import "time"

// Task представляет собой задачу в системе.
// Содержит информацию о заголовке, описании, статусе и временных метках задачи.
type Task struct {
	TID         string    `json:"tid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Stsatus     string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DoneAt      time.Time `json:"done_at"`
}

// User представляет собой пользователя системы.
// Содержит информацию о пользователе, включая учетные данные для аутентификации.
type User struct {
	UID      string `json:"uid"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}

// UserRequest используется для передачи данных при аутентификации пользователя.
// Содержит логин и пароль, необходимые для входа в систему.
type UserRequest struct {
	Login    string `json:"login"    validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}
