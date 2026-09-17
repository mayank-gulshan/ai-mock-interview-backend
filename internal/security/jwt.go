package security

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Token validity, copied verbatim from Security/JWTService.kt:
//
//	private val accessTokenValidityMs = 15L * 60L * 1000L               // 15 minutes
//	private val refreshTokenValidityMs = 30L * 24 * 60 * 60L * 1000L    // 30 days
const (
	accessTokenValidity  = 15 * time.Minute
	refreshTokenValidity = 30 * 24 * time.Hour
)

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

// JWTService mirrors Security/JWTService.kt.
type JWTService struct {
	secretKey []byte
}

// NewJWTService takes the same base64-encoded secret as
// `@Value("\${app.jwt.secret}") private val secret: String`, decoded exactly
// like `Keys.hmacShaKeyFor(Base64.getDecoder().decode(secret))`.
func NewJWTService(base64Secret string) (*JWTService, error) {
	decoded, err := base64.StdEncoding.DecodeString(base64Secret)
	if err != nil {
		return nil, err
	}
	return &JWTService{secretKey: decoded}, nil
}

type claims struct {
	Type string `json:"type"`
	jwt.RegisteredClaims
}

// generateToken mirrors the private `generateToken(userId, type, expiry)`.
func (s *JWTService) generateToken(userID, tokenType string, validity time.Duration) (string, error) {
	now := time.Now()
	c := claims{
		Type: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(validity)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(s.secretKey)
}

// GenerateAccessToken mirrors `fun generateAccessToken(userId: String): String`.
func (s *JWTService) GenerateAccessToken(userID string) (string, error) {
	return s.generateToken(userID, tokenTypeAccess, accessTokenValidity)
}

// GenerateRefreshToken mirrors `fun generateRefreshToken(userId: String): String`.
func (s *JWTService) GenerateRefreshToken(userID string) (string, error) {
	return s.generateToken(userID, tokenTypeRefresh, refreshTokenValidity)
}

func (s *JWTService) parseClaims(token string) (*claims, error) {
	c := &claims{}
	parsed, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})
	if err != nil || !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	return c, nil
}

// ValidateToken mirrors `fun validateToken(token: String): Boolean`.
func (s *JWTService) ValidateToken(token string) bool {
	_, err := s.parseClaims(token)
	return err == nil
}

// ExtractUserID mirrors `fun extractUserId(token: String): UUID`.
func (s *JWTService) ExtractUserID(token string) (uuid.UUID, error) {
	c, err := s.parseClaims(token)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(c.Subject)
}

// ExtractTokenType mirrors `fun extractTokenType(token: String): String`.
func (s *JWTService) ExtractTokenType(token string) (string, error) {
	c, err := s.parseClaims(token)
	if err != nil {
		return "", err
	}
	return c.Type, nil
}

// IsAccessToken mirrors `fun isAccessToken(token: String): Boolean`.
func (s *JWTService) IsAccessToken(token string) bool {
	t, err := s.ExtractTokenType(token)
	return err == nil && t == tokenTypeAccess
}

// IsRefreshToken mirrors `fun isRefreshToken(token: String): Boolean`.
func (s *JWTService) IsRefreshToken(token string) bool {
	t, err := s.ExtractTokenType(token)
	return err == nil && t == tokenTypeRefresh
}
