package model

import (
	"time"

	"gorm.io/gorm"
)

type CorrectionResults struct {
	ID 							string `json:"id" gorm:"primaryKey;type:char(36)"`
	QuestionAnswerID 			string `json:"question_answer_id" gorm:"type:char(36);not null"`
	QuestionTemplateMasterID 	string `json:"question_template_master_id" gorm:"type:char(36);not null"`
	GetPoints 					int `json:"get_points" gorm:"type:int;not null;default:0"`
	ExampleCorrection 			string `json:"example_correction" gorm:"type:text;null"`
	CorrectRate 				int `json:"correct_rate" gorm:"type:int;null"`
	Advice 						string `json:"advice" gorm:"type:text;null"`
	Status 						string `json:"status" gorm:"type:varchar(20);not null;default:PROCESSING"`
	ChallengeCount 				int `json:"challenge_count" gorm:"type:int;not null;default:1"`
	CreatedBy 					string `json:"created_by" gorm:"type:char(36);not null"`
	UpdatedBy 					string `json:"updated_by" gorm:"type:char(36);not null"`
	DeletedBy 					string `json:"deleted_by" gorm:"type:char(36);null"`
	CreatedAt 					time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt 					time.Time `json:"updated_at" gorm:"not null"`
	DeletedAt 					gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type CorrectionResultsSummary struct {
	ID 							string `json:"id"`
	QuestionAnswerID 			string `json:"question_answer_id"`
	QuestionTemplateMasterID 	string `json:"question_template_master_id"`
	GetPoints 					int `json:"get_points"`
	ExampleCorrection 			string `json:"example_correction"`
	CorrectRate 				int `json:"correct_rate"`
	Advice 						string `json:"advice"`
	Status 						string `json:"status"`
	ChallengeCount 				int `json:"challenge_count"`
	QuestionAnswer 				QuestionAnswersSummary `json:"question_answer"`
	QuestionTemplateMaster 		QuestionTemplateMastersSummary `json:"question_template_master"`
}

type CreateCorrectionResultRequest struct {
	QuestionAnswerID string `json:"question_answer_id" binding:"required"`
	QuestionTemplateMasterID string `json:"question_template_master_id" binding:"required"`
	GetPoints int `json:"get_points"`
	ExampleCorrection string `json:"example_correction"`
	CorrectRate int `json:"correct_rate"`
	Advice string `json:"advice"`
	Status string `json:"status"`
	ChallengeCount int `json:"challenge_count"`
}

type UpdateCorrectionResultRequest struct {
	ID string `json:"id" binding:"required"`
	GetPoints int `json:"get_points"`
	ExampleCorrection string `json:"example_correction"`
	CorrectRate int `json:"correct_rate"`
	Advice string `json:"advice"`
	Status string `json:"status"`
}

type CreateCorrectionResultResponse struct {
	ID string `json:"id"`
	QuestionAnswerID string `json:"question_answer_id"`
	QuestionTemplateMasterID string `json:"question_template_master_id"`
	GetPoints int `json:"get_points"`
	ExampleCorrection string `json:"example_correction"`
	CorrectRate int `json:"correct_rate"`
	Advice string `json:"advice"`
	Status string `json:"status"`
	ChallengeCount int `json:"challenge_count"`
}

type UpdateCorrectionResultResponse struct {
	ID string `json:"id"`
	QuestionAnswerID string `json:"question_answer_id"`
	QuestionTemplateMasterID string `json:"question_template_master_id"`
	GetPoints int `json:"get_points"`
	ExampleCorrection string `json:"example_correction"`
	CorrectRate int `json:"correct_rate"`
	Advice string `json:"advice"`
	Status string `json:"status"`
	ChallengeCount int `json:"challenge_count"`
}

type GrandCorrectResultRequest struct {
	ID string `json:"id" binding:"required"`
}

type GrandCorrectResultResponse struct {
	ID string `json:"id"`
	QuestionAnswerID string `json:"question_answer_id"`
	QuestionTemplateMasterID string `json:"question_template_master_id"`
	GetPoints int `json:"get_points"`
	ExampleCorrection string `json:"example_correction"`
	CorrectRate int `json:"correct_rate"`
	Advice string `json:"advice"`
	Status string `json:"status"`
	ChallengeCount int `json:"challenge_count"`
}