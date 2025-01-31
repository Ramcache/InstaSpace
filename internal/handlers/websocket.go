package handlers

import (
	"InstaSpace/internal/services"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WebSocketHandler struct {
	Clients        map[*websocket.Conn]bool
	Mutex          sync.Mutex
	Logger         *zap.Logger
	MessageService *services.MessageService
}

func NewWebSocketHandler(logger *zap.Logger, messageService *services.MessageService) *WebSocketHandler {
	return &WebSocketHandler{
		Clients:        make(map[*websocket.Conn]bool),
		Logger:         logger,
		MessageService: messageService,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleWS обрабатывает WebSocket соединение
//
// @Summary Установить WebSocket соединение
// @Description Устанавливает WebSocket соединение и отправляет/получает сообщения в режиме реального времени
// @Tags WebSocket
// @Accept json
// @Produce json
// @Success 101 {string} string "WebSocket connection established"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /ws [get]
func (h *WebSocketHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Logger.Error("WebSocket upgrade failed", zap.Error(err))
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	h.Mutex.Lock()
	h.Clients[conn] = true
	h.Mutex.Unlock()

	defer func() {
		h.Mutex.Lock()
		delete(h.Clients, conn)
		h.Mutex.Unlock()
		h.Logger.Info("WebSocket connection closed")
	}()

	h.Logger.Info("New WebSocket connection established")

	for {
		var msg struct {
			ConversationID int    `json:"conversation_id"`
			SenderID       int    `json:"sender_id"`
			Content        string `json:"content"`
		}

		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.Logger.Error("Unexpected WebSocket disconnection", zap.Error(err))
			} else {
				h.Logger.Warn("WebSocket client disconnected")
			}
			break
		}

		if msg.ConversationID == 0 || msg.SenderID == 0 || msg.Content == "" {
			h.Logger.Warn("Received invalid message data", zap.Any("message", msg))
			_ = conn.WriteJSON(map[string]string{"error": "Invalid message data"})
			continue
		}

		messageID, err := h.MessageService.SendMessage(r.Context(), msg.ConversationID, msg.SenderID, msg.Content)
		if err != nil {
			h.Logger.Error("Failed to save message", zap.Error(err))
			_ = conn.WriteJSON(map[string]string{"error": "Failed to save message"})
			continue
		}

		h.Logger.Info("Message received and saved",
			zap.Int("conversation_id", msg.ConversationID),
			zap.Int("sender_id", msg.SenderID),
			zap.String("content", msg.Content),
			zap.Int("message_id", messageID),
		)

		response := map[string]interface{}{
			"message_id":      messageID,
			"conversation_id": msg.ConversationID,
			"sender_id":       msg.SenderID,
			"content":         msg.Content,
		}

		h.Mutex.Lock()
		for client := range h.Clients {
			err = client.WriteJSON(response)
			if err != nil {
				h.Logger.Error("Failed to send message", zap.Error(err))
				client.Close()
				delete(h.Clients, client)
			}
		}
		h.Mutex.Unlock()
	}
}
