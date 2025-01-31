package models

// Photo представляет собой фотографию, загруженную пользователем
//
// @swagger:model
type Photo struct {
	// ID фотографии
	ID int `json:"id" example:"1"`
	// ID пользователя, загрузившего фото
	UserID int `json:"user_id" example:"42"`
	// URL изображения
	URL string `json:"url" example:"https://example.com/uploads/photo1.jpg"`
	// Описание фотографии
	Description string `json:"description" example:"Закат на пляже"`
	// Дата загрузки фото (в формате ISO 8601)
	CreatedAt string `json:"created_at" example:"2024-02-01T16:00:00Z"`
}
