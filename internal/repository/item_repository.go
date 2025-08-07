package repository

import (
	"FocusList/internal/model"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) CreateItem(item *model.TodoListItem, userEmail string) error {
	itemCreateQuery := `
		INSERT INTO todo_list_items (id, bucket_id, description, created_at, updated_at, is_complete)
		VALUES ($1, $2, $3, $4, $5, $6)
		`

	_, err := r.db.Exec(itemCreateQuery, item.ID, item.BucketID, item.Description, item.CreatedAt, item.UpdatedAt, false)
	return err
}

func (r *ItemRepository) GetItemsByBucketID(bucketID string) ([]*model.TodoListItem, error) {
	query := `
		SELECT id, bucket_id, description, created_at, updated_at, is_complete FROM todo_list_items
		WHERE bucket_id = $1`
	rows, err := r.db.Query(query, bucketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.TodoListItem
	for rows.Next() {
		var item model.TodoListItem
		if err := rows.Scan(&item.ID, &item.BucketID, &item.Description, &item.CreatedAt, &item.UpdatedAt, &item.IsCompleted); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

func (r *ItemRepository) GetItemByID(itemID string) (*model.TodoListItem, error) {
	query := `SELECT id, bucket_id, description, created_at, updated_at, is_complete FROM todo_list_items
		WHERE id = $1`
	row := r.db.QueryRow(query, itemID)
	var item model.TodoListItem
	if err := row.Scan(&item.ID, &item.BucketID, &item.Description, &item.CreatedAt, &item.UpdatedAt, &item.IsCompleted); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepository) UpdateItem(todoItem *model.TodoListItem) error {
	// Get the Item
	query := `SELECT id, bucket_id, description, created_at, updated_at, is_complete FROM todo_list_items
		WHERE id = $1`
	row := r.db.QueryRow(query, todoItem.ID)
	var existing model.TodoListItem
	if err := row.Scan(&existing.ID, &existing.BucketID, &existing.Description, &existing.CreatedAt, &existing.UpdatedAt, &existing.IsCompleted); err != nil {
		return nil
	}

	var updates []string
	var args []interface{}
	argIdx := 1

	if todoItem.Description != existing.Description {
		updates = append(updates, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, todoItem.Description)
		argIdx++
	}
	if todoItem.IsCompleted != existing.IsCompleted {
		updates = append(updates, fmt.Sprintf("is_complete = $%d", argIdx))
		args = append(args, todoItem.IsCompleted)
		argIdx++
	}

	// If no changes, skip update
	if len(updates) == 0 {
		fmt.Println("There are no updates to the item")
		return nil
	}

	// Always update the updated_at timestamp
	updates = append(updates, fmt.Sprintf("updated_at = $%d", argIdx))
	args = append(args, time.Now())
	argIdx++

	// Final update query
	updateQuery := fmt.Sprintf("UPDATE todo_list_items SET %s WHERE id = $%d", strings.Join(updates, ", "), argIdx)
	args = append(args, todoItem.ID)

	_, err := r.db.Exec(updateQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

func (r *ItemRepository) DeleteItemByID(itemID string) error {
	query := `DELETE FROM todo_list_items WHERE id = $1`
	_, err := r.db.Exec(query, itemID)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}
	return nil
}
