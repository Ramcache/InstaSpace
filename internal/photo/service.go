package photo

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type HTTPError struct {
	Code    int    `json:"code"`    // Код ошибки
	Message string `json:"message"` // Сообщение об ошибке
}

// Service defines methods for managing photos.
type Service struct {
	repo      *Repository
	uploadDir string
}

func NewService(repo *Repository, uploadDir string) *Service {
	return &Service{repo: repo, uploadDir: uploadDir}
}

func (s *Service) GetPhotoByID(ctx context.Context, id int) (*Photo, error) {
	return s.repo.GetPhotoByID(ctx, id)
}

func (s *Service) GetUserIDByEmail(ctx context.Context, email string) (int, error) {
	return s.repo.GetUserIDByEmail(ctx, email)
}

func (s *Service) UploadPhoto(ctx context.Context, userID int, fileHeader *multipart.FileHeader, description string) (*Photo, error) {
	allowedFormats := []string{".jpg", ".jpeg", ".png"}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !contains(allowedFormats, ext) {
		return nil, errors.New("unsupported file format")
	}

	maxSize := int64(2 << 20) // 2 MB
	if fileHeader.Size > maxSize {
		return nil, errors.New("file too large")
	}

	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	filePath := filepath.Join(s.uploadDir, fileName)

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, err
	}

	filePath = strings.ReplaceAll(filePath, "\\", "/")
	if description == "" {
		description = "Описание не указано"
	}

	photo := &Photo{
		UserID:      userID,
		ImageURL:    filePath,
		Description: description,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.SavePhoto(ctx, photo); err != nil {
		return nil, err
	}

	return photo, nil
}

// contains проверяет, содержится ли строка в массиве строк
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
