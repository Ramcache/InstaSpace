package repositories

import (
	"InstaSpace/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PhotoRepository struct {
	DB *pgxpool.Pool
}

func NewPhotoRepository(db *pgxpool.Pool) *PhotoRepository {
	return &PhotoRepository{DB: db}
}

type PhotoRepositoryInterface interface {
	Create(photo *models.Photo) error
}

func (r *PhotoRepository) Create(photo *models.Photo) error {
	query := `
		INSERT INTO photos (user_id, url, description, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id`
	return r.DB.QueryRow(context.Background(), query, photo.UserID, photo.URL, photo.Description).Scan(&photo.ID)
}
