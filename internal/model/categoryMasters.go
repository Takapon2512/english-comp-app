package model

import (
	"time"

	"gorm.io/gorm"
)

type CategoryMasters struct {
	ID        string         `json:"id" gorm:"primaryKey;type:char(36)"`
	Name      string         `json:"name" gorm:"type:varchar(30);not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null"`
	CreatedBy string         `json:"created_by" gorm:"type:char(36);not null"`
	UpdatedBy string         `json:"updated_by" gorm:"type:char(36);not null"`
	DeletedBy string         `json:"deleted_by" gorm:"type:char(36)"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type CategoryMastersSummary struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
}

type GetCategoryMastersSearchRequest struct {
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}

type GetCategoryMastersByIDRequest struct {
	ID string `json:"id"`
}

type GetCategoryMastersByIDResponse struct {
	CategoryMasters CategoryMastersSummary `json:"category_masters"`
}

type GetCategoryMastersSearchResponse struct {
	CategoryMasters []CategoryMastersSummary `json:"category_masters"`
	Total           int                      `json:"total"`
	Page            int                      `json:"page"`
	PerPage         int                      `json:"per_page"`
}
