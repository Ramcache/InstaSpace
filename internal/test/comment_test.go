package test

import (
	"InstaSpace/internal/models"
	"InstaSpace/internal/services"
	"bytes"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

var (
	commentService *services.CommentService
)

func setupTestData(t *testing.T, db *pgxpool.Pool) {
	t.Helper()

	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE comments, photos, users RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO users (id, email, password, username) VALUES (1, 'test@example.com', 'test_password', 'testuser')")
	require.NoError(t, err, "Не удалось добавить запись в таблицу users")

	_, err = db.Exec(ctx, "INSERT INTO photos (id, user_id, url) VALUES (1, 1, 'http://example.com/test-photo.jpg')")
	require.NoError(t, err, "Не удалось добавить запись в таблицу photos")

	_, err = db.Exec(ctx, `
		INSERT INTO comments (user_id, photo_id, content) VALUES
		(1, 1, 'Original content'),
		(1, 1, 'Second comment')`)
	require.NoError(t, err, "Не удалось добавить записи в таблицу comments")
}

func TestCreateCommentHandler(t *testing.T) {
	setupTestData(t, db)

	testCases := []struct {
		Name         string
		Payload      models.Comment
		ExpectedCode int
		ShouldError  bool
	}{
		{
			Name: "Успешное создание комментария",
			Payload: models.Comment{
				PhotoID: 1,
				UserID:  1,
				Content: "New comment",
			},
			ExpectedCode: http.StatusCreated,
			ShouldError:  false,
		},
		{
			Name: "Ошибка: Неверный photo_id",
			Payload: models.Comment{
				PhotoID: 0,
				UserID:  1,
				Content: "Invalid photo_id",
			},
			ExpectedCode: http.StatusBadRequest,
			ShouldError:  true,
		},
		{
			Name: "Ошибка: Неверный user_id",
			Payload: models.Comment{
				PhotoID: 1,
				UserID:  0,
				Content: "Invalid user_id",
			},
			ExpectedCode: http.StatusBadRequest,
			ShouldError:  true,
		},
		{
			Name: "Ошибка: Пустой текст комментария",
			Payload: models.Comment{
				PhotoID: 1,
				UserID:  1,
				Content: "",
			},
			ExpectedCode: http.StatusBadRequest,
			ShouldError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			payload, err := json.Marshal(tc.Payload)
			require.NoError(t, err, "Ошибка сериализации payload")

			req, err := http.NewRequest("POST", testServer.URL+"/api/comments", bytes.NewReader(payload))
			require.NoError(t, err, "Ошибка создания HTTP запроса")
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err, "Ошибка выполнения HTTP запроса")
			defer resp.Body.Close()

			assert.Equal(t, tc.ExpectedCode, resp.StatusCode, "Некорректный HTTP код ответа")
		})
	}
}

func TestGetCommentsByPhotoID(t *testing.T) {
	setupTestData(t, db)

	testCases := []struct {
		Name         string
		PhotoID      string
		ExpectedCode int
		ExpectedLen  int
		ShouldError  bool
	}{
		{
			Name:         "Успешное получение комментариев",
			PhotoID:      "1",
			ExpectedCode: http.StatusOK,
			ExpectedLen:  2,
			ShouldError:  false,
		},
		{
			Name:         "Ошибка: Некорректный photoID",
			PhotoID:      "abc",
			ExpectedCode: http.StatusBadRequest,
			ExpectedLen:  0,
			ShouldError:  true,
		},
		{
			Name:         "Ошибка: Комментарии не найдены",
			PhotoID:      "99",
			ExpectedCode: http.StatusOK,
			ExpectedLen:  0,
			ShouldError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			req, err := http.NewRequest("GET", testServer.URL+"/api/comments/"+tc.PhotoID, nil)
			require.NoError(t, err, "Ошибка создания HTTP запроса")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err, "Ошибка выполнения HTTP запроса")
			defer resp.Body.Close()

			assert.Equal(t, tc.ExpectedCode, resp.StatusCode, "Некорректный HTTP код ответа")

			if !tc.ShouldError {
				var comments []struct {
					Content  string `json:"content"`
					Username string `json:"username"`
				}
				err := json.NewDecoder(resp.Body).Decode(&comments)
				require.NoError(t, err, "Ошибка декодирования ответа")

				assert.Equal(t, tc.ExpectedLen, len(comments), "Некорректное количество комментариев")

				if len(comments) > 0 {
					assert.Equal(t, "testuser", comments[0].Username, "Некорректное имя пользователя")
				}
			}
		})
	}
}

func TestUpdateComment(t *testing.T) {
	setupTestData(t, db)

	testCases := []struct {
		Name         string
		CommentID    string
		Payload      string
		ExpectedCode int
		ExpectedBody string
		ShouldError  bool
	}{
		{
			Name:         "Успешное обновление комментария",
			CommentID:    "1",
			Payload:      `{"content": "Updated content", "user_id": 1}`,
			ExpectedCode: http.StatusOK,
			ExpectedBody: `comment updated successfully`,
			ShouldError:  false,
		},
		{
			Name:         "Ошибка: Некорректный ID комментария",
			CommentID:    "abc",
			Payload:      `{"content": "Updated content", "user_id": 1}`,
			ExpectedCode: http.StatusBadRequest,
			ShouldError:  true,
		},
		{
			Name:         "Ошибка: Некорректное тело запроса",
			CommentID:    "1",
			Payload:      `{"invalid_field": "value"}`,
			ExpectedCode: http.StatusBadRequest,
			ShouldError:  true,
		},
		{
			Name:         "Ошибка: Комментарий не найден",
			CommentID:    "99",
			Payload:      `{"content": "Updated content", "user_id": 1}`,
			ExpectedCode: http.StatusInternalServerError,
			ShouldError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", testServer.URL+"/api/comments/"+tc.CommentID+"/edit", bytes.NewBuffer([]byte(tc.Payload)))
			require.NoError(t, err, "Ошибка создания HTTP запроса")
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err, "Ошибка выполнения HTTP запроса")
			defer resp.Body.Close()

			assert.Equal(t, tc.ExpectedCode, resp.StatusCode, "Некорректный HTTP код ответа")

			if !tc.ShouldError {
				var body map[string]string
				err := json.NewDecoder(resp.Body).Decode(&body)
				require.NoError(t, err, "Ошибка декодирования ответа")

				assert.Equal(t, tc.ExpectedBody, body["message"], "Некорректное сообщение ответа")
			}
		})
	}
}

func TestDeleteComment(t *testing.T) {
	setupTestData(t, db)

	testCases := []struct {
		Name         string
		CommentID    string
		UserID       string
		ExpectedCode int
		ExpectedBody string
		ShouldError  bool
	}{
		{
			Name:         "Успешное удаление комментария",
			CommentID:    "1",
			UserID:       "1",
			ExpectedCode: http.StatusOK,
			ExpectedBody: `{"message":"comment deleted successfully"}`,
			ShouldError:  false,
		},
		{
			Name:         "Ошибка: Некорректный comment_id",
			CommentID:    "abc",
			UserID:       "1",
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "invalid comment_id",
			ShouldError:  true,
		},
		{
			Name:         "Ошибка: Некорректный user_id",
			CommentID:    "1",
			UserID:       "abc",
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "invalid user_id",
			ShouldError:  true,
		},
		{
			Name:         "Ошибка: Комментарий не найден",
			CommentID:    "99",
			UserID:       "1",
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: "no rows deleted",
			ShouldError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", testServer.URL+"/api/comments/"+tc.CommentID+"/delete?user_id="+tc.UserID, nil)
			require.NoError(t, err, "Ошибка создания HTTP запроса")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err, "Ошибка выполнения HTTP запроса")
			defer resp.Body.Close()

			assert.Equal(t, tc.ExpectedCode, resp.StatusCode, "Некорректный HTTP код ответа")

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "Ошибка чтения тела ответа")

			if !tc.ShouldError {
				var response map[string]string
				err := json.Unmarshal(body, &response)
				require.NoError(t, err, "Ошибка декодирования JSON ответа")
				assert.Equal(t, "comment deleted successfully", response["message"], "Некорректное сообщение ответа")
			} else {
				assert.Contains(t, string(body), tc.ExpectedBody, "Сообщение об ошибке не совпадает")
			}
		})
	}
}
