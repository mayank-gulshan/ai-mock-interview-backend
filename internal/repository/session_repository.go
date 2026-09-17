package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/models"
)

// SessionRepository mirrors Repository/SessionRepository.kt. NOTE: not
// currently called by any service/handler in the original backend (see the
// NOTE in internal/models/session.go); kept for 1:1 source parity.
type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// FindByUserOrderByCreatedAtDesc mirrors
// `fun findByUserOrderByCreatedAtDesc(user: User): List<Session>`.
func (r *SessionRepository) FindByUserOrderByCreatedAtDesc(userID uuid.UUID) ([]models.Session, error) {
	var sessions []models.Session
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&sessions).Error
	return sessions, err
}

func (r *SessionRepository) Save(session *models.Session) (*models.Session, error) {
	if err := r.db.Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (r *SessionRepository) FindByID(id uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := r.db.Preload("Answers").Where("id = ?", id).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}
