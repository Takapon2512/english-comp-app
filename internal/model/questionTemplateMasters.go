package model

import (
	"time"

	"gorm.io/gorm"
)

type QuestionTemplateMasters struct {
	ID         		string         `json:"id" gorm:"primaryKey;type:char(36)"`
	CategoryID 		string         `json:"category_id" gorm:"type:char(36);not null"`
	QuestionType 	string         `json:"question_type" gorm:"type:varchar(10);not null"`
	English    		string         `json:"english" gorm:"type:text;not null"`
	Japanese   		string         `json:"japanese" gorm:"type:text;not null"`
	Status     		string         `json:"status" gorm:"type:varchar(10);not null"`
	Level      		string         `json:"level" gorm:"type:varchar(10);not null"`
	EstimatedTime 	int            `json:"estimated_time" gorm:"type:int;not null"`
	Points       	int            `json:"points" gorm:"type:int;not null"`
	CreatedBy  		string         `json:"created_by" gorm:"type:char(36);not null"`
	UpdatedBy  		string         `json:"updated_by" gorm:"type:char(36);not null"`
	DeletedBy  		string         `json:"deleted_by" gorm:"type:char(36)"`
	CreatedAt  		time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt  		time.Time      `json:"updated_at" gorm:"not null"`
	DeletedAt  		gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type QuestionTemplateMastersSummary struct {
	ID        string `json:"id"`
	CategoryID string `json:"category_id"`
	QuestionType string `json:"question_type"`
	English    string `json:"english"`
	Japanese   string `json:"japanese"`
	Status     string `json:"status"`
	Level      string `json:"level"`
	EstimatedTime int `json:"estimated_time"`
	Points int `json:"points"`
}

type GetQuestionTemplateMastersSearchRequest struct {
	CategoryID string `json:"category_id"`
	QuestionType string `json:"question_type"`
	Status     string `json:"status"`
	Level      string `json:"level"`
	EstimatedTime int `json:"estimated_time"`
	Points int `json:"points"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
}

type GetQuestionTemplateMastersSearchResponse struct {
	QuestionTemplateMasters []QuestionTemplateMastersSummary `json:"question_template_masters"`
	Total                   int                              `json:"total"`
	Page                    int                              `json:"page"`
	PerPage                 int                              `json:"per_page"`
}
