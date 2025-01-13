package repositories

import (
	"InstaSpace/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB: db}
}

type AuthRepositoryInterface interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
}

func (r *UserRepository) Create(user *models.User) error {
	query := "INSERT INTO users (email, password, verified) VALUES ($1, $2, $3) RETURNING id"
	return r.DB.QueryRow(context.Background(), query, user.Email, user.Password, user.Verified).Scan(&user.ID)
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := "SELECT id, email, password, verified FROM users WHERE email = $1"
	row := r.DB.QueryRow(context.Background(), query, email)

	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Verified); err != nil {
		return nil, err
	}

	return &user, nil
}
