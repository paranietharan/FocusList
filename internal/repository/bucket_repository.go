package repository

import (
	"FocusList/internal/model"
	"database/sql"
	"fmt"
)

type TodoListBucketRepository struct {
	db *sql.DB
}

func NewTodoListBucketRepository(db *sql.DB) *TodoListBucketRepository {
	return &TodoListBucketRepository{db: db}
}

func (r *TodoListBucketRepository) CreateBucket(bucket *model.TodoListBucket, userEmail string) error {
	fmt.Println("Creating bucket with ID:", bucket.ID, "and name:", bucket.Name)
	bucketCreatQuery := `
	INSERT INTO todo_list_buckets (id, name, created_at, updated_at)
	VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(bucketCreatQuery, bucket.ID, bucket.Name, bucket.CreatedAt, bucket.UpdatedAt)
	if err != nil {
		fmt.Println("Error creating bucket:", err)
		return err
	}

	bucketUserQuery := `
	INSERT INTO todo_list_bucket_users (bucket_id, user_email, permission)
	VALUES ($1, $2, $3)
	`
	_, err = r.db.Exec(bucketUserQuery, bucket.ID, userEmail, model.ExecutePermission)
	if err != nil {
		fmt.Println("Error adding user to bucket:", err)
		return err
	}

	return nil
}

func (r *TodoListBucketRepository) GetBucketsByUserEmail(email string) ([]*model.TodoListBucket, error) {
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

	var buckets []*model.TodoListBucket
	for rows.Next() {
		var bucket model.TodoListBucket
		if err := rows.Scan(&bucket.ID, &bucket.Name, &bucket.CreatedAt, &bucket.UpdatedAt); err != nil {
			return nil, err
		}
		buckets = append(buckets, &bucket)
	}
	return buckets, nil
}
