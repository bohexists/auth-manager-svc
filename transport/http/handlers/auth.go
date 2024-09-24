package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bohexists/auth-manager-svc/domain"
	"github.com/bohexists/auth-manager-svc/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles user registration requests
func (h *AuthHandler) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	// Call the domain validation function
	if err := user.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service layer to register the user
	if err := h.authService.Register(&user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
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

	// Call the service layer to login the user
	user, err := h.authService.Login(request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Generate the JWT access and refresh tokens
	accessToken, err := h.authService.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := h.authService.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	// Set refresh token in a secure cookie
	c.SetCookie("refresh_token", refreshToken, 60*60*24*7, "/", "", false, true) // 7 days expiration, HTTPOnly

	// Return the JWT access token
	c.JSON(http.StatusOK, gin.H{"message": "login successful", "access_token": accessToken})
}

// RefreshToken обновляет access токен
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Получаем refresh токен из cookies
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token found"})
		return
	}

	// Валидация рефреш токена
	claims, err := h.authService.ValidateToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	// Генерация нового access токена
	newAccessToken, err := h.authService.GenerateToken(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate new access token"})
		return
	}

	// Возвращаем новый access token
	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}
