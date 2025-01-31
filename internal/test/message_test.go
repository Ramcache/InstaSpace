package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
)

func TestSendMessage(t *testing.T) {
	ctx := context.Background()

	// Очистка таблиц перед тестами
	_, err := db.Exec(ctx, "TRUNCATE TABLE messages, conversations, users RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	// Добавление тестовых пользователей
	_, err = db.Exec(ctx, `
		INSERT INTO users (id, email, password, username) VALUES 
		(1, 'user1@example.com', 'hashedpassword1', 'user1'),
		(2, 'user2@example.com', 'hashedpassword2', 'user2')
	`)
	require.NoError(t, err, "Не удалось создать тестовых пользователей")

	// Добавление тестовой переписки
	_, err = db.Exec(ctx, "INSERT INTO conversations (id, user1_id, user2_id) VALUES (1, 1, 2)")
	require.NoError(t, err, "Не удалось создать тестовую переписку")

	testCases := []struct {
		Name         string
		Payload      string
		ExpectedCode int
		ShouldError  bool
	}{
		{
			Name:         "Успешная отправка сообщения",
			Payload:      `{"conversation_id": 1, "sender_id": 1, "content": "Hello!"}`,
			ExpectedCode: http.StatusOK,
			ShouldError:  false,
		},
		{
			Name:         "Ошибка: Пустое сообщение",
			Payload:      `{"conversation_id": 1, "sender_id": 1, "content": ""}`,
			ExpectedCode: http.StatusBadRequest,
			ShouldError:  true,
		},
		{
			Name:         "Ошибка: Отправка в несуществующую переписку",
			Payload:      `{"conversation_id": 999, "sender_id": 1, "content": "Test"}`,
			ExpectedCode: http.StatusBadRequest,
			ShouldError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := http.Post(testServer.URL+"/api/messages", "application/json", strings.NewReader(tc.Payload))
			require.NoError(t, err, "Ошибка отправки запроса")
			defer resp.Body.Close()

			assert.Equal(t, tc.ExpectedCode, resp.StatusCode, "Некорректный HTTP код ответа")
		})
	}
}

func TestGetMessages(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE messages, conversations, users RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, `
		INSERT INTO users (id, email, password, username) VALUES 
		(1, 'user1@example.com', 'hashedpassword1', 'user1'),
		(2, 'user2@example.com', 'hashedpassword2', 'user2')
	`)
	require.NoError(t, err, "Не удалось создать тестовых пользователей")

	_, err = db.Exec(ctx, "INSERT INTO conversations (id, user1_id, user2_id) VALUES (1, 1, 2)")
	require.NoError(t, err, "Не удалось создать тестовую переписку")

	_, err = db.Exec(ctx, "INSERT INTO messages (conversation_id, sender_id, content) VALUES ($1, $2, $3)", 1, 1, "Test message")
	require.NoError(t, err, "Не удалось добавить сообщение")

	testCases := []struct {
		Name         string
		URL          string
		ExpectedCode int
		ShouldError  bool
	}{
		{
			Name:         "Успешное получение сообщений",
			URL:          "/api/messages/1", // ✅ Исправлено
			ExpectedCode: http.StatusOK,
			ShouldError:  false,
		},
		{
			Name:         "Ошибка: Переписка не найдена",
			URL:          "/api/messages/999", // ✅ Ошибочный ID в URL
			ExpectedCode: http.StatusBadRequest,
			ShouldError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := http.Get(testServer.URL + tc.URL) // ✅ Передаем conversationID в URL
			require.NoError(t, err, "Ошибка выполнения HTTP запроса")
			defer resp.Body.Close()

			assert.Equal(t, tc.ExpectedCode, resp.StatusCode, "Некорректный HTTP код ответа")
		})
	}
}

func TestDeleteMessage(t *testing.T) {
	ctx := context.Background()

	// Очистка таблиц перед тестами
	_, err := db.Exec(ctx, "TRUNCATE TABLE messages, conversations, users RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	// Добавление тестовых пользователей
	_, err = db.Exec(ctx, `
		INSERT INTO users (id, email, password, username) VALUES 
		(1, 'user1@example.com', 'hashedpassword1', 'user1'),
		(2, 'user2@example.com', 'hashedpassword2', 'user2')
	`)
	require.NoError(t, err, "Не удалось создать тестовых пользователей")

	// Добавление тестовой переписки
	_, err = db.Exec(ctx, "INSERT INTO conversations (id, user1_id, user2_id) VALUES (1, 1, 2)")
	require.NoError(t, err, "Не удалось создать тестовую переписку")

	// Вставка тестового сообщения
	_, err = db.Exec(ctx, "INSERT INTO messages (id, conversation_id, sender_id, content) VALUES (1, 1, 1, 'Test message')")
	require.NoError(t, err, "Не удалось добавить сообщение")

	testCases := []struct {
		Name         string
		URL          string
		Method       string
		ExpectedCode int
		ShouldError  bool
	}{
		{
			Name:         "Успешное удаление сообщения",
			URL:          "/api/messages/1",
			Method:       "DELETE",
			ExpectedCode: http.StatusOK,
			ShouldError:  false,
		},
		{
			Name:         "Ошибка: Удаление несуществующего сообщения",
			URL:          "/api/messages/999",
			Method:       "DELETE",
			ExpectedCode: http.StatusBadRequest,
			ShouldError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			req, err := http.NewRequest(tc.Method, testServer.URL+tc.URL, nil)
			require.NoError(t, err, "Ошибка создания HTTP запроса")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err, "Ошибка выполнения HTTP запроса")
			defer resp.Body.Close()

			assert.Equal(t, tc.ExpectedCode, resp.StatusCode, "Некорректный HTTP код ответа")
		})
	}
}
