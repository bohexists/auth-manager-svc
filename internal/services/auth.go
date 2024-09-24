package services

import (
	"errors"

	"github.com/bohexists/auth-manager-svc/domain"
	"github.com/bohexists/auth-manager-svc/internal/user"
)

// AuthServiceInterface определяет интерфейс для AuthService
type AuthServiceInterface interface {
	Register(user *domain.User) error
	Login(email, password string) (*domain.User, error)
}

type AuthService struct {
	userRepo     user.UserRepositoryInterface
	tokenService JWTServiceInterface
}

func NewAuthService(userRepo user.UserRepositoryInterface, tokenService JWTServiceInterface) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// Register handles user registration logic
func (s *AuthService) Register(user *domain.User) error {
	// Check if a user with the given email already exists
	_, err := s.userRepo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("user already exists")
	}

	// Hash the user's password before saving
	if err := user.HashPassword(); err != nil {
		return err
	}

	// Save the new user to the database
	return s.userRepo.Create(user)
}

// Login handles user authentication logic
func (s *AuthService) Login(email, password string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Verify the user's password
	if err := user.CheckPassword(password); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
