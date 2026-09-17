package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Session mirrors Entity/Session.kt.
//
// NOTE: as of the current backend, no controller actually creates or reads
// Sessions/Answers yet (InterviewController only calls GeminiService and
// AnswerEvaluationService directly, stateless). The Dto/InterviewDto.kt
// (StartSessionRequest/Response, SubmitAnswerRequest) and Dto/HistoryDto.kt
// types, plus SessionRepository/AnswerRepository, appear to be groundwork
// for a not-yet-wired-up "save session history" feature. We still translate
// the entities, repositories, and DTOs faithfully (per the "don't skip any
// file" requirement) so the Go codebase has full parity, even though they
// are currently unused/dead code exactly as in the original.
type Session struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	User           User      `gorm:"foreignKey:UserID;references:ID"`
	Role           string    `gorm:"not null;default:''"`
	Difficulty     string    `gorm:"not null;default:''"`
	TotalQuestions int       `gorm:"not null;default:5"`
	AverageScore   float64   `gorm:"not null;default:0"`
	Completed      bool      `gorm:"not null;default:false"`
	CreatedAt      time.Time `gorm:"not null"`

	// mappedBy = "session", cascade = [CascadeType.ALL] -> delete children
	// when the parent session is deleted.
	Answers []Answer `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE"`
}

func (Session) TableName() string { return "sessions" }

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now().UTC()
	}
	return nil
}
