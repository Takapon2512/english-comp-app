package model

import (
	"time"

	"gorm.io/gorm"
)

type QuestionAnswers struct {
	ID                       string         `json:"id" gorm:"primaryKey;type:char(36)"`
	UserID                   string         `json:"user_id" gorm:"type:char(36);not null"`
	ProjectID                string         `json:"project_id" gorm:"type:char(36);not null"`
	QuestionTemplateMasterID string         `json:"question_template_master_id" gorm:"type:char(36);not null"`
	UserAnswer               string         `json:"user_answer" gorm:"type:text"`
	ChallengeCount           int            `json:"challenge_count" gorm:"type:int;not null;default:1"`
	Status                   string         `json:"status" gorm:"type:varchar(20);not null;default:'PROCESSING'"`
	CreatedAt                time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt                time.Time      `json:"updated_at" gorm:"not null"`
	DeletedAt                gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	DeletedBy                string         `json:"deleted_by" gorm:"type:char(36)"`
	CreatedBy                string         `json:"created_by" gorm:"type:char(36);not null"`
	UpdatedBy                string         `json:"updated_by" gorm:"type:char(36);not null"`
}

type QuestionAnswersSummary struct {
	ID                       string `json:"id"`
	UserID                   string `json:"user_id"`
	ProjectID                string `json:"project_id"`
	QuestionTemplateMasterID string `json:"question_template_master_id"`
	UserAnswer               string `json:"user_answer"`
}

type CreateQuestionAnswersRequest struct {
	ProjectID                string `json:"project_id" binding:"required"`
	QuestionTemplateMasterID string `json:"question_template_master_id" binding:"required"`
	UserAnswer               string `json:"user_answer"`
}

type CreateQuestionAnswersResponse struct {
	ID                       string `json:"id"`
	UserID                   string `json:"user_id"`
	ProjectID                string `json:"project_id"`
	QuestionTemplateMasterID string `json:"question_template_master_id"`
	UserAnswer               string `json:"user_answer"`
	ChallengeCount           int    `json:"challenge_count"`
}

type GetQuestionAnswersRequest struct {
	ProjectID string `json:"project_id" binding:"required"`
}

type GetQuestionAnswersResponse struct {
	QuestionAnswers []QuestionAnswers `json:"question_answers"`
}

type UpdateQuestionAnswersFinishResponse struct {
	QuestionAnswers []QuestionAnswers `json:"question_answers"`
}

type GetProjectQuestionToAnswerResponse struct {
	Question QuestionTemplateMastersSummary `json:"question"`
	NowQuestionNumber int `json:"now_question_number"`
}