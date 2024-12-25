package models

type Photo struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	URL         string `json:"url"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}
