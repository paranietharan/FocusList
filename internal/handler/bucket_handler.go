package handler

import (
	"FocusList/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

type BucketHandler struct {
	TodoListBucketService *service.TodoListBucketService
}

func (h *BucketHandler) CreateBucket(c *gin.Context) {
	var input struct {
		Name string `json:"bucket_name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// TODO: DEBUG: Print the input for debugging purposes
	fmt.Printf("Received input: %+v\n", input)

	if input.Name == "" {
		c.JSON(400, gin.H{"error": "Bucket name is required"})
		return
	}

	email := c.GetString("userEmail")
	err := h.TodoListBucketService.CreateBucket(input.Name, email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create bucket"})
		return
	}
	c.JSON(201, gin.H{"message": "Bucket created successfully"})
}

func (h *BucketHandler) GetBuckets(c *gin.Context) {
	// Get the user email from the add from the token
	email := c.GetString("userEmail")
	buckets, err := h.TodoListBucketService.GetBucketsByUserEmail(email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve buckets"})
		return
	}
	c.JSON(200, gin.H{"buckets": buckets})
}
