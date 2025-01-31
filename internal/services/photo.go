package services

import (
	"InstaSpace/internal/models"
	"InstaSpace/internal/repositories"
	"errors"
)

var ErrInvalidPhotoData = errors.New("invalid photo data: user_id or URL is missing")

type PhotoService struct {
	Repository repositories.PhotoRepositoryInterface
}

func NewPhotoService(repo repositories.PhotoRepositoryInterface) *PhotoService {
	return &PhotoService{Repository: repo}
}

type PhotoServiceInterface interface {
	SavePhoto(photo *models.Photo) error
}

func (s *PhotoService) SavePhoto(photo *models.Photo) error {
	if photo.UserID == 0 || photo.URL == "" {
		return ErrInvalidPhotoData
	}
	return s.Repository.Create(photo)
}
