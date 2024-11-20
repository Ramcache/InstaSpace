package auth

import (
	"InstaSpace/internal/jwt"
	"InstaSpace/internal/models"
	"InstaSpace/pkg/utils"
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strconv"
)

type AuthService struct {
	repo Repository
}

func NewAuthService(repo Repository) *AuthService {
	return &AuthService{repo: repo}
}

func sendConfirmationEmail(email, token string) error {
	domain := os.Getenv("APP_DOMAIN")

	confirmationLink := fmt.Sprintf("%s/confirm?token=%s", domain, token)
	messageBody := fmt.Sprintf("Пожалуйста, подтвердите вашу регистрацию, перейдя по следующей ссылке: %s", confirmationLink)

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_USER"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Подтверждение регистрации")
	m.SetBody("text/plain", messageBody)

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		log.Printf("Ошибка при преобразовании порта SMTP: %v", err)
		return err
	}

	d := gomail.NewDialer(smtpHost, port, smtpUser, smtpPassword)
	if port == 465 {
		d.SSL = true
	} else {
		d.SSL = false
	}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Ошибка отправки письма: %v", err)
		return err
	}

	log.Printf("Отправлено письмо на %s с подтверждающей ссылкой: %s\n", email, confirmationLink)
	return nil
}

func (s *AuthService) ResendConfirmationEmail(email string) error {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.IsActive {
		return errors.New("already active")
	}

	confirmationToken, err := jwt.GenerateConfirmationToken(user.Email)
	if err != nil {
		return err
	}

	err = sendConfirmationEmail(user.Email, confirmationToken)
	if err != nil {
		return err
	}

	return nil
}

// RegisterUser godoc
// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя по email, имени пользователя и паролю. Отправляет письмо для подтверждения email.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   email     body string true "Email пользователя"
// @Param   username  body string true "Имя пользователя"
// @Param   password  body string true "Пароль пользователя"
// @Success 200 {string} string "Сообщение о подтверждении регистрации"
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
		IsActive: false,
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		return "", err
	}

	confirmationToken, err := jwt.GenerateConfirmationToken(user.Email)
	if err != nil {
		return "", err
	}

	// Отправляем email с подтверждающей ссылкой
	err = sendConfirmationEmail(user.Email, confirmationToken)
	if err != nil {
		return "", err
	}

	return "Регистрация успешна! Пожалуйста, проверьте ваш email для подтверждения.", nil
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
func (s *AuthService) AuthenticateUser(email, password string) (string, string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil || !utils.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		return "", "", errors.New("account not activated")
	}

	token, err := jwt.GenerateJWT(user.Email)
	if err != nil {
		return "", "", errors.New("failed to generate token")
	}

	return token, user.Username, nil
}

// ConfirmEmail godoc
// @Summary Подтверждение email пользователя
// @Description Подтверждает email пользователя на основе предоставленного токена.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   token query string true "Токен подтверждения"
// @Success 200 {string} string "Email успешно подтверждён"
// @Failure 400 {string} string "Ошибка подтверждения"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /confirm [get]
func (s *AuthService) ConfirmEmail(token string) (string, error) {
	// Проверка токена и извлечение email
	email, err := jwt.ValidateConfirmationToken(token) // Реализуйте эту функцию в пакете jwt
	if err != nil {
		return "", errors.New("invalid or expired token")
	}

	// Активируем пользователя в базе данных
	err = s.repo.ActivateUser(email)
	if err != nil {
		return "", err
	}

	return "Email успешно подтверждён!", nil
}

func (s *AuthService) ActivateUser(email string) error {
	return s.repo.ActivateUser(email) // Реализуйте этот метод в репозитории
}

// SendConfirmationEmail отправляет подтверждающее письмо с токеном
func (s *AuthService) SendConfirmationEmail(email, token string) error {
	return sendConfirmationEmail(email, token)
}
