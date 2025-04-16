package models

import "time"

type Task struct {
	TID         string    `json:"tid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Stsatus     string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DoneAt      time.Time `json:"done_at"`
}

type User struct {
	UID      string `json:"uid"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserRequest struct {
	Login    string `json:"login" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}
