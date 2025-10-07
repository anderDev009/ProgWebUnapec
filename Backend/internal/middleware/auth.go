package middleware

import (
	"net/http"
	"strings"

	"petmatch/internal/models"
	"petmatch/internal/services"

	"github.com/gin-gonic/gin"
)

const (
	userContextKey = "currentUser"
)

func Authentication(auth *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		user, err := auth.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set(userContextKey, user)
		c.Next()
	}
}

func RequireRoles(roles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get(userContextKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		user, ok := value.(*models.User)
		if !ok || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		for _, role := range roles {
			if user.Role == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}

func CurrentUser(c *gin.Context) *models.User {
	value, exists := c.Get(userContextKey)
	if !exists {
		return nil
	}
	user, ok := value.(*models.User)
	if !ok {
		return nil
	}
	return user
}
