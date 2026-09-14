package router

import (
	"net/http"
	"strings"

	"restapirian/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID   = "auth.user_id"
	ContextUsername = "auth.username"
	ContextRole     = "auth.role"
)

func JWTAuth(tokens service.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		parts := strings.Fields(strings.TrimSpace(c.GetHeader("Authorization")))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access denied"})
			return
		}
		claims, err := tokens.Parse(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired access token"})
			return
		}
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUsername, claims.Username)
		c.Set(ContextRole, claims.Role)
		c.Next()
	}
}

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}
	return func(c *gin.Context) {
		roleValue, ok := c.Get(ContextRole)
		role, valid := roleValue.(string)
		if !ok || !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}
		if _, permitted := allowed[role]; !permitted {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient role permission"})
			return
		}
		c.Next()
	}
}

func CurrentUserID(c *gin.Context) int64 {
	value, exists := c.Get(ContextUserID)
	if !exists {
		return 0
	}
	userID, _ := value.(int64)
	return userID
}
