package auth

import (
	"context"
	"errors"

	"InstaSpace/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateUser добавляет нового пользователя в базу данных
func (r *AuthRepository) CreateUser(user *models.User) error {
	query := "INSERT INTO users (email, username, password, is_active) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(context.Background(), query, user.Email, user.Username, user.Password, user.IsActive)
	return err
}

// GetUserByEmail получает пользователя по email из базы данных
func (r *AuthRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, email, username, password, is_active FROM users WHERE email=$1"
	row := r.db.QueryRow(context.Background(), query, email)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.IsActive)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// ActivateUser обновляет статус пользователя на активный
func (r *AuthRepository) ActivateUser(email string) error {
	query := "UPDATE users SET is_active = true WHERE email = $1"
	cmdTag, err := r.db.Exec(context.Background(), query, email)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("user not found or already active")
	}
	return nil
}
