package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/repository"
)

type UserTagsService interface {
	CreateUserTags(userID string, req *model.CreateUserTagsRequest) (*model.CreateUserTagsResponse, error)
	GetUserTags(userID string) (*model.GetUserTagsResponse, error)
	UpdateUserTags(userID string, req *model.UpdateUserTagsRequest) (*model.UpdateUserTagsResponse, error)
	DeleteUserTags(userID string, req *model.DeleteUserTagsRequest) (*model.DeleteUserTagsResponse, error)
}

type userTagsService struct {
	db   *gorm.DB
	repo repository.UserTagsRepository
}

func NewUserTagsService(db *gorm.DB, repo repository.UserTagsRepository) UserTagsService {
	return &userTagsService{
		db:   db,
		repo: repo,
	}
}

// createUserTags ユーザータグを作成する
func (s *userTagsService) CreateUserTags(userID string, req *model.CreateUserTagsRequest) (*model.CreateUserTagsResponse, error) {
	now := time.Now()
	userTags := &model.UserTags{
		UserID:    userID,
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {

		if err := s.repo.CreateUserTags(userTags); err != nil {
			return fmt.Errorf("ユーザータグの作成に失敗しました: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &model.CreateUserTagsResponse{
		ID:        userTags.ID,
		Name:      userTags.Name,
		CreatedAt: userTags.CreatedAt,
	}, nil
}

// getUserTags ユーザータグを取得する
func (s *userTagsService) GetUserTags(userID string) (*model.GetUserTagsResponse, error) {
	userTags, err := s.repo.GetUserTags(userID)
	if err != nil {
		return nil, err
	}

	// UserTagsからUserTagsSummaryに変換
	userTagsSummary := make([]model.UserTagsSummary, len(userTags))
	for i, tag := range userTags {
		userTagsSummary[i] = model.UserTagsSummary{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	return &model.GetUserTagsResponse{UserTags: userTagsSummary}, nil
}

// updateUserTags ユーザータグを更新する
func (s *userTagsService) UpdateUserTags(userID string, req *model.UpdateUserTagsRequest) (*model.UpdateUserTagsResponse, error) {
	// 既存のレコードを取得
	existingUserTags, err := s.repo.GetUserTagsByID(req.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("ユーザータグが見つかりません: %w", err)
	}

	now := time.Now()
	userTags := &model.UserTags{
		ID:        req.ID,
		UserID:    userID,
		Name:      req.Name,
		CreatedAt: existingUserTags.CreatedAt, // 既存の値を保持
		UpdatedAt: now,
		CreatedBy: existingUserTags.CreatedBy, // 既存の値を保持
		UpdatedBy: userID,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.UpdateUserTags(userTags); err != nil {
			return fmt.Errorf("ユーザータグの更新に失敗しました: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &model.UpdateUserTagsResponse{
		ID:        userTags.ID,
		Name:      userTags.Name,
		CreatedAt: userTags.CreatedAt,
	}, nil
}

// deleteUserTags ユーザータグを削除する
func (s *userTagsService) DeleteUserTags(userId string, req *model.DeleteUserTagsRequest) (*model.DeleteUserTagsResponse, error) {
	now := time.Now()

	// 既存のレコードを取得
	existingUserTags, err := s.repo.GetUserTagsByID(req.ID, userId)
	if err != nil {
		return nil, fmt.Errorf("ユーザータグが見つかりません: %w", err)
	}

	userTags := &model.UserTags{
		ID:        req.ID,
		UserID:    userId,
		Name:      existingUserTags.Name,
		CreatedAt: existingUserTags.CreatedAt,
		UpdatedAt: existingUserTags.UpdatedAt,
		CreatedBy: existingUserTags.CreatedBy,
		UpdatedBy: existingUserTags.UpdatedBy,
		DeletedAt: gorm.DeletedAt{Time: now},
		DeletedBy: userId,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.DeleteUserTags(userTags); err != nil {
			return fmt.Errorf("ユーザータグの削除に失敗しました: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &model.DeleteUserTagsResponse{
		ID: userTags.ID,
	}, nil
}
