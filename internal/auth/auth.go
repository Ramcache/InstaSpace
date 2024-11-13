package auth

import (
	"InstaSpace/internal/jwt"
	"InstaSpace/internal/models"
	"InstaSpace/pkg/utils"
	"errors"
)

type AuthService struct {
	repo Repository
}

func NewAuthService(repo Repository) *AuthService {
	return &AuthService{repo: repo}
}

// RegisterUser godoc
// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя по email, имени пользователя и паролю. Возвращает JWT токен после успешной регистрации.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   email     body string true "Email пользователя"
// @Param   username  body string true "Имя пользователя"
// @Param   password  body string true "Пароль пользователя"
// @Success 200 {string} string "JWT токен после успешной регистрации"
// @Failure 400 {string} string "Ошибка регистрации"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /register [post]
func (s *AuthService) RegisterUser(email, username, password string) (string, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	user := &models.User{
		Email:    email,
		Username: username,
		Password: hashedPassword,
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		return "", err
	}

	token, err := jwt.GenerateJWT(user.Email)
	return token, err
}

// AuthenticateUser godoc
// @Summary Аутентификация пользователя
// @Description Выполняет аутентификацию пользователя по email и паролю. Возвращает JWT токен, если учетные данные верны.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   email     body string true "Email пользователя"
// @Param   password  body string true "Пароль пользователя"
// @Success 200 {string} string "JWT токен после успешного входа"
// @Failure 400 {string} string "Неверные учетные данные"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /login [post]
func (s *AuthService) AuthenticateUser(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil || !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return jwt.GenerateJWT(user.Email)
}
