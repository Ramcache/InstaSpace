package middleware

import (
	"context"
	"net/http"
	"strings"

	"InstaSpace/internal/jwt"
)

type contextKey string

const UserContextKey contextKey = "user"

// JWTAuth - middleware для проверки JWT токена.
// @Summary Проверка JWT токена в заголовке Authorization
// @Description Эта middleware функция проверяет наличие и валидность JWT токена в заголовке запроса. Если токен валиден, в контекст запроса добавляется email пользователя. В противном случае возвращается ошибка "Unauthorized".
// @Tags middleware
// @Accept  */*
// @Produce  json
// @Param Authorization header string true "Bearer {токен}"
// @Success 200 "Запрос обработан успешно"
// @Failure 401 {object} map[string]string "Неавторизован: некорректный токен или отсутствующий заголовок"
// @Router / [middleware]

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		email, err := jwt.ValidateJWT(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
