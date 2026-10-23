package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	appconfig "github.com/abubakar-sadiq001/learning-go-shop/internal/config"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/dto"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/events"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/models"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/notifications"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/repositories"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/utils"
)

var _ AuthServiceInterface = (*AuthService)(nil) // makes the compiler prove that *AuthService satisfies AuthServiceInterface.

type AuthService struct {
	userRepo       repositories.UserRepositoryInterface
	cartRepo       repositories.CartRepositoryInterface
	config         *appconfig.Config
	eventPublisher events.Publisher
}

func NewAuthService(config *appconfig.Config, eventPublisher events.Publisher, userRepo repositories.UserRepositoryInterface, cartRepo repositories.CartRepositoryInterface) *AuthService {
	return &AuthService{
		config:         config,
		eventPublisher: eventPublisher,
		userRepo:       userRepo,
		cartRepo:       cartRepo,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// 1. Check if user exists
	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("you cannot register with this email")
	}

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
	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	// 4. Create a cart
	cart := models.Cart{UserID: user.ID}
	if err := s.cartRepo.Create(&cart); err != nil {
		fmt.Println("Unable to create cart")
		_ = err
	}

	// 5. Return generated token
	return s.generateAuthResponse(&user)
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// 1. Find user
	user, err := s.userRepo.GetByEmailAndActive(req.Email, true)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 2. Verify password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentilas")
	}

	// 3. Return generated token
	return s.generateAuthResponse(user)
}

func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	// 1. Validate refresh token
	claims, err := utils.ValidateToken(req.RefreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 2. Retrieve and check if token exist or expired
	refreshToken, err := s.userRepo.GetValidRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	// 3. Retrieve user
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 4. Delete token if user is not found
	if err := s.userRepo.DeleteRefreshTokenByID(refreshToken.ID); err != nil {

	}

	// 5. Return generated token
	return s.generateAuthResponse(user)
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.userRepo.DeleteRefreshToken(refreshToken)
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
	if err = s.userRepo.DeleteExpiredTokenByID(user.ID); err != nil {
		return nil, err
	}

	// 3. Insert the new token
	if err := s.userRepo.CreateRefreshToken(&refreshTokenModel); err != nil {
		log.Println()
		_ = err
	}

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
