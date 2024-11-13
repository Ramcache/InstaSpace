// /internal/auth/auth.go
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

func (s *AuthService) AuthenticateUser(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil || !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return jwt.GenerateJWT(user.Email)
}
