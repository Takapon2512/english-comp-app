package model

import (
	"time"

	"gorm.io/gorm"
)

type UserTags struct {
	ID        string         `json:"id" gorm:"primaryKey;type:char(36)"`
	UserID    string         `json:"user_id" gorm:"type:char(36);not null"`
	Name      string         `json:"name" gorm:"type:varchar(30);not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null"`
	CreatedBy string         `json:"created_by" gorm:"type:char(36);not null"`
	UpdatedBy string         `json:"updated_by" gorm:"type:char(36);not null"`
	DeletedBy *string        `json:"deleted_by" gorm:"type:char(36)"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// createUserTagsRequestはユーザータグ作成リクエストを表す構造体です
type CreateUserTagsRequest struct {
	Name string `json:"name" binding:"required,max=30"`
}

// createUserTagsResponseはユーザータグ作成レスポンスを表す構造体です
type CreateUserTagsResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// getUserTagsRequestはユーザータグ取得リクエストを表す構造体です
type GetUserTagsRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// UserTagsSummary はユーザータグの要約情報を表す構造体です
type UserTagsSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// getUserTagsResponseはユーザータグ取得レスポンスを表す構造体です
type GetUserTagsResponse struct {
	UserTags []UserTagsSummary `json:"user_tags"`
}

// updateUserTagsRequestはユーザータグ更新リクエストを表す構造体です
type UpdateUserTagsRequest struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required,max=30"`
}

// updateUserTagsResponseはユーザータグ更新レスポンスを表す構造体です
type UpdateUserTagsResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// deleteUserTagsRequestはユーザータグ削除リクエストを表す構造体です
type DeleteUserTagsRequest struct {
	ID string `json:"id" binding:"required"`
}

// deleteUserTagsResponseはユーザータグ削除レスポンスを表す構造体です
type DeleteUserTagsResponse struct {
	ID string `json:"id"`
}
