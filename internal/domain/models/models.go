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
