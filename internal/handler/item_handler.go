package handler

import (
	"FocusList/internal/model"
	"FocusList/internal/service"
	"FocusList/internal/utils"
	"github.com/gin-gonic/gin"
)

type ItemHandler struct {
	ItemService *service.ItemService
}

func (h *ItemHandler) CreateItem(c *gin.Context) {
	var input struct {
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if input.Description == "" {
		c.JSON(400, gin.H{"error": "Description is required"})
		return
	}

	bucketID := c.Param("bucketID")
	if bucketID == "" {
		c.JSON(400, gin.H{"error": "Bucket ID is required"})
		return
	}

	email := c.GetString("userEmail")
	item := &model.TodoListItem{
		ID:          utils.GenerateUniqueUUID(),
		BucketID:    bucketID,
		Description: input.Description,
	}
	err := h.ItemService.CreateItem(item, email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create item"})
		return
	}
	c.JSON(201, gin.H{"message": "Item created successfully"})
}

func (h *ItemHandler) GetItems(c *gin.Context) {
	bucketID := c.Param("bucketID")
	if bucketID == "" {
		c.JSON(400, gin.H{"error": "Bucket ID is required"})
		return
	}

	items, err := h.ItemService.GetItemsByBucketID(bucketID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve items"})
		return
	}
	c.JSON(200, gin.H{"items": items})
}

func (h *ItemHandler) GetItemByID(c *gin.Context) {
	itemID := c.Param("itemID")
	if itemID == "" {
		c.JSON(400, gin.H{"error": "Item ID is required"})
		return
	}

	item, err := h.ItemService.GetItemByID(itemID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve item"})
		return
	}
	c.JSON(200, gin.H{"item": item})
}

func (h *ItemHandler) UpdateItem(c *gin.Context) {
	itemID := c.Param("itemID")
	if itemID == "" {
		c.JSON(400, gin.H{"error": "Item ID is required"})
		return
	}

	var input struct {
		Description string `json:"description"`
		IsCompleted bool   `json:"is_completed"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if input.Description == "" {
		c.JSON(400, gin.H{"error": "Description is required"})
		return
	}

	item := &model.TodoListItem{
		ID:          itemID,
		Description: input.Description,
		IsCompleted: input.IsCompleted,
	}
	err := h.ItemService.UpdateItem(item)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update item"})
		return
	}
	c.JSON(200, gin.H{"message": "Item updated successfully"})
}

func (h *ItemHandler) DeleteItem(c *gin.Context) {
	itemID := c.Param("itemID")
	if itemID == "" {
		c.JSON(400, gin.H{"error": "Item ID is required"})
		return
	}

	err := h.ItemService.DeleteItemByID(itemID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete item"})
		return
	}
	c.JSON(200, gin.H{"message": "Item deleted successfully"})
}
