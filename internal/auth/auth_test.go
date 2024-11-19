package auth_test

import (
	"InstaSpace/internal/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) FindUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

type MockAuthService struct {
	repo *MockAuthRepository
}

func NewService(repo *MockAuthRepository) *MockAuthService {
	return &MockAuthService{repo: repo}
}

func (s *MockAuthService) RegisterUser(user *models.User) error {
	existingUser, err := s.repo.FindUserByEmail(user.Email)
	if existingUser != nil {
		return errors.New("user already exists")
	}
	if err != nil {
		return err
	}
	return s.repo.CreateUser(user)
}

func (s *MockAuthService) LoginUser(email, password string) (string, error) {
	storedUser, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return "", err
	}
	if !CheckPasswordHash(password, storedUser.Password) {
		return "", errors.New("invalid password")
	}
	token := "dummy-token"
	return token, nil
}

func HashPassword(password string) string {
	return "hashed-" + password
}

func CheckPasswordHash(password, hash string) bool {
	return hash == "hashed-"+password
}

func TestRegisterUser(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	authService := NewService(mockRepo)

	inputUser := &models.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.On("FindUserByEmail", "test@example.com").Return(nil, nil)
	mockRepo.On("CreateUser", inputUser).Return(nil)

	err := authService.RegisterUser(inputUser)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_ExistingEmail(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	authService := NewService(mockRepo)

	existingUser := &models.User{
		Email:    "existing@example.com",
		Password: "password123",
	}

	mockRepo.On("FindUserByEmail", "existing@example.com").Return(existingUser, nil)

	err := authService.RegisterUser(existingUser)

	assert.Error(t, err)
	assert.Equal(t, "user already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_Success(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	authService := NewService(mockRepo)

	inputUser := &models.User{
		Email:    "test@example.com",
		Password: "password123",
	}
	storedUser := &models.User{
		Email:    "test@example.com",
		Password: HashPassword("password123"),
	}

	mockRepo.On("FindUserByEmail", "test@example.com").Return(storedUser, nil)

	token, err := authService.LoginUser(inputUser.Email, inputUser.Password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_InvalidPassword(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	authService := NewService(mockRepo)

	inputUser := &models.User{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	storedUser := &models.User{
		Email:    "test@example.com",
		Password: HashPassword("password123"),
	}

	mockRepo.On("FindUserByEmail", "test@example.com").Return(storedUser, nil)

	token, err := authService.LoginUser(inputUser.Email, inputUser.Password)

	assert.Error(t, err)
	assert.Equal(t, "invalid password", err.Error())
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_UserNotFound(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	authService := NewService(mockRepo)

	email := "nonexistent@example.com"

	mockRepo.On("FindUserByEmail", email).Return(nil, errors.New("user not found"))

	token, err := authService.LoginUser(email, "password123")

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestHashPassword(t *testing.T) {
	password := "password123"
	hashedPassword := HashPassword(password)

	assert.NotEqual(t, password, hashedPassword)
	assert.True(t, CheckPasswordHash(password, hashedPassword))
}

func TestCheckPasswordHash(t *testing.T) {
	password := "password123"
	hashedPassword := HashPassword(password)

	isValid := CheckPasswordHash(password, hashedPassword)
	assert.True(t, isValid)

	isInvalid := CheckPasswordHash("wrongpassword", hashedPassword)
	assert.False(t, isInvalid)
}
