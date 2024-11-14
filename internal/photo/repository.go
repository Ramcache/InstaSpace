package photo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

// Photo представляет структуру фото
type Photo struct {
	ID          int       `json:"id"`          // Идентификатор фото
	UserID      int       `json:"user_id"`     // Идентификатор пользователя, загрузившего фото
	ImageURL    string    `json:"image_url"`   // URL изображения
	Description string    `json:"description"` // Описание изображения
	CreatedAt   time.Time `json:"created_at"`  // Дата и время создания фото
}

// Error представляет структуру ошибки для API
type Error struct {
	Message string `json:"message"` // Описание ошибки
}

// Repository предоставляет методы для работы с фото
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository создает новый репозиторий для работы с фото
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SavePhoto(ctx context.Context, photo *Photo) error {
	query := `INSERT INTO photos (user_id, image_url, description) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, photo.UserID, photo.ImageURL, photo.Description).Scan(&photo.ID, &photo.CreatedAt)
	return err
}

func (r *Repository) GetPhotoByID(ctx context.Context, id int) (*Photo, error) {
	photo := &Photo{}
	query := `SELECT id, user_id, image_url, description, created_at FROM photos WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&photo.ID,
		&photo.UserID,
		&photo.ImageURL,
		&photo.Description,
		&photo.CreatedAt,
	)
	if err != nil {
		return nil, errors.New("photo not found")
	}
	return photo, nil
}

func (r *Repository) GetUserIDByEmail(ctx context.Context, email string) (int, error) {
	var userID int
	query := `SELECT id FROM users WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(&userID)
	if err != nil {
		return 0, errors.New("user not found")
	}
	return userID, nil
}
