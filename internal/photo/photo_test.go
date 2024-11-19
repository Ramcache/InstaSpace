package photo_test

import (
	"InstaSpace/internal/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockPhotoRepository struct {
	mock.Mock
}

func (m *MockPhotoRepository) UploadPhoto(photo *models.Photo) error {
	args := m.Called(photo)
	return args.Error(0)
}

func (m *MockPhotoRepository) GetPhotoByID(photoID uint) (*models.Photo, error) {
	args := m.Called(photoID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Photo), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPhotoRepository) DeletePhoto(photoID uint) error {
	args := m.Called(photoID)
	return args.Error(0)
}

type MockPhotoService struct {
	repo *MockPhotoRepository
}

func NewPhotoService(repo *MockPhotoRepository) *MockPhotoService {
	return &MockPhotoService{repo: repo}
}

func (s *MockPhotoService) UploadPhoto(photo *models.Photo) error {
	return s.repo.UploadPhoto(photo)
}

func (s *MockPhotoService) GetPhotoByID(photoID uint) (*models.Photo, error) {
	return s.repo.GetPhotoByID(photoID)
}

func (s *MockPhotoService) DeletePhoto(photoID uint) error {
	return s.repo.DeletePhoto(photoID)
}

func TestUploadPhoto_Success(t *testing.T) {
	mockRepo := new(MockPhotoRepository)
	photoService := NewPhotoService(mockRepo)

	// Define input
	inputPhoto := &models.Photo{
		ID:     123,
		URL:    "http://example.com/photo.jpg",
		UserID: 456,
	}

	mockRepo.On("UploadPhoto", inputPhoto).Return(nil)

	err := photoService.UploadPhoto(inputPhoto)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetPhotoByID_Success(t *testing.T) {
	mockRepo := new(MockPhotoRepository)
	photoService := NewPhotoService(mockRepo)

	photoID := uint(123)
	expectedPhoto := &models.Photo{
		ID:     photoID,
		URL:    "http://example.com/photo.jpg",
		UserID: 456,
	}

	mockRepo.On("GetPhotoByID", photoID).Return(expectedPhoto, nil)

	photo, err := photoService.GetPhotoByID(photoID)

	assert.NoError(t, err)
	assert.Equal(t, expectedPhoto, photo)
	mockRepo.AssertExpectations(t)
}

func TestGetPhotoByID_NotFound(t *testing.T) {
	mockRepo := new(MockPhotoRepository)
	photoService := NewPhotoService(mockRepo)

	photoID := uint(999)

	mockRepo.On("GetPhotoByID", photoID).Return(nil, errors.New("photo not found"))

	photo, err := photoService.GetPhotoByID(photoID)

	assert.Error(t, err)
	assert.Equal(t, "photo not found", err.Error())
	assert.Nil(t, photo)
	mockRepo.AssertExpectations(t)
}

func TestDeletePhoto_Success(t *testing.T) {
	mockRepo := new(MockPhotoRepository)
	photoService := NewPhotoService(mockRepo)

	photoID := uint(123)

	mockRepo.On("DeletePhoto", photoID).Return(nil)

	err := photoService.DeletePhoto(photoID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeletePhoto_NotFound(t *testing.T) {
	mockRepo := new(MockPhotoRepository)
	photoService := NewPhotoService(mockRepo)

	photoID := uint(999)

	mockRepo.On("DeletePhoto", photoID).Return(errors.New("photo not found"))

	err := photoService.DeletePhoto(photoID)

	assert.Error(t, err)
	assert.Equal(t, "photo not found", err.Error())
	mockRepo.AssertExpectations(t)
}
