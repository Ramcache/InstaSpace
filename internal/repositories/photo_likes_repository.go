package repositories

import (
	"InstaSpace/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LikeRepository struct {
	DB *pgxpool.Pool
}

func NewLikeRepository(db *pgxpool.Pool) *LikeRepository {
	return &LikeRepository{DB: db}
}

var (
	ErrInvalidPhotoID = errors.New("invalid photo ID")
	ErrInvalidUserID  = errors.New("invalid user ID")
)

func (r *LikeRepository) AddLike(ctx context.Context, photoID, userID int) error {
	// Проверяем, существует ли запись
	var exists bool
	err := r.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM photos WHERE id=$1)", photoID).Scan(&exists)
	if err != nil || !exists {
		return ErrInvalidPhotoID
	}

	err = r.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE id=$1)", userID).Scan(&exists)
	if err != nil || !exists {
		return ErrInvalidUserID
	}

	// Выполняем добавление лайка
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "INSERT INTO photo_likes (photo_id, user_id) VALUES ($1, $2)", photoID, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE photos SET likes_count = likes_count + 1 WHERE id = $1", photoID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// RemoveLike удаляет лайк из таблицы photo_likes и уменьшает счетчик
func (r *LikeRepository) RemoveLike(ctx context.Context, photoID, userID int) error {
	var exists bool
	err := r.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM photo_likes WHERE photo_id=$1 AND user_id=$2)", photoID, userID).Scan(&exists)
	if err != nil || !exists {
		return errors.New("like not found")
	}

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "DELETE FROM photo_likes WHERE photo_id = $1 AND user_id = $2", photoID, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE photos SET likes_count = likes_count - 1 WHERE id = $1", photoID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GetLikes возвращает количество лайков и список пользователей
func (r *LikeRepository) GetLikes(ctx context.Context, photoID int) ([]models.User, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT u.id, u.username FROM photo_likes pl
		JOIN users u ON pl.user_id = u.id
		WHERE pl.photo_id = $1
	`, photoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetLikeCount возвращает количество лайков для фотографии
func (r *LikeRepository) GetLikeCount(ctx context.Context, photoID int) (int, error) {
	var count int
	err := r.DB.QueryRow(ctx, "SELECT likes_count FROM photos WHERE id = $1", photoID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
