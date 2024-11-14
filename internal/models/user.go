package models

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"-"`
	IsActive bool   `json:"isActive"`
}

type Photo struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	UserID      uint   `json:"user_id"`
	URL         string `json:"url"`
	Description string `json:"description"`
}
