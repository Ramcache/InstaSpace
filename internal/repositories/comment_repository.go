package repositories

import (
	"InstaSpace/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepository struct {
	DB *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{DB: db}
}

type CommentRepositoryInterface interface {
	CreateComment(ctx context.Context, comment *models.Comment) (int, error)
	GetCommentsByPhotoID(ctx context.Context, photoID int) ([]models.Comment, error)
	UpdateComment(ctx context.Context, comment *models.Comment) error
	DeleteComment(ctx context.Context, commentID, userID int) error
}

func (r *CommentRepository) CreateComment(ctx context.Context, comment *models.Comment) (int, error) {
	query := `
		INSERT INTO comments (user_id, photo_id, content, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id`
	err := r.DB.QueryRow(ctx, query, comment.UserID, comment.PhotoID, comment.Content).Scan(&comment.ID)
	if err != nil {
		return 0, err
	}
	return comment.ID, nil
}

func (r *CommentRepository) GetCommentsByPhotoID(ctx context.Context, photoID int) ([]models.Comment, error) {
	query := `
        SELECT id, user_id, photo_id, content, created_at
        FROM comments
        WHERE photo_id = $1`
	rows, err := r.DB.Query(ctx, query, photoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]models.Comment, 0)
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.PhotoID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *CommentRepository) UpdateComment(ctx context.Context, comment *models.Comment) error {
	query := `
        UPDATE comments
        SET content = $1, updated_at = NOW()
        WHERE id = $2 AND user_id = $3`
	cmdTag, err := r.DB.Exec(ctx, query, comment.Content, comment.ID, comment.UserID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no rows updated")
	}

	return nil
}

func (r *CommentRepository) DeleteComment(ctx context.Context, commentID, userID int) error {
	query := `
		DELETE FROM comments
		WHERE id = $1 AND user_id = $2`
	cmdTag, err := r.DB.Exec(ctx, query, commentID, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no rows deleted")
	}

	return err
}
