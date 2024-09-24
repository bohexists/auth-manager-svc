package middleware

import (
	"net/http"
	"strings"

	"github.com/bohexists/auth-manager-svc/internal/auth"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(jwtService auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header not provided"})
			c.Abort()
			return
		}

		// Extract token from header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || token.Valid() != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
