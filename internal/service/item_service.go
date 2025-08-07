package service

import (
	"FocusList/internal/model"
	"FocusList/internal/repository"
	"time"
)

type ItemService struct {
	ItemRepo *repository.ItemRepository
}

func NewItemService(itemRepo *repository.ItemRepository) *ItemService {
	return &ItemService{ItemRepo: itemRepo}
}

func (s *ItemService) CreateItem(item *model.TodoListItem, userEmail string) error {
	now := time.Now().Format(time.RFC3339)
	item.CreatedAt = now
	item.UpdatedAt = now
	item.IsCompleted = false
	return s.ItemRepo.CreateItem(item, userEmail)
}

func (s *ItemService) GetItemsByBucketID(bucketID string) ([]*model.TodoListItem, error) {
	return s.ItemRepo.GetItemsByBucketID(bucketID)
}

func (s *ItemService) GetItemByID(itemID string) (*model.TodoListItem, error) {
	return s.ItemRepo.GetItemByID(itemID)
}

func (s *ItemService) UpdateItem(item *model.TodoListItem) error {
	item.UpdatedAt = time.Now().Format(time.RFC3339)
	return s.ItemRepo.UpdateItem(item)
}

func (s *ItemService) DeleteItemByID(itemID string) error {
	return s.ItemRepo.DeleteItemByID(itemID)
}
