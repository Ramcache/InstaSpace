package models

// User представляет пользователя системы.
// @description Это структура, описывающая пользователя, который регистрируется в приложении.
// @example
//
//	{
//	  "id": 1,
//	  "email": "example@example.com",
//	  "username": "john_doe"
//	}
type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"-"`
	IsActive bool   `json:"isActive"`
}
