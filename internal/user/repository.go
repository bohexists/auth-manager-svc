package user

import (
	"gorm.io/gorm"

	"github.com/bohexists/auth-manager-svc/domain"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create stores a new user in the database
func (r *UserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// FindByEmail searches for a user by their email
func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
