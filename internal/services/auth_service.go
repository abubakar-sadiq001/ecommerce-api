package services

import (
	"errors"
	"fmt"
	"time"

	appconfig "github.com/abubakar-sadiq001/learning-go-shop/internal/config"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/dto"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/events"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/models"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/notifications"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db             *gorm.DB
	config         *appconfig.Config
	eventPublisher events.Publisher
}

func NewAuthService(db *gorm.DB, config *appconfig.Config, eventPublisher events.Publisher) *AuthService {
	return &AuthService{
		db:             db,
		config:         config,
		eventPublisher: eventPublisher,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// 1. Check if user exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("you cannot register with this email")
	}
	// if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
	// 	return nil, errors.New("you cannot register with this email")
	// }

	// 2. Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 3. Create user
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

	// 4. Create a cart
	cart := models.Cart{UserID: user.ID}
	if err := s.db.Create(&cart).Error; err != nil {
		fmt.Println("Unable to create cart")
	}

	// 5. Return generated token
	return s.generateAuthResponse(&user)
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// 1. Find user
	var user models.User
	if err := s.db.Where("email = ? AND is_active = ?", req.Email, true).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 2. Verify password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentilas")
	}

	// 3. Return generated token
	return s.generateAuthResponse(&user)
}

func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	// 1. Validate refresh token
	claims, err := utils.ValidateToken(req.RefreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 2. Retrieve and check if token exist or expired
	var refreshToken models.RefreshToken
	if err := s.db.Where("token = ? AND expires_at > ?", req.RefreshToken, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	// 3. Retrieve user
	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// 4. Delete token if user is not found
	s.db.Delete(&refreshToken)

	// 5. Return generated token
	return s.generateAuthResponse(&user)
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.db.Where("token = ?", refreshToken).Delete(&models.RefreshToken{}).Error
}

func (s *AuthService) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	// 1. Generate token pairs
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

	// 2. Delete expired tokens or ALL previous tokens for this user first
	s.db.Where("user_id = ?", user.ID).Delete(&models.RefreshToken{})

	// 3. Insert the new token
	if err := s.db.Create(&refreshTokenModel).Error; err != nil {
		return nil, err
	}

	// if err := s.db.Model(&models.RefreshToken{}).Where("user_id = ?", user.ID).Update("token", refreshToken).Error; err != nil {
	// 	return nil, err
	// }

	// 4. Publish login in event
	err = s.eventPublisher.Publish(notifications.UserLoggedIn, user, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("unable to publish user login event: %w", err)
	}

	// 5. Return auth response
	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
