package services

import (
	"InstaSpace/internal/models"
	"InstaSpace/internal/repositories"
	"context"
)

type LikeService struct {
	Repo *repositories.LikeRepository
}

func NewLikeService(repo *repositories.LikeRepository) *LikeService {
	return &LikeService{Repo: repo}
}

// AddLike добавляет лайк к фотографии
func (s *LikeService) AddLike(ctx context.Context, photoID, userID int) error {
	return s.Repo.AddLike(ctx, photoID, userID)
}

// RemoveLike удаляет лайк с фотографии
func (s *LikeService) RemoveLike(ctx context.Context, photoID, userID int) error {
	return s.Repo.RemoveLike(ctx, photoID, userID)
}

// GetLikes возвращает список пользователей, которые поставили лайк
func (s *LikeService) GetLikes(ctx context.Context, photoID int) ([]models.User, error) {
	return s.Repo.GetLikes(ctx, photoID)
}

// GetLikeCount возвращает количество лайков на фотографии
func (s *LikeService) GetLikeCount(ctx context.Context, photoID int) (int, error) {
	return s.Repo.GetLikeCount(ctx, photoID)
}
