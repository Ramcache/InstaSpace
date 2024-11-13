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

func (r *AuthRepository) CreateUser(user *models.User) error {
	query := "INSERT INTO users (email, username, password) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(context.Background(), query, user.Email, user.Username, user.Password)
	return err
}

func (r *AuthRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, email, username, password FROM users WHERE email=$1"
	row := r.db.QueryRow(context.Background(), query, email)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
