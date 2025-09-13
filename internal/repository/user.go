package repository

import (
	"errors"
	"time"

	"github.com/Takanpon2512/english-app/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	FindByID(id string) (*model.User, error)
	Create(user *model.User) error
	CreateRefreshToken(token *model.RefreshToken) error
	FindRefreshTokenByHash(tokenHash string) (*model.RefreshToken, error)
	RevokeRefreshToken(tokenID string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	result := r.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) CreateRefreshToken(token *model.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *userRepository) FindRefreshTokenByHash(tokenHash string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	result := r.db.Joins("JOIN users ON users.id = refresh_tokens.user_id").
		Where("refresh_tokens.token_hash = ? AND refresh_tokens.revoked_at IS NULL", tokenHash).
		Preload("User").
		First(&token)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &token, nil
}

func (r *userRepository) RevokeRefreshToken(tokenHash string) error {
	return r.db.Model(&model.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked_at", time.Now()).
		Error
}
