package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/models"
)

// AnswerRepository mirrors Repository/AnswerRepository.kt. NOTE: not
// currently called by any service/handler in the original backend; kept for
// 1:1 source parity.
type AnswerRepository struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) *AnswerRepository {
	return &AnswerRepository{db: db}
}

// FindBySession mirrors `fun findBySession(session: Session): List<Answer>`.
func (r *AnswerRepository) FindBySession(sessionID uuid.UUID) ([]models.Answer, error) {
	var answers []models.Answer
	err := r.db.Where("session_id = ?", sessionID).Find(&answers).Error
	return answers, err
}

func (r *AnswerRepository) Save(answer *models.Answer) (*models.Answer, error) {
	if err := r.db.Create(answer).Error; err != nil {
		return nil, err
	}
	return answer, nil
}
