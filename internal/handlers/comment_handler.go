package handlers

import (
	"InstaSpace/internal/models"
	"InstaSpace/internal/services"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type CommentHandler struct {
	Service services.CommentServiceInterface
	Logger  *zap.Logger
}

func NewCommentHandler(service services.CommentServiceInterface, logger *zap.Logger) *CommentHandler {
	return &CommentHandler{Service: service, Logger: logger}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Получен запрос на создание комментария")

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		h.Logger.Error("Ошибка декодирования тела запроса", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.Service.CreateComment(r.Context(), &comment)
	if err != nil {
		h.Logger.Error("Ошибка создания комментария", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Комментарий успешно создан", zap.Int("id", id))
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "comment created successfully",
		"id":      id,
	})
}

func (h *CommentHandler) GetCommentsByPhotoID(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Получен запрос на получение комментариев по photo_id")

	photoID := r.URL.Query().Get("photo_id")
	if photoID == "" {
		h.Logger.Warn("Отсутствует параметр photo_id")
		http.Error(w, "missing photo_id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(photoID)
	if err != nil {
		h.Logger.Error("Неверный формат photo_id", zap.String("photo_id", photoID), zap.Error(err))
		http.Error(w, "invalid photo_id", http.StatusBadRequest)
		return
	}

	comments, err := h.Service.GetCommentsByPhotoID(r.Context(), id)
	if err != nil {
		h.Logger.Error("Ошибка получения комментариев", zap.Int("photo_id", id), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Комментарии успешно получены", zap.Int("photo_id", id), zap.Int("count", len(comments)))
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(comments)
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Получен запрос на обновление комментария")

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		h.Logger.Error("Ошибка декодирования тела запроса", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.Service.UpdateComment(r.Context(), &comment)
	if err != nil {
		h.Logger.Error("Ошибка обновления комментария", zap.Int("id", comment.ID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Комментарий успешно обновлён", zap.Int("id", comment.ID))
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "comment updated successfully"})
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Получен запрос на удаление комментария")

	commentID := r.URL.Query().Get("comment_id")
	userID := r.URL.Query().Get("user_id")
	if commentID == "" || userID == "" {
		h.Logger.Warn("Отсутствует comment_id или user_id")
		http.Error(w, "missing comment_id or user_id", http.StatusBadRequest)
		return
	}

	cID, err := strconv.Atoi(commentID)
	if err != nil {
		h.Logger.Error("Неверный формат comment_id", zap.String("comment_id", commentID), zap.Error(err))
		http.Error(w, "invalid comment_id", http.StatusBadRequest)
		return
	}

	uID, err := strconv.Atoi(userID)
	if err != nil {
		h.Logger.Error("Неверный формат user_id", zap.String("user_id", userID), zap.Error(err))
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteComment(r.Context(), cID, uID)
	if err != nil {
		h.Logger.Error("Ошибка удаления комментария", zap.Int("comment_id", cID), zap.Int("user_id", uID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Комментарий успешно удалён", zap.Int("comment_id", cID), zap.Int("user_id", uID))
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "comment deleted successfully"})
}
