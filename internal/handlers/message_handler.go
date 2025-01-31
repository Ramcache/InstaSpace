package handlers

import (
	"InstaSpace/internal/services"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type MessageHandler struct {
	Service *services.MessageService
	Logger  *zap.Logger
}

func NewMessageHandler(service *services.MessageService, logger *zap.Logger) *MessageHandler {
	return &MessageHandler{Service: service, Logger: logger}
}

func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ConversationID int    `json:"conversation_id"`
		SenderID       int    `json:"sender_id"`
		Content        string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		h.Logger.Error("Failed to parse request", zap.Error(err))
		return
	}

	if req.Content == "" {
		http.Error(w, "Message content cannot be empty", http.StatusBadRequest)
		h.Logger.Error("Empty message content")
		return
	}

	messageID, err := h.Service.SendMessage(r.Context(), req.ConversationID, req.SenderID, req.Content)
	if err != nil {
		if errors.Is(err, services.ErrConversationNotFound) {
			http.Error(w, "Conversation not found", http.StatusBadRequest)
			h.Logger.Error("Failed to send message - conversation not found", zap.Error(err))
			return
		}
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		h.Logger.Error("Failed to send message", zap.Error(err))
		return
	}

	h.Logger.Info("Message sent successfully", zap.Int("messageID", messageID))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"message_id": messageID})
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	conversationID, err := strconv.Atoi(r.URL.Query().Get("conversationID"))
	if err != nil || conversationID <= 0 {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		h.Logger.Error("Invalid conversation ID", zap.Error(err))
		return
	}

	messages, err := h.Service.GetMessages(r.Context(), conversationID)
	if err != nil {
		if errors.Is(err, services.ErrConversationNotFound) {
			http.Error(w, "Conversation not found", http.StatusBadRequest)
			h.Logger.Error("Failed to retrieve messages - conversation not found", zap.Error(err))
			return
		}
		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
		h.Logger.Error("Failed to retrieve messages", zap.Error(err))
		return
	}

	h.Logger.Info("Messages retrieved successfully", zap.Int("conversationID", conversationID))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}
