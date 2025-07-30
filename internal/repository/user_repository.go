package repository

import (
	"FocusList/internal/model"
	"database/sql"
	"fmt"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User

	query := `
	SELECT email, first_name, last_name, password, created_at, updated_at, is_active, role
	FROM users where email = $1
	`
	row := r.db.QueryRow(query, email)
	var createdAt, updatedAt time.Time
	if err := row.Scan(&user.Email, &user.FirstName, &user.LastName, &user.Password, &createdAt, &updatedAt, &user.IsActive, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}
	user.CreatedAt = createdAt.Format(time.RFC3339)
	user.UpdatedAt = updatedAt.Format(time.RFC3339)
	return &user, nil
}

func (r *UserRepository) CreateUser(user *model.User) error {
	query := `
	INSERT INTO users (email, first_name, last_name, password, created_at, updated_at, is_active, role)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	createdAt := time.Now().Format(time.RFC3339)
	updatedAt := createdAt
	_, err := r.db.Exec(query, user.Email, user.FirstName, user.LastName, user.Password, createdAt, updatedAt, user.IsActive, user.Role)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}
