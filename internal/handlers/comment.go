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

// CreateComment создает новый комментарий
//
// @Summary Создать комментарий
// @Description Создает новый комментарий к фото
// @Tags Comments
// @Accept json
// @Produce json
// @Param comment body models.Comment true "Данные комментария"
// @Success 201 {object} map[string]interface{} "message: comment created successfully, id: 1"
// @Failure 400 {string} string "Некорректный ввод"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/comments [post]
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
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "comment created successfully",
		"id":      id,
	})
}

// GetCommentsByPhotoID возвращает список комментариев к фото
//
// @Summary Получить комментарии
// @Description Возвращает список комментариев по photo_id
// @Tags Comments
// @Produce json
// @Param photoID path int true "ID фото"
// @Success 200 {array} models.Comment
// @Failure 400 {string} string "Неверный photo_id"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/comments/{photoID} [get]
func (h *CommentHandler) GetCommentsByPhotoID(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Получен запрос на получение комментариев по photo_id")

	vars := mux.Vars(r)
	photoID := vars["photoID"]
	if photoID == "" {
		h.Logger.Warn("Отсутствует параметр photoID")
		http.Error(w, "missing photo_id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(photoID)
	if err != nil {
		h.Logger.Error("Неверный формат photoID", zap.String("photoID", photoID), zap.Error(err))
		http.Error(w, "invalid photo_id", http.StatusBadRequest)
		return
	}

	comments, err := h.Service.GetCommentsByPhotoID(r.Context(), id)
	if err != nil {
		h.Logger.Error("Ошибка получения комментариев", zap.Int("photoID", id), zap.Error(err))
		http.Error(w, "failed to get comments", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Комментарии успешно получены", zap.Int("photoID", id), zap.Int("count", len(comments)))
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(comments)
}

// UpdateComment обновляет комментарий
//
// @Summary Обновить комментарий
// @Description Обновляет текст комментария
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path int true "ID комментария"
// @Param comment body models.Comment true "Обновленные данные комментария"
// @Success 200 {object} map[string]string "message: comment updated successfully"
// @Failure 400 {string} string "Некорректные данные"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/comments/{id}/edit [put]
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
		h.Logger.Error("Ошибка обновления комментария", zap.Int("id", comment.ID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Комментарий успешно обновлён", zap.Int("id", comment.ID))
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "comment updated successfully"})
}

// DeleteComment удаляет комментарий
//
// @Summary Удалить комментарий
// @Description Удаляет комментарий по ID
// @Tags Comments
// @Produce json
// @Param id path int true "ID комментария"
// @Param user_id query int true "ID пользователя"
// @Success 200 {object} map[string]string "message: comment deleted successfully"
// @Failure 400 {string} string "Некорректный ID"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/comments/{id}/delete [delete]
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
		h.Logger.Error("Ошибка удаления комментария", zap.Int("comment_id", commentID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Комментарий успешно удалён", zap.Int("comment_id", commentID))
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "comment deleted successfully"})
}
