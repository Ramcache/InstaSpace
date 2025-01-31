package handlers

import (
	"InstaSpace/internal/repositories"
	"InstaSpace/internal/services"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type LikeHandler struct {
	Service *services.LikeService
	Logger  *zap.Logger
}

func NewLikeHandler(service *services.LikeService, logger *zap.Logger) *LikeHandler {
	return &LikeHandler{Service: service, Logger: logger}
}

// AddLikeHandler обрабатывает запрос на добавление лайка
func (h *LikeHandler) AddLikeHandler(w http.ResponseWriter, r *http.Request) {
	photoID, err := strconv.Atoi(r.URL.Query().Get("photoID"))
	if err != nil || photoID <= 0 {
		http.Error(w, "Invalid photo ID", http.StatusBadRequest)
		h.Logger.Error("Failed to parse photo ID", zap.Error(err))
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("userID"))
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		h.Logger.Error("Failed to parse user ID", zap.Error(err))
		return
	}

	err = h.Service.AddLike(r.Context(), photoID, userID)
	if errors.Is(err, repositories.ErrInvalidPhotoID) {
		http.Error(w, "Invalid photo ID", http.StatusBadRequest)
		return
	}
	if errors.Is(err, repositories.ErrInvalidUserID) {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	if err != nil {
		h.Logger.Error("Failed to add like", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Like added successfully", zap.Int("photoID", photoID), zap.Int("userID", userID))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Like added successfully"})
}

// RemoveLikeHandler обрабатывает запрос на удаление лайка
func (h *LikeHandler) RemoveLikeHandler(w http.ResponseWriter, r *http.Request) {
	photoID, err := strconv.Atoi(r.URL.Query().Get("photoID"))
	if err != nil || photoID <= 0 {
		http.Error(w, "Invalid photo ID", http.StatusBadRequest)
		h.Logger.Error("Failed to parse photo ID", zap.Error(err))
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("userID"))
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		h.Logger.Error("Failed to parse user ID", zap.Error(err))
		return
	}

	err = h.Service.RemoveLike(r.Context(), photoID, userID)
	if err != nil {
		if err.Error() == "like not found" {
			http.Error(w, "Like not found", http.StatusBadRequest)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		h.Logger.Error("Failed to remove like", zap.Error(err))
		return
	}

	h.Logger.Info("Like removed successfully", zap.Int("photoID", photoID), zap.Int("userID", userID))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Like removed successfully"})
}

// GetLikesHandler обрабатывает запрос на получение списка пользователей, поставивших лайк
func (h *LikeHandler) GetLikesHandler(w http.ResponseWriter, r *http.Request) {
	photoID, err := strconv.Atoi(r.URL.Query().Get("photoID"))
	if err != nil {
		http.Error(w, "invalid photoID", http.StatusBadRequest)
		h.Logger.Error("Invalid photoID", zap.Error(err))
		return
	}

	users, err := h.Service.GetLikes(r.Context(), photoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error("Failed to get likes", zap.Error(err))
		return
	}

	h.Logger.Info("Fetched likes successfully", zap.Int("photoID", photoID))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"users": users})
}

// GetLikeCountHandler обрабатывает запрос на получение количества лайков
func (h *LikeHandler) GetLikeCountHandler(w http.ResponseWriter, r *http.Request) {
	photoID, err := strconv.Atoi(r.URL.Query().Get("photoID"))
	if err != nil {
		http.Error(w, "invalid photoID", http.StatusBadRequest)
		h.Logger.Error("Invalid photoID", zap.Error(err))
		return
	}

	count, err := h.Service.GetLikeCount(r.Context(), photoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error("Failed to get like count", zap.Error(err))
		return
	}

	h.Logger.Info("Fetched like count successfully", zap.Int("photoID", photoID), zap.Int("likeCount", count))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"likes_count": count})
}
