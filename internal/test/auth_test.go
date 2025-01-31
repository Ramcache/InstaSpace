package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		Name         string
		Payload      string
		ExpectedCode int
		ShouldError  bool
	}{
		{
			Name:         "Успешная регистрация",
			Payload:      `{"email": "test@example.com", "password": "securepassword", "username": "testuser"}`,
			ExpectedCode: http.StatusCreated,
			ShouldError:  false,
		},
		{
			Name:         "Регистрация с существующим email",
			Payload:      `{"email": "test@example.com", "password": "securepassword", "username": "testuser2"}`,
			ExpectedCode: http.StatusInternalServerError,
			ShouldError:  true,
		},
		{
			Name:         "Регистрация без username",
			Payload:      `{"email": "test2@example.com", "password": "securepassword"}`,
			ExpectedCode: http.StatusBadRequest,
			ShouldError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := db.Exec(ctx, "TRUNCATE TABLE users RESTART IDENTITY CASCADE")
			require.NoError(t, err, "Не удалось очистить таблицу пользователей")

			if tc.Name == "Регистрация с существующим email" {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("securepassword"), bcrypt.DefaultCost)
				_, err = db.Exec(ctx, "INSERT INTO users (email, password, username) VALUES ($1, $2, $3)", "test@example.com", hashedPassword, "testuser")
				require.NoError(t, err, "Не удалось добавить существующего пользователя")
			}

			resp, err := http.Post(testServer.URL+"/register", "application/json", strings.NewReader(tc.Payload))
			require.NoError(t, err, "Ошибка отправки запроса")
			defer resp.Body.Close()

			assert.Equal(t, tc.ExpectedCode, resp.StatusCode, "Некорректный HTTP код ответа")
		})
	}
}

func TestLoginUser(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицу пользователей")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("securepassword"), bcrypt.DefaultCost)
	require.NoError(t, err, "Не удалось хэшировать пароль")
	_, err = db.Exec(ctx, "INSERT INTO users (email, password, username) VALUES ($1, $2, $3)", "test@example.com", hashedPassword, "testuser")
	require.NoError(t, err, "Не удалось добавить пользователя")

	testCases := []struct {
		Name         string
		Payload      string
		ExpectedCode int
		ShouldError  bool
	}{
		{
			Name:         "Успешный вход",
			Payload:      `{"email": "test@example.com", "password": "securepassword"}`,
			ExpectedCode: http.StatusOK,
			ShouldError:  false,
		},
		{
			Name:         "Неверный пароль",
			Payload:      `{"email": "test@example.com", "password": "wrongpassword"}`,
			ExpectedCode: http.StatusUnauthorized,
			ShouldError:  true,
		},
		{
			Name:         "Пользователь не найден",
			Payload:      `{"email": "unknown@example.com", "password": "securepassword"}`,
			ExpectedCode: http.StatusUnauthorized,
			ShouldError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := http.Post(testServer.URL+"/login", "application/json", strings.NewReader(tc.Payload))
			require.NoError(t, err, "Ошибка отправки запроса")
			defer resp.Body.Close()

			assert.Equal(t, tc.ExpectedCode, resp.StatusCode, "Некорректный HTTP код ответа")
		})
	}
}
