package test

import (
	"InstaSpace/internal/handlers"
	"context"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var (
	wsHandler *handlers.WebSocketHandler
)

func TestWebSocketMessaging(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE messages, conversations, users RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, `INSERT INTO users (id, email, password, username) VALUES 
		(1, 'user1@example.com', 'password1', 'user1'), 
		(2, 'user2@example.com', 'password2', 'user2')`)
	require.NoError(t, err, "Не удалось создать пользователей")

	_, err = db.Exec(ctx, "INSERT INTO conversations (id, user1_id, user2_id) VALUES (1, 1, 2)")
	require.NoError(t, err, "Не удалось создать тестовую переписку")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wsHandler.HandleWS(w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "Не удалось подключиться к WebSocket")
	defer conn.Close()

	testMessage := map[string]interface{}{
		"conversation_id": 1,
		"sender_id":       1,
		"content":         "Hello WebSocket!",
	}

	err = conn.WriteJSON(testMessage)
	require.NoError(t, err, "Ошибка при отправке WebSocket-сообщения")

	var response map[string]interface{}
	err = conn.ReadJSON(&response)
	require.NoError(t, err, "Ошибка при чтении WebSocket-сообщения")

	// Проверяем, что ответ содержит message_id и совпадает с отправленным содержимым
	assert.Equal(t, testMessage["content"], response["content"], "Сообщения не совпадают")
	assert.Contains(t, response, "message_id", "Ожидалось, что сообщение содержит message_id")
}

func TestWebSocketInvalidMessage(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wsHandler.HandleWS(w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "Не удалось подключиться к WebSocket")
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("invalid data"))
	require.NoError(t, err, "Ошибка при отправке некорректного сообщения")

	time.Sleep(1 * time.Second)

	_, _, err = conn.ReadMessage()
	assert.Error(t, err, "Ожидалось закрытие соединения из-за некорректных данных")
}

func TestWebSocketMultipleClients(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wsHandler.HandleWS(w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	conn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "Не удалось подключиться к WebSocket клиенту 1")
	defer conn1.Close()

	conn2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "Не удалось подключиться к WebSocket клиенту 2")
	defer conn2.Close()

	testMessage := map[string]interface{}{
		"conversation_id": 1,
		"sender_id":       1,
		"content":         "Hello from Client 1!",
	}

	err = conn1.WriteJSON(testMessage)
	require.NoError(t, err, "Ошибка при отправке WebSocket-сообщения клиентом 1")

	var response map[string]interface{}
	err = conn2.ReadJSON(&response)
	require.NoError(t, err, "Ошибка при чтении WebSocket-сообщения клиентом 2")

	assert.Equal(t, testMessage["content"], response["content"], "Сообщения не совпадают")
	assert.Contains(t, response, "message_id", "Ожидалось, что сообщение содержит message_id")
}
