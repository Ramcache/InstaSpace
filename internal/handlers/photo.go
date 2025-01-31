package handlers

import (
	"InstaSpace/internal/models"
	"InstaSpace/internal/services"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type PhotoHandler struct {
	Service services.PhotoServiceInterface
	Logger  *zap.Logger
}

func NewPhotoHandler(service services.PhotoServiceInterface, logger *zap.Logger) *PhotoHandler {
	return &PhotoHandler{Service: service, Logger: logger}
}

// UploadPhoto загружает фото
//
// @Summary Загрузить фото
// @Description Загружает фото в систему и сохраняет в базе данных
// @Tags Photos
// @Accept multipart/form-data
// @Produce json
// @Param user_id header int true "ID пользователя"
// @Param file formData file true "Файл изображения"
// @Param description formData string false "Описание изображения"
// @Success 201 {object} models.Photo
// @Failure 400 {string} string "Некорректный ввод"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/photos [post]
func (h *PhotoHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Начало загрузки фото")

	userID, err := strconv.Atoi(r.Header.Get("user_id"))
	if err != nil || userID == 0 {
		h.Logger.Warn("Некорректный user_id", zap.Error(err))
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.Logger.Warn("Файл не найден", zap.Error(err))
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	const MaxUploadSize = 5 * 1024 * 1024
	if header.Size > MaxUploadSize {
		h.Logger.Warn("Файл превышает максимальный размер")
		http.Error(w, "File size exceeds 5MB", http.StatusBadRequest)
		return
	}

	allowedFormats := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	fileExt := filepath.Ext(header.Filename)
	if !allowedFormats[fileExt] {
		h.Logger.Warn("Неподдерживаемый формат файла", zap.String("format", fileExt))
		http.Error(w, "Unsupported file format", http.StatusBadRequest)
		return
	}

	uploadDir := "uploads/"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		h.Logger.Error("Не удалось создать директорию", zap.Error(err))
		http.Error(w, "Could not create directory", http.StatusInternalServerError)
		return
	}

	uniqueFileName := time.Now().Format("20060102150405") + "_" + header.Filename
	filePath := filepath.Join(uploadDir, uniqueFileName)

	if _, err := os.Stat(uploadDir); os.IsPermission(err) {
		h.Logger.Error("Недостаточно прав для записи в директорию", zap.Error(err))
		http.Error(w, "Permission denied for directory", http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(filePath)
	if err != nil {
		h.Logger.Error("Не удалось сохранить файл", zap.Error(err))
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := file.Seek(0, 0); err != nil {
		h.Logger.Error("Ошибка при обработке файла", zap.Error(err))
		http.Error(w, "File processing error", http.StatusInternalServerError)
		return
	}
	if _, err := dst.ReadFrom(file); err != nil {
		h.Logger.Error("Ошибка записи файла", zap.Error(err))
		http.Error(w, "File write error", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Файл успешно загружен", zap.String("file", filePath))

	description := r.FormValue("description")

	photo := models.Photo{
		UserID:      userID,
		URL:         filePath,
		Description: description,
	}

	if err := h.Service.SavePhoto(&photo); err != nil {
		if errors.Is(err, services.ErrInvalidPhotoData) {
			h.Logger.Warn("Некорректные данные фото", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.Logger.Error("Ошибка сохранения фото в базе данных", zap.Error(err))
		http.Error(w, "Could not save photo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(photo)
}
