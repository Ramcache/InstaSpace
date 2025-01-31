package models

// User представляет собой модель пользователя
//
// @swagger:model
type User struct {
	// ID пользователя
	ID int `json:"id" example:"1"`
	// Имя пользователя
	Username string `json:"username" example:"johndoe"`
	// Email пользователя
	Email string `json:"email" example:"johndoe@example.com"`
	// Пароль пользователя (не возвращается в ответах)
	Password string `json:"password,omitempty" example:"securepassword"`
	// Флаг подтверждения email
	Verified bool `json:"verified" example:"true"`
}
