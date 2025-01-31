package models

import "time"

type Like struct {
	ID        int       `json:"id"`
	PhotoID   int       `json:"photo_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
