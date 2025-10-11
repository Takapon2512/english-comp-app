package model

import (
	"time"

	"gorm.io/gorm"
)

// ProjectTag はプロジェクトタグを表す構造体です
type ProjectTag struct {
	ID         string         `json:"id" gorm:"primaryKey;type:char(36)"`
	ProjectID  string         `json:"project_id" gorm:"type:char(36);not null"`
	UserTagsID string         `json:"user_tags_id" gorm:"type:char(36);not null"`
	CreatedAt  time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"not null"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	CreatedBy  string         `json:"created_by" gorm:"type:char(36);not null"`
	UpdatedBy  string         `json:"updated_by" gorm:"type:char(36);not null"`
	DeletedBy  string         `json:"deleted_by" gorm:"type:char(36)"`

	// リレーション
	Project *Project  `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	UserTag *UserTags `json:"user_tag,omitempty" gorm:"foreignKey:UserTagsID"`
}
