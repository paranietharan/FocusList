package repository

import (
	"FocusList/internal/model"
	"database/sql"
	"fmt"
	"time"
)

type Repo struct {
	db *sql.DB
}

// Get the user
func (r *Repo) GetUserByEmail(email string) (*model.User, error) {
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

// Create a new user
func (r *Repo) CreateUser(user *model.User) error {
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

func (r *Repo) GetTodoListBucketsByUserEmail(email string) ([]model.TodoListBucket, error) {
	var buckets []model.TodoListBucket
	query := `
		SELECT b.id, b.name, b.created_at, b.updated_at
		FROM todo_list_buckets b
		JOIN todo_list_bucket_users u ON b.id = u.bucket_id
		WHERE u.user_email = $1
	`
	rows, err := r.db.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b model.TodoListBucket
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&b.ID, &b.Name, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		b.CreatedAt = createdAt.Format(time.RFC3339)
		b.UpdatedAt = updatedAt.Format(time.RFC3339)

		buckets = append(buckets, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return buckets, nil
}
