package middleware

import (
	"FocusList/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)

		if len(requiredRoles) > 0 {
			roleAllowed := false
			for _, allowedRole := range requiredRoles {
				if claims.Role == allowedRole {
					roleAllowed = true
					break
				}
			}
			if !roleAllowed {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient role"})
				return
			}
		}

		c.Next()
	}
}
