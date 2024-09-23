package routes

import (
	"github.com/bohexists/auth-manager-svc/internal/auth"
	"github.com/bohexists/auth-manager-svc/pkg/middleware"
	"github.com/bohexists/auth-manager-svc/transport/http/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *handlers.AuthHandler, jwtService auth.JWTService) *gin.Engine {
	router := gin.Default()

	// Auth routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	// Protected routes
	protected := router.Group("/protected")
	protected.Use(middleware.JWTAuthMiddleware(jwtService))
	protected.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "this is a protected route"})
	})

	return router
}
