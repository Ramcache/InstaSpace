package handlers

import (
	"InstaSpace/internal/models"
	"InstaSpace/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
)

type CommentHandler struct {
	Service services.CommentServiceInterface
}

// Конструктор
func NewCommentHandler(service services.CommentServiceInterface) *CommentHandler {
	return &CommentHandler{Service: service}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.Service.CreateComment(r.Context(), &comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "comment created successfully",
		"id":      id,
	})
}

func (h *CommentHandler) GetCommentsByPhotoID(w http.ResponseWriter, r *http.Request) {
	photoID := r.URL.Query().Get("photo_id")
	if photoID == "" {
		http.Error(w, "missing photo_id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(photoID)
	if err != nil {
		http.Error(w, "invalid photo_id", http.StatusBadRequest)
		return
	}

	comments, err := h.Service.GetCommentsByPhotoID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(comments)
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.Service.UpdateComment(r.Context(), &comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "comment updated successfully"})
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID := r.URL.Query().Get("comment_id")
	userID := r.URL.Query().Get("user_id")
	if commentID == "" || userID == "" {
		http.Error(w, "missing comment_id or user_id", http.StatusBadRequest)
		return
	}

	cID, err := strconv.Atoi(commentID)
	if err != nil {
		http.Error(w, "invalid comment_id", http.StatusBadRequest)
		return
	}

	uID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteComment(r.Context(), cID, uID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "comment deleted successfully"})
}
