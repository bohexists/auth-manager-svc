package routes

import (
	"github.com/bohexists/auth-manager-svc/transport/http/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the HTTP router
func SetupRouter(authHandler *handlers.AuthHandler, jwtMiddleware gin.HandlerFunc) *gin.Engine {
	router := gin.Default()

	// Auth routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	// Protected routes
	protected := router.Group("/protected")
	protected.Use(jwtMiddleware)
	protected.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "this is a protected route"})
	})

	return router
}
