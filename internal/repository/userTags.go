package repository

import (
	"fmt"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserTagsRepository interface {
	CreateUserTags(userTags *model.UserTags) error
	GetUserTags(userID string) ([]model.UserTags, error)
	GetUserTagsByID(id string, userID string) (*model.UserTags, error)
	UpdateUserTags(userTags *model.UserTags) error
	DeleteUserTags(userTags *model.UserTags) error
}

type userTagsRepository struct {
	db *gorm.DB
}

func NewUserTagsRepository(db *gorm.DB) UserTagsRepository {
	return &userTagsRepository{db: db}
}

func (r *userTagsRepository) CreateUserTags(userTags *model.UserTags) error {
	if userTags.ID == "" {
		userTags.ID = uuid.New().String()
	}
	return r.db.Create(userTags).Error
}

func (r *userTagsRepository) GetUserTags(userID string) ([]model.UserTags, error) {
	var userTags []model.UserTags
	if err := r.db.Where("user_id = ?", userID).Find(&userTags).Error; err != nil {
		return nil, fmt.Errorf("ユーザータグの取得に失敗しました: %w", err)
	}
	return userTags, nil
}

func (r *userTagsRepository) UpdateUserTags(userTags *model.UserTags) error {
	return r.db.Save(userTags).Error
}

func (r *userTagsRepository) GetUserTagsByID(id string, userID string) (*model.UserTags, error) {
	var userTags model.UserTags
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&userTags).Error; err != nil {
		return nil, fmt.Errorf("ユーザータグの取得に失敗しました: %w", err)
	}
	return &userTags, nil
}

func (r *userTagsRepository) DeleteUserTags(userTags *model.UserTags) error {
	return r.db.Delete(userTags).Error
}
