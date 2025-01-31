package models

import "time"

// Comment представляет собой комментарий пользователя к фото
//
// @swagger:model
type Comment struct {
	// ID комментария
	ID int `json:"id" example:"1"`
	// ID пользователя, оставившего комментарий
	UserID int `json:"user_id" example:"42"`
	// ID фото, к которому относится комментарий
	PhotoID int `json:"photo_id" example:"101"`
	// Текст комментария
	Content string `json:"content" example:"Отличное фото!"`
	// Дата создания комментария
	CreatedAt time.Time `json:"created_at" example:"2024-01-31T12:45:00Z"`
	// Дата последнего обновления комментария
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-31T12:50:00Z"`
	// Имя пользователя, оставившего комментарий
	Username string `json:"username" example:"johndoe"`
}
