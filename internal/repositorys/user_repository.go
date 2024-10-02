package repositorys

import (
	"gorm.io/gorm"

	"github.com/bohexists/auth-manager-svc/domain"
)

// UserRepository is a repository for users
type UserRepositoryInterface interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{db: db}
}

// Create stores a new repositorys in the database
func (r *UserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// FindByEmail searches for a repositorys by their email
func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
