// Package auth mirrors Security/AuthController.kt + Security/AuthService.kt.
package auth

import (
	"errors"

	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/apperrors"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/dto"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/models"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/repository"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/security"
)

// Service mirrors Security/AuthService.kt.
type Service struct {
	userRepo    *repository.UserRepository
	hashEncoder *security.HashEncoder
	jwtService  *security.JWTService
}

func NewService(userRepo *repository.UserRepository, hashEncoder *security.HashEncoder, jwtService *security.JWTService) *Service {
	return &Service{userRepo: userRepo, hashEncoder: hashEncoder, jwtService: jwtService}
}

// Register mirrors `fun register(request: RegisterRequest): AuthResponse`.
func (s *Service) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, &apperrors.UserAlreadyExistsError{Msg: "User with email " + req.Email + " already exists"}
	}

	hashed, err := s.hashEncoder.Encode(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashed,
	}

	saved, err := s.userRepo.Save(user)
	if err != nil {
		return nil, err
	}

	return s.buildAuthResponse(saved)
}

// Login mirrors `fun login(request: LoginRequest): AuthResponse`.
func (s *Service) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, &apperrors.InvalidCredentialsError{Msg: "Invalid email or password"}
	}

	if !s.hashEncoder.Matches(req.Password, user.PasswordHash) {
		return nil, &apperrors.InvalidCredentialsError{Msg: "Invalid email or password"}
	}

	return s.buildAuthResponse(user)
}

// Refresh mirrors `fun refresh(refreshToken: String): AuthResponse`.
func (s *Service) Refresh(refreshToken string) (*dto.AuthResponse, error) {
	if !s.jwtService.ValidateToken(refreshToken) || !s.jwtService.IsRefreshToken(refreshToken) {
		return nil, &apperrors.InvalidTokenError{Msg: "Invalid or expired refresh token"}
	}

	userID, err := s.jwtService.ExtractUserID(refreshToken)
	if err != nil {
		return nil, &apperrors.InvalidTokenError{Msg: "Invalid or expired refresh token"}
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, &apperrors.InvalidTokenError{Msg: "User no longer exists"}
	}

	return s.buildAuthResponse(user)
}

func (s *Service) buildAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	if user.ID.String() == "" {
		return nil, errors.New("user ID missing")
	}
	access, err := s.jwtService.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, err
	}
	refresh, err := s.jwtService.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		Name:         user.Name,
	}, nil
}
