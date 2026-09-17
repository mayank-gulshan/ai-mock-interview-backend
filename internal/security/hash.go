// Package security mirrors everything under Security/*.kt except
// AuthController/AuthService (moved to internal/auth to group the feature
// together) and the exception types (moved to internal/apperrors).
package security

import "golang.org/x/crypto/bcrypt"

// HashEncoder mirrors Security/HashEncoder.kt, which wraps Spring Security's
// BCryptPasswordEncoder with its default cost factor (10).
type HashEncoder struct{}

func NewHashEncoder() *HashEncoder {
	return &HashEncoder{}
}

// Encode mirrors `fun encode(raw: String): String`.
func (h *HashEncoder) Encode(raw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Matches mirrors `fun matches(raw: String, hashed: String): Boolean`.
func (h *HashEncoder) Matches(raw, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw)) == nil
}
