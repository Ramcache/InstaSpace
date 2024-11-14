package auth

import "InstaSpace/internal/models"

// Repository представляет интерфейс для операций с базой данных
type Repository interface {
	CreateUser(user *models.User) error                // Создает нового пользователя
	GetUserByEmail(email string) (*models.User, error) // Получает пользователя по email
	ActivateUser(email string) error                   // Активирует пользователя по email
}

// Service представляет интерфейс для бизнес-логики авторизации
type Service interface {
	ResendConfirmationEmail(email string) error                    // Повторно отправляет подтверждающее письмо
	RegisterUser(email, username, password string) (string, error) // Регистрирует нового пользователя и отправляет подтверждение
	AuthenticateUser(email, password string) (string, error)       // Аутентифицирует пользователя и возвращает токен
	SendConfirmationEmail(email, token string) error               // Отправляет подтверждающее письмо по email
	ActivateUser(email string) error                               // Активирует пользователя, используя токен
}
