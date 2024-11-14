package photo

import (
	"InstaSpace/pkg/middleware"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"path"
	"strconv"
)

type Handler struct {
	service *Service
	baseURL string
}

func NewHandler(service *Service, baseURL string) *Handler {
	return &Handler{service: service, baseURL: baseURL}
}

// UploadPhoto handles the photo upload
// @Summary Uploads a new photo
// @Description Upload a photo with an optional description. The uploaded photo is associated with the authenticated user.
// @Tags Photos
// @Accept multipart/form-data
// @Produce json
// @Param photo formData file true "Photo file to upload"
// @Param Description formData string false "Description for the photo"
// @Success 200 {object} Photo "Returns the uploaded photo details including the generated URL"
// @Failure 400 {string} string "Error uploading file or invalid input"
// @Failure 401 {string} string "Unauthorized access"
// @Failure 500 {string} string "Internal server error"
// @Router /photos/upload [post]
func (h *Handler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok || email == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := h.service.GetUserIDByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	file, fileHeader, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Error uploading file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	description := r.PostFormValue("Description")
	if description == "" {
		description = "Описание не указано"
	}

	log.Printf("Полученное описание: %s", description)

	photo, err := h.service.UploadPhoto(r.Context(), userID, fileHeader, description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	photo.ImageURL = path.Join(h.baseURL, photo.ImageURL)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(photo)
}

// GetPhoto handles fetching photo by ID
// @Summary Retrieves a photo by ID
// @Description Fetch photo details by its unique identifier.
// @Tags Photos
// @Produce json
// @Param id path int true "ID of the photo to retrieve"
// @Success 200 {object} Photo "Returns the photo details including the full URL"
// @Failure 400 {string} string "Invalid photo ID"
// @Failure 404 {string} string "Photo not found"
// @Router /photos/{id} [get]
func (h *Handler) GetPhoto(w http.ResponseWriter, r *http.Request) {
	// Получаем id из переменной пути
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Photo ID is missing in the URL", http.StatusBadRequest)
		return
	}

	// Преобразуем id из строки в число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid photo ID", http.StatusBadRequest)
		return
	}

	// Получаем фото из базы данных
	photo, err := h.service.GetPhotoByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Photo not found", http.StatusNotFound)
		return
	}

	// Создаем полный URL для изображения
	photo.ImageURL = path.Join(h.baseURL, photo.ImageURL)

	// Возвращаем JSON-ответ с информацией о фото
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(photo)
}
