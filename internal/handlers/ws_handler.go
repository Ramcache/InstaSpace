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

// WebSocket-соединение
func (h *WebSocketHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Logger.Error("WebSocket upgrade failed", zap.Error(err))
		return
	}
	defer conn.Close()

	h.Mutex.Lock()
	h.Clients[conn] = true
	h.Mutex.Unlock()

	h.Logger.Info("New WebSocket connection established")

	for {
		var msg struct {
			ConversationID int    `json:"conversation_id"`
			SenderID       int    `json:"sender_id"`
			Content        string `json:"content"`
		}

		err := conn.ReadJSON(&msg)
		if err != nil {
			h.Logger.Error("Failed to read message", zap.Error(err))
			break
		}

		// Сохраняем сообщение в БД
		_, err = h.MessageService.SendMessage(r.Context(), msg.ConversationID, msg.SenderID, msg.Content)
		if err != nil {
			h.Logger.Error("Failed to save message", zap.Error(err))
			continue
		}

		h.Logger.Info("Message received and saved", zap.Int("conversation_id", msg.ConversationID), zap.Int("sender_id", msg.SenderID), zap.String("content", msg.Content))

		// Рассылаем сообщение всем клиентам
		h.Mutex.Lock()
		for client := range h.Clients {
			err = client.WriteJSON(msg)
			if err != nil {
				h.Logger.Error("Failed to send message", zap.Error(err))
				client.Close()
				delete(h.Clients, client)
			}
		}
		h.Mutex.Unlock()
	}
}
