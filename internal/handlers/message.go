package handlers

import (
	"InstaSpace/internal/services"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
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

// SendMessage отправляет сообщение в беседу
//
// @Summary Отправить сообщение
// @Description Отправляет новое сообщение в указанную беседу
// @Tags Messages
// @Accept json
// @Produce json
// @Param message body models.Message true "Данные сообщения"
// @Success 200 {object} map[string]int "message_id: ID созданного сообщения"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/messages [post]
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

// GetMessages возвращает сообщения из беседы
//
// @Summary Получить сообщения
// @Description Возвращает список сообщений по conversation_id
// @Tags Messages
// @Produce json
// @Param conversationID path int true "ID беседы"
// @Success 200 {array} models.Message
// @Failure 400 {string} string "Некорректный conversation_id"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/messages/{conversationID} [get]
func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationIDStr, ok := vars["conversationID"]
	if !ok || conversationIDStr == "" {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		h.Logger.Error("Invalid conversation ID")
		return
	}

	conversationID, err := strconv.Atoi(conversationIDStr)
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		h.Logger.Error("Invalid conversation ID", zap.Error(err))
		return
	}

	messages, err := h.Service.GetMessages(r.Context(), conversationID)
	if err != nil {
		if err.Error() == "conversation not found" {
			http.Error(w, "Conversation not found", http.StatusBadRequest)
			h.Logger.Warn("Conversation not found", zap.Int("conversationID", conversationID))
			return
		}
		http.Error(w, "Failed to get messages", http.StatusInternalServerError)
		h.Logger.Error("Failed to get messages", zap.Error(err))
		return
	}

	h.Logger.Info("Messages retrieved successfully", zap.Int("conversationID", conversationID), zap.Int("count", len(messages)))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}

// DeleteMessageHandler удаляет сообщение
//
// @Summary Удалить сообщение
// @Description Удаляет сообщение по его ID
// @Tags Messages
// @Produce json
// @Param messageID path int true "ID сообщения"
// @Success 200 {object} map[string]string "message: Message deleted successfully"
// @Failure 400 {string} string "Некорректный ID"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/messages/{messageID} [delete]
func (h *MessageHandler) DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	messageID, err := strconv.Atoi(mux.Vars(r)["messageID"])
	if err != nil {
		h.Logger.Error("Invalid message ID", zap.Error(err))
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteMessage(r.Context(), messageID)
	if err != nil {
		if err.Error() == "message not found" {
			http.Error(w, "Message not found", http.StatusBadRequest)
			h.Logger.Warn("Message not found", zap.Int("messageID", messageID))
			return
		}
		h.Logger.Error("Failed to delete message", zap.Error(err))
		http.Error(w, "Failed to delete message", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Message deleted successfully", zap.Int("messageID", messageID))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Message deleted successfully"})
}
