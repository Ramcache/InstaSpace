package models

import "time"

// Like представляет собой лайк на фото
//
// @swagger:model
type Like struct {
	// ID лайка
	ID int `json:"id" example:"1"`
	// ID фотографии, на которую поставлен лайк
	PhotoID int `json:"photo_id" example:"101"`
	// ID пользователя, который поставил лайк
	UserID int `json:"user_id" example:"42"`
	// Дата и время добавления лайка
	CreatedAt time.Time `json:"created_at" example:"2024-02-01T16:30:00Z"`
}
