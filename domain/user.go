package domain

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
	"unicode"
)

// User represents the user entity in the system
type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email" binding:"required,email"`
	Password  string    `gorm:"not null" json:"password" binding:"required,min=8"`
	IsActive  bool      `gorm:"default:false" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Validate performs validation checks on the User fields
func (u *User) Validate() error {
	// Check if email is valid
	if !isValidEmail(u.Email) {
		return errors.New("invalid email format")
	}

	// Check password complexity
	if err := validatePassword(u.Password); err != nil {
		return err
	}

	return nil
}

// HashPassword hashes the user's password before saving to the database
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword compares the provided password with the hashed password
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// Helper function to validate email format using a regex pattern
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(email)
}

// validatePassword ensures the password meets complexity requirements
func validatePassword(password string) error {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case isSpecialCharacter(char):
			hasSpecial = true
		}
	}

	if !hasMinLen || !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("password must be at least 8 characters long and include at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	return nil
}

// Helper function to check for special characters
func isSpecialCharacter(char rune) bool {
	specialChars := "!@#$%^&*()_+[]{}|;:,.<>?/\\"
	return unicode.In(char, unicode.Punct, unicode.Symbol) || containsRune(specialChars, char)
}

// Helper function to check if a rune is in a string
func containsRune(s string, r rune) bool {
	for _, v := range s {
		if v == r {
			return true
		}
	}
	return false
}
