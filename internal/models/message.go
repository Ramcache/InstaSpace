package models

import "time"

// Message представляет собой сообщение в беседе
//
// @swagger:model
type Message struct {
	// ID сообщения
	ID int `json:"id" example:"1"`
	// ID беседы, к которой относится сообщение
	ConversationID int `json:"conversation_id" example:"101"`
	// ID отправителя сообщения
	SenderID int `json:"sender_id" example:"42"`
	// Содержимое сообщения
	Content string `json:"content" example:"Привет! Как дела?"`
	// Дата и время отправки сообщения
	CreatedAt time.Time `json:"created_at" example:"2024-02-01T15:45:00Z"`
}
