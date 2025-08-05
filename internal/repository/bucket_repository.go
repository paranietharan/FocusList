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

func (r *TodoListBucketRepository) GetBucketByID(bucketID, userEmail string) (*model.TodoListBucket, error) {
	if bucketID == "" {
		return nil, fmt.Errorf("bucket ID is required")
	}

	query := `
	SELECT b.id, b.name, b.created_at, b.updated_at
	FROM todo_list_buckets b
	JOIN todo_list_bucket_users u ON b.id = u.bucket_id
	WHERE b.id = $1 AND u.user_email = $2
	`
	row := r.db.QueryRow(query, bucketID, userEmail)

	var bucket model.TodoListBucket
	if err := row.Scan(&bucket.ID, &bucket.Name, &bucket.CreatedAt, &bucket.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bucket not found")
		}
		return nil, err
	}
	return &bucket, nil
}

func (r *TodoListBucketRepository) UpdateBucketName(bucketID, newName, userEmail string) error {
	if bucketID == "" || newName == "" {
		return fmt.Errorf("bucket ID and new name are required")
	}

	query := `
	UPDATE todo_list_buckets
	SET name = $1, updated_at = NOW()
	WHERE id = $2 AND EXISTS (
		SELECT 1 FROM todo_list_bucket_users
		WHERE bucket_id = $2 AND user_email = $3
	)
	`
	_, err := r.db.Exec(query, newName, bucketID, userEmail)
	if err != nil {
		fmt.Println("Error updating bucket name:", err)
		return err
	}
	return nil
}

func (r *TodoListBucketRepository) AddUserToBucket(bucketID, userEmail, email string) error {
	if bucketID == "" || userEmail == "" {
		return fmt.Errorf("bucket ID and user email are required")
	}

	query := `
	INSERT INTO todo_list_bucket_users (bucket_id, user_email, permission)
	VALUES ($1, $2, $3)
	ON CONFLICT (bucket_id, user_email) DO NOTHING
	`
	_, err := r.db.Exec(query, bucketID, userEmail, model.ExecutePermission)
	if err != nil {
		fmt.Println("Error adding user to bucket:", err)
		return err
	}
	return nil
}

func (r *TodoListBucketRepository) DeleteBucket(bucketID, userEmail string) error {
	if bucketID == "" {
		return fmt.Errorf("bucket ID is required")
	}

	query := `
	DELETE FROM todo_list_buckets
	WHERE id = $1 AND EXISTS (
		SELECT 1 FROM todo_list_bucket_users
		WHERE bucket_id = $1 AND user_email = $2
	)
	`
	_, err := r.db.Exec(query, bucketID, userEmail)
	if err != nil {
		fmt.Println("Error deleting bucket:", err)
		return err
	}
	return nil
}
