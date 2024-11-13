package auth

import "InstaSpace/internal/models"

type Repository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}

type Service interface {
	RegisterUser(email, username, password string) (string, error)
	AuthenticateUser(email, password string) (string, error)
}
