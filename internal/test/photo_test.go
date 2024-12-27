package test

import (
	"InstaSpace/internal/models"
	"InstaSpace/internal/services"
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

var (
	photoService *services.PhotoService
)

func TestSavePhoto(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE photos RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицу фото")

	testCases := []struct {
		Name         string
		Input        *models.Photo
		ShouldError  bool
		ExpectedCode int
	}{
		{
			Name: "Успешное сохранение фото",
			Input: &models.Photo{
				UserID:      1,
				URL:         "uploads/test1.jpg",
				Description: "Тестовое фото",
			},
			ShouldError: false,
		},
		{
			Name: "Ошибка: Некорректные данные фото",
			Input: &models.Photo{
				UserID:      0,
				URL:         "",
				Description: "",
			},
			ShouldError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := photoService.SavePhoto(tc.Input)

			if tc.ShouldError {
				require.Error(t, err, "Ожидалась ошибка, но её не было")
				return
			}

			require.NoError(t, err, "Ошибка сохранения фото")
		})
	}
}

func TestUploadPhotoHandler(t *testing.T) {
	filePath := "testdata/test_image.jpg"

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		require.NoError(t, err, "Не удалось создать директорию testdata")

		file, err := os.Create(filePath)
		require.NoError(t, err, "Не удалось создать тестовый файл")

		_, err = file.Write([]byte("test image content"))
		require.NoError(t, err, "Не удалось записать содержимое в файл")
		require.NoError(t, file.Close(), "Не удалось закрыть файл")
	}

	testCases := []struct {
		Name         string
		UserID       string
		FilePath     string
		Description  string
		ExpectedCode int
		ShouldError  bool
	}{
		{
			Name:         "Успешная загрузка фото",
			UserID:       "1",
			FilePath:     filePath,
			Description:  "Тестовое описание",
			ExpectedCode: http.StatusCreated,
			ShouldError:  false,
		},
		{
			Name:         "Ошибка: Некорректный user_id",
			UserID:       "",
			FilePath:     filePath,
			Description:  "Тестовое описание",
			ExpectedCode: http.StatusBadRequest,
			ShouldError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			file, err := os.Open(tc.FilePath)
			if tc.ShouldError && err != nil {
				return
			}
			defer file.Close()

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			part, err := writer.CreateFormFile("file", filepath.Base(tc.FilePath))
			require.NoError(t, err, "Ошибка создания файла в multipart")
			_, err = io.Copy(part, file)
			require.NoError(t, err, "Ошибка копирования файла в multipart")

			writer.WriteField("description", tc.Description)
			require.NoError(t, writer.Close(), "Ошибка закрытия writer")

			req, err := http.NewRequest("POST", testServer.URL+"/api/photos", body)
			require.NoError(t, err, "Ошибка создания запроса")
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.Header.Set("user_id", tc.UserID)

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err, "Ошибка выполнения запроса")
			defer resp.Body.Close()

			assert.Equal(t, tc.ExpectedCode, resp.StatusCode, "Некорректный HTTP код ответа")
		})
	}
}
