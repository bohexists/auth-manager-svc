package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TokenClaims represents the structure of claims inside the JWT token
type TokenClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

// TokenService interface for generating and validating tokens
type TokenServiceInterface interface {
	GenerateAccessToken(userID int64) (string, error)
	GenerateRefreshToken(userID int64) (string, error)
	ValidateToken(token string) (*TokenClaims, error)
}

// JWTService implements TokenService
type JWTService struct {
	secretKey        string
	refreshSecretKey string
	issuer           string
}

func (s *JWTService) GenerateToken(userID int64) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewJWTService(secretKey, refreshSecretKey, issuer string) *JWTService {
	return &JWTService{
		secretKey:        secretKey,
		refreshSecretKey: refreshSecretKey,
		issuer:           issuer,
	}
}

// GenerateAccessToken creates a new access token for a given user
func (s *JWTService) GenerateAccessToken(userID int64) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(), // Access token is valid for 15 minutes
			Issuer:    s.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// GenerateRefreshToken creates a new refresh token for a given user
func (s *JWTService) GenerateRefreshToken(userID int64) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(), // Refresh token valid for 7 days
			Issuer:    s.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.refreshSecretKey))
}

// ValidateToken checks if the provided JWT token is valid
func (s *JWTService) ValidateToken(token string) (*TokenClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate the token and its claims
	if claims, ok := parsedToken.Claims.(*TokenClaims); ok && parsedToken.Valid {
		return claims, nil
	} else {
		return nil, jwt.ErrSignatureInvalid
	}
}
