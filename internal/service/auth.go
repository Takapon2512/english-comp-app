package service

import (
	"errors"
	"log"
	"time"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("無効な認証情報です")
	ErrUserExists         = errors.New("このメールアドレスは既に登録されています")
	ErrInvalidToken       = errors.New("無効なトークンです")
)

type AuthService interface {
	Login(email, password string) (*model.User, error)
	Signup(email, password, name string) (*model.User, error)
	CreateRefreshToken(userID string) (*model.RefreshToken, error)
	ValidateRefreshToken(tokenHash string) (*model.User, error)
	RotateRefreshToken(tokenHash string) (*model.User, *model.RefreshToken, error)
	Logout(tokenHash string) error
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Login(email, password string) (*model.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *authService) Signup(email, password, name string) (*model.User, error) {
	existingUser, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUserID := uuid.New().String()

	user := &model.User{
		ID:           newUserID,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Name:         name,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		CreatedBy:    newUserID,
		UpdatedBy:    newUserID,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) CreateRefreshToken(userID string) (*model.RefreshToken, error) {
	token := &model.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		TokenHash: uuid.New().String(),                 // 実際のプロダクションでは、より安全なトークン生成方法を使用すべき
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30日
		CreatedAt: time.Now(),
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	err := s.userRepo.CreateRefreshToken(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *authService) ValidateRefreshToken(tokenHash string) (*model.User, error) {
	if tokenHash == "" {
		log.Printf("Empty token hash provided")
		return nil, ErrInvalidToken
	}

	token, err := s.userRepo.FindRefreshTokenByHash(tokenHash)
	if err != nil {
		log.Printf("Error finding refresh token: %v", err)
		return nil, err
	}
	if token == nil {
		log.Printf("Token not found: %s", tokenHash)
		return nil, ErrInvalidToken
	}
	if token.ExpiresAt.Before(time.Now()) {
		log.Printf("Token expired: %s, expires: %v", tokenHash, token.ExpiresAt)
		return nil, ErrInvalidToken
	}
	if token.User.ID == "" {
		log.Printf("User not found in token relation: %s", tokenHash)
		return nil, ErrInvalidToken
	}

	return &token.User, nil
}

func (s *authService) RotateRefreshToken(tokenHash string) (*model.User, *model.RefreshToken, error) {
	// 現在のトークンを検証
	user, err := s.ValidateRefreshToken(tokenHash)
	if err != nil {
		return nil, nil, err
	}

	// 現在のトークンを無効化
	err = s.userRepo.RevokeRefreshToken(tokenHash)
	if err != nil {
		log.Printf("Error revoking refresh token: %v", err)
		return nil, nil, err
	}

	// 新しいリフレッシュトークンを生成
	newToken, err := s.CreateRefreshToken(user.ID)
	if err != nil {
		log.Printf("Error creating new refresh token: %v", err)
		return nil, nil, err
	}

	return user, newToken, nil
}

func (s *authService) Logout(tokenHash string) error {
	err := s.userRepo.RevokeRefreshToken(tokenHash)
	if err != nil {
		log.Printf("Error revoking refresh token: %v", err)
		return err
	}

	return nil
}
