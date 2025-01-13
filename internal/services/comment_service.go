package services

import (
	"InstaSpace/internal/models"
	"InstaSpace/internal/repositories"
	"context"
	"errors"
)

type CommentService struct {
	Repo repositories.CommentRepositoryInterface
}

func NewCommentService(repo repositories.CommentRepositoryInterface) *CommentService {
	return &CommentService{Repo: repo}
}

type CommentServiceInterface interface {
	CreateComment(ctx context.Context, comment *models.Comment) (int, error)
	GetCommentsByPhotoID(ctx context.Context, photoID int) ([]models.Comment, error)
	UpdateComment(ctx context.Context, comment *models.Comment) error
	DeleteComment(ctx context.Context, commentID, userID int) error
}

func (s *CommentService) CreateComment(ctx context.Context, comment *models.Comment) (int, error) {
	// Простая валидация
	if comment.UserID == 0 || comment.PhotoID == 0 || comment.Content == "" {
		return 0, errors.New("missing required fields")
	}

	// Создаем комментарий через репозиторий
	return s.Repo.CreateComment(ctx, comment)
}

func (s *CommentService) GetCommentsByPhotoID(ctx context.Context, photoID int) ([]models.Comment, error) {
	if photoID <= 0 {
		return nil, errors.New("invalid photo ID")
	}

	return s.Repo.GetCommentsByPhotoID(ctx, photoID)
}

func (s *CommentService) UpdateComment(ctx context.Context, comment *models.Comment) error {
	if comment.ID <= 0 || comment.UserID <= 0 || comment.Content == "" {
		return errors.New("invalid comment data")
	}

	return s.Repo.UpdateComment(ctx, comment)
}

func (s *CommentService) DeleteComment(ctx context.Context, commentID, userID int) error {
	if commentID <= 0 || userID <= 0 {
		return errors.New("invalid IDs")
	}

	return s.Repo.DeleteComment(ctx, commentID, userID)
}
