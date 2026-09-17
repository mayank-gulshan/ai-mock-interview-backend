// Package models mirrors the Kotlin @Entity classes in Entity/*.kt.
//
// NOTE on UserPrincipal (Entity/UserPrinciple.kt): that class exists purely
// to satisfy Spring Security's UserDetails interface (getAuthorities,
// isAccountNonExpired, etc.). Go/Gin has no equivalent interface — the
// authenticated *User is stored directly in the gin.Context by the auth
// middleware (see internal/security/middleware.go), so UserPrincipal has no
// standalone Go equivalent and is intentionally not translated 1:1.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User mirrors Entity/User.kt.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name         string    `gorm:"not null;default:''"`
	Email        string    `gorm:"unique;not null;default:''"`
	PasswordHash string    `gorm:"column:password_hash;not null;default:''"`
	CreatedAt    time.Time `gorm:"not null"`
}

func (User) TableName() string { return "users" }

// BeforeCreate mimics @GeneratedValue(strategy = GenerationType.UUID) and the
// `createdAt: Instant = Instant.now()` default on the Kotlin data class.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}
	return nil
}
