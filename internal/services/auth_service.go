package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/abubakar-sadiq001/learning-go-shop/internal/config"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/dto"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/models"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db     *gorm.DB
	config *config.Config
}

func NewAuthService(db *gorm.DB, config *config.Config) *AuthService {
	return &AuthService{
		db:     db,
		config: config,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      "customer",
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Create a cart
	cart := models.Cart{UserID: user.ID}
	if err := s.db.Create(&cart).Error; err != nil {
		fmt.Println("Unable to create cart")
	}

	// Generate token
	return s.generateAuthResponse(&user)
}

func (s *AuthService) Login(dtoReq *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user
	var user models.User
	if err := s.db.Where("email = ? AND is_active = ?", dtoReq.Email, true).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if !utils.CheckPasswordHash(dtoReq.Password, user.Password) {
		return nil, errors.New("invalid credentilas")
	}

	// Generate token
	return s.generateAuthResponse(&user)
}

func (s *AuthService) RefreshToken(dtoReq *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	// Validate refresh token
	claims, err := utils.ValidateToken(dtoReq.RefreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Retrieve and check if token exist or expired
	var refreshToken models.RefreshToken
	if err := s.db.Where("token = ? AND expires > ?", dtoReq.RefreshToken, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	// Retrieve user
	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Delete token if user is not found
	s.db.Delete(&refreshToken)

	// Generate token
	return s.generateAuthResponse(&user)
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.db.Where("token = ?", refreshToken).Delete(&models.RefreshToken{}).Error
}

func (s *AuthService) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	// Generate token pairs
	accessToken, refreshToken, err := utils.GenerateTokenPair(
		&s.config.JWT,
		user.ID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.config.JWT.RefreshTokenExpiresIn),
	}

	// Insert token in to the DB
	s.db.Create(&refreshTokenModel)

	// Return auth response
	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
