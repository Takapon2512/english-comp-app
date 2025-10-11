package model

import (
	"time"

	"gorm.io/gorm"
)

// Project はプロジェクトを表す構造体です
type Project struct {
	ID          string         `json:"id" gorm:"primaryKey;type:char(36)"`
	UserID      string         `json:"user_id" gorm:"type:char(36);not null"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"not null"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	DeletedBy   string         `json:"deleted_by" gorm:"type:char(36)"`
	User        *User          `json:"-" gorm:"foreignKey:UserID"`
	CreatedBy   string         `json:"created_by" gorm:"type:char(36);not null"`
	UpdatedBy   string         `json:"updated_by" gorm:"type:char(36);not null"`

	// リレーション
	Tags []ProjectTag `json:"tags,omitempty" gorm:"foreignKey:ProjectID"`
}

// CreateProjectRequest はプロジェクト作成リクエストを表す構造体です
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=1000"`
}

// CreateProjectResponse はプロジェクト作成レスポンスを表す構造体です
type CreateProjectResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// プロジェクト一覧取得リクエスト
type GetProjectsRequest struct {
	UserID  string `json:"-"` // 内部使用のため、JSONにはシリアライズしない
	Page    int    `form:"page,default=1" binding:"min=0"`
	PerPage int    `form:"per_page,default=10" binding:"min=0,max=100"`
	Tag     string `form:"tag" binding:"omitempty,max=30"`
}

// プロジェクト一覧取得レスポンス
type GetProjectsResponse struct {
	Projects []Project `json:"projects"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PerPage  int       `json:"per_page"`
}

// プロジェクト詳細取得リクエスト
type GetProjectDetailRequest struct {
	ID string `param:"id" binding:"required"`
}

// プロジェクト詳細取得レスポンス
type GetProjectDetailResponse struct {
	Project Project `json:"project"`
}
