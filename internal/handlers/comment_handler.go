package handlers

import (
	"InstaSpace/internal/models"
	"InstaSpace/internal/services"
	"encoding/json"
	"github.com/gorilla/mux"
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

	if comment.PhotoID <= 0 || comment.UserID <= 0 || len(comment.Content) == 0 {
		h.Logger.Warn("Некорректные данные для комментария", zap.Any("comment", comment))
		http.Error(w, "invalid input data", http.StatusBadRequest)
		return
	}

	id, err := h.Service.CreateComment(r.Context(), &comment)
	if err != nil {
		if err == services.ErrInvalidForeignKey {
			h.Logger.Warn("Ошибка внешнего ключа", zap.Error(err))
			http.Error(w, "invalid photo_id or user_id", http.StatusBadRequest)
		} else {
			h.Logger.Error("Ошибка создания комментария", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.Logger.Info("Комментарий успешно создан", zap.Int("id", id))
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "comment created successfully",
		"id":      id,
	}); err != nil {
		h.Logger.Error("Ошибка кодирования ответа", zap.Error(err))
	}
}

func (h *CommentHandler) GetCommentsByPhotoID(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Получен запрос на получение комментариев по photo_id")

	vars := mux.Vars(r)
	photoID := vars["photoID"]
	if photoID == "" {
		h.Logger.Warn("Отсутствует параметр photoID")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing photo_id"})
		return
	}

	id, err := strconv.Atoi(photoID)
	if err != nil {
		h.Logger.Error("Неверный формат photoID", zap.String("photoID", photoID), zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid photo_id"})
		return
	}

	comments, err := h.Service.GetCommentsByPhotoID(r.Context(), id)
	if err != nil {
		h.Logger.Error("Ошибка получения комментариев", zap.Int("photoID", id), zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get comments"})
		return
	}

	h.Logger.Info("Комментарии успешно получены", zap.Int("photoID", id), zap.Int("count", len(comments)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(comments)
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Получен запрос на обновление комментария")

	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.Logger.Error("Некорректный ID комментария", zap.Error(err))
		http.Error(w, "invalid comment ID", http.StatusBadRequest)
		return
	}

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		h.Logger.Error("Ошибка декодирования тела запроса", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	comment.ID = commentID

	err = h.Service.UpdateComment(r.Context(), &comment)
	if err != nil {
		if err.Error() == "invalid comment data" {
			h.Logger.Error("Некорректные данные комментария", zap.Error(err))
			http.Error(w, "invalid comment data", http.StatusBadRequest)
			return
		}
		h.Logger.Error("Ошибка обновления комментария", zap.Int("id", comment.ID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Комментарий успешно обновлён", zap.Int("id", comment.ID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "comment updated successfully"})
	if err != nil {
		h.Logger.Error("Ошибка при кодировании JSON ответа", zap.Error(err))
	}
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Получен запрос на удаление комментария")

	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.Logger.Error("Неверный формат comment_id", zap.String("comment_id", vars["id"]), zap.Error(err))
		http.Error(w, "invalid comment_id", http.StatusBadRequest)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		h.Logger.Warn("Отсутствует user_id")
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.Logger.Error("Неверный формат user_id", zap.String("user_id", userIDStr), zap.Error(err))
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteComment(r.Context(), commentID, userID)
	if err != nil {
		if err.Error() == "no rows deleted" {
			h.Logger.Warn("Комментарий не найден", zap.Int("comment_id", commentID), zap.Int("user_id", userID))
			http.Error(w, "no rows deleted", http.StatusInternalServerError)
			return
		}
		h.Logger.Error("Ошибка удаления комментария", zap.Int("comment_id", commentID), zap.Int("user_id", userID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Комментарий успешно удалён", zap.Int("comment_id", commentID), zap.Int("user_id", userID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "comment deleted successfully"})
}
