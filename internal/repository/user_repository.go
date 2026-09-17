package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/models"
)

// UserRepository mirrors Repository/UserRepository.kt (a Spring Data
// JpaRepository<User, UUID> with two derived query methods).
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByEmail mirrors `fun findByEmail(email: String): User?`.
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ExistsByEmail mirrors `fun existsByEmail(email: String): Boolean`.
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindByID mirrors JpaRepository's inherited `findById(id): Optional<User>`.
func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Save mirrors JpaRepository's inherited `save(user): User`.
func (r *UserRepository) Save(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
