package handler

import "github.com/gin-gonic/gin"

type TestHandler struct {
}

func (t *TestHandler) UserDetailsCheck(c *gin.Context) {
	email := c.GetString("userEmail")
	c.JSON(200, gin.H{
		"email":   email,
		"message": "Hello, user!",
		"role":    c.GetString("userRole"),
	})
}
