package service

import (
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
