package model

import (
	"time"

	"gorm.io/gorm"
)

type ProjectQuestions struct {
	ID                       string         `json:"id" gorm:"primaryKey;type:char(36)"`
	ProjectID                string         `json:"project_id" gorm:"type:char(36);not null"`
	QuestionTemplateMasterID string         `json:"question_template_master_id" gorm:"type:char(36);not null"`
	CreatedAt                time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt                time.Time      `json:"updated_at" gorm:"not null"`
	DeletedAt                gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	DeletedBy                string         `json:"deleted_by" gorm:"type:char(36)"`
	CreatedBy                string         `json:"created_by" gorm:"type:char(36);not null"`
	UpdatedBy                string         `json:"updated_by" gorm:"type:char(36);not null"`
}

type ProjectQuestionsSummary struct {
	ID                       string                         `json:"id"`
	ProjectID                string                         `json:"project_id"`
	QuestionTemplateMasterID string                         `json:"question_template_master_id"`
	Questions                QuestionTemplateMastersSummary `json:"questions"`
}

type CreateProjectQuestionsRequest struct {
	UserID                    string   `json:"user_id"`
	ProjectID                 string   `json:"project_id" binding:"required"`
	QuestionTemplateMasterIDs []string `json:"question_template_master_ids" binding:"required"`
}

type CreateProjectQuestionsResponse struct {
	ProjectQuestions []ProjectQuestionsSummary `json:"project_questions"`
}
