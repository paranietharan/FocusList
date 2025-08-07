package middleware

import (
	"FocusList/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var RoleRank = map[string]int{
	"user":        1,
	"moderator":   2,
	"admin":       3,
	"super_admin": 4,
}

// Commented out the old AuthMiddleware for reference
// this version has many roles and if the required role is not in the list, it will return forbidden
// func AuthMiddleware(requiredRoles ...string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 			return
// 		}
// 		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
// 		claims, err := utils.ParseJWT(tokenStr)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
// 			return
// 		}

// 		c.Set("userEmail", claims.Email)
// 		c.Set("userRole", claims.Role)

// 		if len(requiredRoles) > 0 {
// 			roleAllowed := false
// 			for _, allowedRole := range requiredRoles {
// 				if claims.Role == allowedRole {
// 					roleAllowed = true
// 					break
// 				}
// 			}
// 			if !roleAllowed {
// 				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient role"})
// 				return
// 			}
// 		}

// 		c.Next()
// 	}
// }

// This auth middle ware works according to the role hierarchy
func AuthMiddleware(minRequiredRole string) gin.HandlerFunc {
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

		// Role hierarchy check
		userRank := RoleRank[claims.Role]
		requiredRank := RoleRank[minRequiredRole]

		if userRank < requiredRank {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient role"})
			return
		}

		c.Next()
	}
}
