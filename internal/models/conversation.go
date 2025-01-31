package models

import "time"

// Conversation представляет собой диалог между двумя пользователями
//
// @swagger:model
type Conversation struct {
	// ID беседы
	ID int `json:"id" example:"1"`
	// ID первого пользователя
	User1ID int `json:"user1_id" example:"42"`
	// ID второго пользователя
	User2ID int `json:"user2_id" example:"58"`
	// Дата создания беседы
	CreatedAt time.Time `json:"created_at" example:"2024-02-01T14:30:00Z"`
}
