package test

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func setupTestLikes(t *testing.T, db *pgxpool.Pool) {
	t.Helper()

	ctx := context.Background()

	// Очистка таблиц
	_, err := db.Exec(ctx, "TRUNCATE TABLE comments, photos, photo_likes, users RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	// Вставка тестового пользователя
	_, err = db.Exec(ctx, `
		INSERT INTO users (id, email, password, username) 
		VALUES (1, 'test@example.com', 'test_password', 'testuser')
	`)
	require.NoError(t, err, "Не удалось добавить запись в таблицу users")

	// Вставка тестовой фотографии
	_, err = db.Exec(ctx, `
		INSERT INTO photos (id, user_id, url, description) 
		VALUES (1, 1, 'http://example.com/test-photo.jpg', 'Test photo')
	`)
	require.NoError(t, err, "Не удалось добавить запись в таблицу photos")

	// Вставка комментариев
	_, err = db.Exec(ctx, `
		INSERT INTO comments (user_id, photo_id, content) 
		VALUES (1, 1, 'Original comment'), (1, 1, 'Second comment')
	`)
	require.NoError(t, err, "Не удалось добавить записи в таблицу comments")
}

func TestLikeHandlers(t *testing.T) {
	setupTestLikes(t, db)

	testCases := []struct {
		Name         string
		Method       string
		URL          string
		ExpectedCode int
		Payload      map[string]string
		ShouldError  bool
	}{
		{
			Name:         "Успешное добавление лайка",
			Method:       "POST",
			URL:          "/api/likes?photoID=1&userID=1",
			ExpectedCode: http.StatusOK,
			ShouldError:  false,
		},
		{
			Name:         "Ошибка: Неверный photoID",
			Method:       "POST",
			URL:          "/api/likes?photoID=0&userID=1",
			ExpectedCode: http.StatusBadRequest,
			ShouldError:  true,
		},
		{
			Name:         "Ошибка: Неверный userID",
			Method:       "POST",
			URL:          "/api/likes?photoID=1&userID=0",
			ExpectedCode: http.StatusBadRequest,
			ShouldError:  true,
		},
		{
			Name:         "Успешное удаление лайка",
			Method:       "DELETE",
			URL:          "/api/likes?photoID=1&userID=1",
			ExpectedCode: http.StatusOK,
			ShouldError:  false,
		},
		{
			Name:         "Ошибка: Удаление отсутствующего лайка",
			Method:       "DELETE",
			URL:          "/api/likes?photoID=1&userID=2",
			ExpectedCode: http.StatusBadRequest,
			ShouldError:  true,
		},

		{
			Name:         "Получение списка пользователей, поставивших лайки",
			Method:       "GET",
			URL:          "/api/likes?photoID=1",
			ExpectedCode: http.StatusOK,
			ShouldError:  false,
		},
		{
			Name:         "Получение количества лайков",
			Method:       "GET",
			URL:          "/api/likes/count?photoID=1",
			ExpectedCode: http.StatusOK,
			ShouldError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tc.Method == "POST" || tc.Method == "DELETE" {
				req, err = http.NewRequest(tc.Method, testServer.URL+tc.URL, nil)
			} else {
				req, err = http.NewRequest(tc.Method, testServer.URL+tc.URL, nil)
			}
			require.NoError(t, err, "Ошибка создания HTTP запроса")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err, "Ошибка выполнения HTTP запроса")
			defer resp.Body.Close()

			assert.Equal(t, tc.ExpectedCode, resp.StatusCode, "Некорректный HTTP код ответа")

			if !tc.ShouldError {
				var result map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				require.NoError(t, err, "Ошибка декодирования ответа")
			}
		})
	}
}
