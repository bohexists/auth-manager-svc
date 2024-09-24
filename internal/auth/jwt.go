package auth

import (
	"github.com/bohexists/auth-manager-svc/config"
	"github.com/bohexists/auth-manager-svc/domain"
)

type JWTService interface {
	GenerateToken(userID int64) (string, error)
	ValidateToken(token string) (*domain.TokenClaims, error)
	GenerateRefreshToken(userID int64) (string, error)
}

type jwtService struct {
	tokenService domain.TokenService
}

func NewJWTService(cfg *config.Config) JWTService {
	tokenService := domain.NewJWTService(cfg.JWTSecret, cfg.RefreshTokenSecret, "auth_manager")
	return &jwtService{
		tokenService: tokenService,
	}
}

// GenerateToken creates a new JWT token for a given user
func (j *jwtService) GenerateToken(userID int64) (string, error) {
	return j.tokenService.GenerateAccessToken(userID)
}

// GenerateRefreshToken creates a new refresh JWT token for a given user
func (j *jwtService) GenerateRefreshToken(userID int64) (string, error) {
	return j.tokenService.GenerateRefreshToken(userID)
}

// ValidateToken checks if the provided JWT token is valid
func (j *jwtService) ValidateToken(token string) (*domain.TokenClaims, error) {
	return j.tokenService.ValidateToken(token)
}
