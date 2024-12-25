package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func JWTMiddleware(secret string, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				logger.Warn("Отсутствует токен или неверный формат заголовка",
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
				)
				http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				logger.Warn("Неверный токен",
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
					zap.Error(err),
				)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			logger.Info("Токен успешно проверен",
				zap.String("path", r.URL.Path),
				zap.String("method", r.Method),
			)

			next.ServeHTTP(w, r)
		})
	}
}
