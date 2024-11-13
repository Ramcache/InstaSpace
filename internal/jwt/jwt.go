package jwt

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// GenerateJWT godoc
// @Summary Генерация JWT токена
// @Description Генерирует новый JWT токен, используя переданный email в качестве данных. Токен будет действителен в течение 24 часов.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   email path string true "Email для включения в JWT токен"
// @Success 200 {string} string "Успешно сгенерированный JWT токен"
// @Router /generate-token [get]
func GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(jwtKey)
}

// ValidateJWT godoc
// @Summary Валидация JWT токена
// @Description Проверяет переданный JWT токен и возвращает email, который был закодирован в токене, если он действителен.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   token path string true "JWT токен для проверки"
// @Success 200 {string} string "Успешная валидация токена, возвращается email"
// @Failure 401 {object} Response "Токен недействителен или истек"
// @Router /validate-token [get]

func ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["email"].(string), nil
	}

	return "", err
}
