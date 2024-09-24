package services

import (
	"github.com/bohexists/auth-manager-svc/config"
	"github.com/bohexists/auth-manager-svc/domain"
)

type JWTServiceInterface interface {
	GenerateToken(userID int64) (string, error)
	ValidateToken(token string) (*domain.TokenClaims, error)
	GenerateRefreshToken(userID int64) (string, error)
}

type JWTService struct {
	tokenService domain.TokenServiceInterface
}

func NewJWTService(cfg *config.Config) *JWTService {
	tokenService := domain.NewJWTService(cfg.JWTSecret, cfg.RefreshTokenSecret, "auth_manager")
	return &JWTService{
		tokenService: tokenService,
	}
}

// GenerateToken creates a new JWT token for a given user
func (j *JWTService) GenerateToken(userID int64) (string, error) {
	return j.tokenService.GenerateAccessToken(userID)
}

// GenerateRefreshToken creates a new refresh JWT token for a given user
func (j *JWTService) GenerateRefreshToken(userID int64) (string, error) {
	return j.tokenService.GenerateRefreshToken(userID)
}

// ValidateToken checks if the provided JWT token is valid
func (j *JWTService) ValidateToken(token string) (*domain.TokenClaims, error) {
	return j.tokenService.ValidateToken(token)
}
