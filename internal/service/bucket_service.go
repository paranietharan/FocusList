package service

import (
	"FocusList/internal/dto"
	"FocusList/internal/model"
	"FocusList/internal/repository"
	"FocusList/internal/utils"
	"fmt"
	"log"
	"time"
)

type TodoListBucketService struct {
	TodoListBucketRepo *repository.TodoListBucketRepository
}

func (s *TodoListBucketService) CreateBucket(bucketName, userEmail string) error {
	var bucket model.TodoListBucket

	bucket.ID = utils.GenerateUniqueUUID()
	bucket.Name = bucketName
	bucket.CreatedAt = time.Now().Format(time.RFC3339)
	bucket.UpdatedAt = bucket.CreatedAt

	fmt.Println("Creating bucket with ID:", bucket.ID, "and name:", bucket.Name)

	err := s.TodoListBucketRepo.CreateBucket(&bucket, userEmail)
	if err != nil {
		log.Println("Error creating bucket:", err)
	}
	return err
}

func (s *TodoListBucketService) GetBucketsByUserEmail(email string) ([]*model.TodoListBucket, error) {
	buckets, err := s.TodoListBucketRepo.GetBucketsByUserEmail(email)
	if err != nil {
		return nil, err
	}
	return buckets, nil
}

func (s *TodoListBucketService) GetBucketByID(bucketID, userEmail string) (*model.TodoListBucket, error) {
	if bucketID == "" {
		return nil, fmt.Errorf("bucket ID is required")
	}

	bucket, err := s.TodoListBucketRepo.GetBucketByID(bucketID, userEmail)
	if err != nil {
		log.Println("Error retrieving bucket by ID:", err)
		return nil, err
	}
	return bucket, nil
}

func (s *TodoListBucketService) UpdateBucketName(bucketID, newName, userEmail string) error {
	if bucketID == "" || newName == "" {
		return fmt.Errorf("bucket ID and new name are required")
	}

	err := s.TodoListBucketRepo.UpdateBucketName(bucketID, newName, userEmail)
	if err != nil {
		log.Println("Error updating bucket name:", err)
		return err
	}
	return nil
}

func (s *TodoListBucketService) AddUserToBucket(bucketID, userEmail, email string) error {
	if bucketID == "" || userEmail == "" {
		return fmt.Errorf("bucket ID and user email are required")
	}

	err := s.TodoListBucketRepo.AddUserToBucket(bucketID, userEmail, email)
	if err != nil {
		log.Println("Error adding user to bucket:", err)
		return err
	}
	return nil
}

func (s *TodoListBucketService) DeleteBucket(bucketID, userEmail string) error {
	if bucketID == "" {
		return fmt.Errorf("bucket ID is required")
	}

	err := s.TodoListBucketRepo.DeleteBucket(bucketID, userEmail)
	if err != nil {
		log.Println("Error deleting bucket:", err)
		return err
	}
	return nil
}

func (s *TodoListBucketService) GetBucketUsers(bucketID string) ([]*dto.TodoListBucketUserDTO, error) {
	if bucketID == "" {
		return nil, fmt.Errorf("bucket ID is required")
	}

	users, err := s.TodoListBucketRepo.GetBucketUsersByBucketID(bucketID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *TodoListBucketService) RemoveUserFromBucket(bucketID, userEmail, email string) error {
	if bucketID == "" || userEmail == "" || email == "" {
		fmt.Errorf("bucket ID and user email are required")
		return fmt.Errorf("bucket ID and user email are required")
	}

	// check if the user have permission to remove the user
	users, err := s.TodoListBucketRepo.GetBucketUsersByBucketID(bucketID)
	if err != nil {
		return err
	}

	deletePermission := false
	for _, user := range users {
		if user.UserEmail == userEmail {
			deletePermission = true
			break
		}
	}
	if !deletePermission {
		fmt.Errorf("user %s is not in bucket %s", userEmail, bucketID)
		return fmt.Errorf("user %s is not in bucket %s", userEmail, bucketID)
	}

	err = s.TodoListBucketRepo.RemoveUserFromBucket(bucketID, email)

	return err
}
