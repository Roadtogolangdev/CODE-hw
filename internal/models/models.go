package models

import "time"

type User struct {
	ID       int
	Username string
	Password string
}

type Note struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UserID    int       `json:"user_id"`
}
