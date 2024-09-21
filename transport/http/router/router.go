package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bohexists/auth-manager-svc/transport/http/handlers"
)

func SetupRouter(authHandler *handlers.AuthHandler) *gin.Engine {
	router := gin.Default()

	// Auth routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	return router
}
