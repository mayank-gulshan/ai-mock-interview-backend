package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Answer mirrors Entity/Answer.kt. See the NOTE in session.go — currently
// unused by any controller/service, translated for parity.
type Answer struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	SessionID     uuid.UUID `gorm:"column:session_id;type:uuid;not null"`
	Session       Session   `gorm:"foreignKey:SessionID;references:ID"`
	QuestionText  string    `gorm:"type:text;not null;default:''"`
	UserAnswer    string    `gorm:"type:text;not null;default:''"`
	Rating        float64   `gorm:"not null;default:0"`
	Feedback      string    `gorm:"type:text;not null;default:''"`
	IdealAnswer   string    `gorm:"type:text;not null;default:''"`
	AccuracyScore int       `gorm:"not null;default:0"`
	DepthScore    int       `gorm:"not null;default:0"`
	ClarityScore  int       `gorm:"not null;default:0"`
	CreatedAt     time.Time `gorm:"not null"`
}

func (Answer) TableName() string { return "answers" }

func (a *Answer) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now().UTC()
	}
	return nil
}
