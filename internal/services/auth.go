package services

import (
	"errors"

	"github.com/bohexists/auth-manager-svc/domain"
	"github.com/bohexists/auth-manager-svc/internal/repositorys"
)

// AuthServiceInterface defines the interface for AuthService
type AuthServiceInterface interface {
	Register(user *domain.User) error
	Login(email, password string) (*domain.User, error)
}

// UserRepositoryInterface defines the interface for UserRepository
type UserRepositoryInterface interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}

// AuthService is a service for repositorys authentication
type AuthService struct {
	userRepo     UserRepositoryInterface
	tokenService JWTServiceInterface
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepo repositorys.UserRepositoryInterface, tokenService JWTServiceInterface) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// Register handles repositorys registration logic
func (s *AuthService) Register(user *domain.User) error {
	// Check if a repositorys with the given email already exists
	_, err := s.userRepo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("repositorys already exists")
	}

	// Hash the repositorys's password before saving
	if err := user.HashPassword(); err != nil {
		return err
	}

	// Save the new repositorys to the database
	return s.userRepo.Create(user)
}

// Login handles repositorys authentication logic
func (s *AuthService) Login(email, password string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("repositorys not found")
	}

	// Verify the repositorys's password
	if err := user.CheckPassword(password); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
