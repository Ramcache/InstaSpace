package services

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"

	"InstaSpace/internal/models"
	"InstaSpace/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repository repositories.AuthRepositoryInterface
	JWTSecret  string
}

func NewAuthService(repo repositories.AuthRepositoryInterface, jwtSecret string) *AuthService {
	return &AuthService{
		Repository: repo,
		JWTSecret:  jwtSecret,
	}
}

type AuthServiceInterface interface {
	RegisterUser(user *models.User) error
	Authenticate(email, password string) (*models.User, error)
	GenerateToken(user *models.User) (string, error)
}

func (s *AuthService) RegisterUser(user *models.User) error {
	existingUser, _ := s.Repository.GetByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email уже зарегистрирован")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.Repository.Create(user)
}

func (s *AuthService) Authenticate(email, password string) (*models.User, error) {
	user, err := s.Repository.GetByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JWTSecret))
}
