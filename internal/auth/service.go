package auth

import (
	"errors"

	"github.com/bohexists/auth-manager-svc/domain"
	"github.com/bohexists/auth-manager-svc/internal/user"
)

type AuthService struct {
	userRepo *user.UserRepository
}

func NewAuthService(userRepo *user.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
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
