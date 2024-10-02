package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bohexists/auth-manager-svc/domain"
	"github.com/bohexists/auth-manager-svc/internal/services"
)

// AuthHandler is a handler for repositorys authentication
type AuthHandler struct {
	authService *services.AuthService
	JWTService  *services.JWTService
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authService *services.AuthService, JWTService *services.JWTService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		JWTService:  JWTService,
	}

}

// Register handles repositorys registration requests
func (h *AuthHandler) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	// Call the services layer to register the repositorys
	if err := h.authService.Register(&user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "repositorys created successfully"})
}

// Login with JWT token generation
func (h *AuthHandler) Login(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	// Call the services layer to login the repositorys
	user, err := h.authService.Login(request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Generate the JWT access and refresh tokens
	accessToken, err := h.JWTService.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := h.JWTService.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	// Set refresh token in a secure cookie
	c.SetCookie("refresh_token", refreshToken, 60*60*24*7, "/", "", false, true) // 7 days expiration, HTTPOnly

	// Return the JWT access token
	c.JSON(http.StatusOK, gin.H{"message": "login successful", "access_token": accessToken})
}

// RefreshToken updates the access token using the refresh token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Resive refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token found"})
		return
	}

	// Validate refresh token
	claims, err := h.JWTService.ValidateToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	// Generate new access token
	newAccessToken, err := h.JWTService.GenerateAccessToken(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate new access token"})
		return
	}

	// Return new access token
	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}
