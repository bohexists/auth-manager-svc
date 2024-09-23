package routes

import (
	"github.com/bohexists/auth-manager-svc/internal/auth"
	"github.com/bohexists/auth-manager-svc/pkg/middleware"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bohexists/auth-manager-svc/transport/http/handlers"
)

func SetupRouter(authHandler *handlers.AuthHandler, jwtService auth.JWTService) *gin.Engine {
	router := gin.Default()

	// Auth routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	// Protected routes
	protected := router.Group("/").Use(middleware.JWTAuthMiddleware(jwtService))
	{
		// Example protected route
		protected.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "this is a protected route"})
		})
	}

	return router
}
