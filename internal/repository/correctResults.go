package repository

import (
	"fmt"
	"time"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CorrectResultsRepository interface {
	CreateCorrectionResult(req *model.CreateCorrectionResultRequest) (*model.CreateCorrectionResultResponse, error)
	UpdateCorrectionResult(req *model.UpdateCorrectionResultRequest) (*model.UpdateCorrectionResultResponse, error)

	GetCorrectionResultById(id string) (*model.CorrectionResults, error)
}

type correctResultsRepository struct {
	db *gorm.DB
}

func NewCorrectResultsRepository(db *gorm.DB) CorrectResultsRepository {
	return &correctResultsRepository{db: db}
}

func (r *correctResultsRepository) CreateCorrectionResult(req *model.CreateCorrectionResultRequest) (*model.CreateCorrectionResultResponse, error) {
	// 
	
	now := time.Now()
	correctionResult := &model.CorrectionResults{
		ID:                       uuid.New().String(),
		QuestionAnswerID:         req.QuestionAnswerID,
		QuestionTemplateMasterID: req.QuestionTemplateMasterID,
		GetPoints:                req.GetPoints,
		ExampleCorrection:        req.ExampleCorrection,
		CorrectRate:              req.CorrectRate,
		Advice:                   req.Advice,
		Status:                   req.Status,
		ChallengeCount:           req.ChallengeCount,
		CreatedAt:                now,
		UpdatedAt:                now,
		CreatedBy:                "system",
		UpdatedBy:                "system",
	}

	if err := r.db.Create(correctionResult).Error; err != nil {
		return nil, fmt.Errorf("修正結果の作成に失敗しました: %w", err)
	}

	return &model.CreateCorrectionResultResponse{
		ID:                       correctionResult.ID,
		QuestionAnswerID:         correctionResult.QuestionAnswerID,
		QuestionTemplateMasterID: correctionResult.QuestionTemplateMasterID,
		GetPoints:                correctionResult.GetPoints,
		ExampleCorrection:        correctionResult.ExampleCorrection,
		CorrectRate:              correctionResult.CorrectRate,
		Advice:                   correctionResult.Advice,
		Status:                   correctionResult.Status,
		ChallengeCount:           correctionResult.ChallengeCount,
	}, nil
}

func (r *correctResultsRepository) UpdateCorrectionResult(req *model.UpdateCorrectionResultRequest) (*model.UpdateCorrectionResultResponse, error) {
	now := time.Now()
	correctionResult := &model.CorrectionResults{
		ID:                       req.ID,
		GetPoints:                req.GetPoints,
		ExampleCorrection:        req.ExampleCorrection,
		CorrectRate:              req.CorrectRate,
		Advice:                   req.Advice,
		Status:                   req.Status,
		UpdatedAt:                now,
		UpdatedBy:                "system",
	}

	if err := r.db.Model(&model.CorrectionResults{}).Where("id = ?", req.ID).Updates(correctionResult).Error; err != nil {
		return nil, fmt.Errorf("添削結果の更新に失敗しました: %w", err)
	}

	return &model.UpdateCorrectionResultResponse{
		ID:                       correctionResult.ID,
		QuestionAnswerID:         correctionResult.QuestionAnswerID,
		QuestionTemplateMasterID: correctionResult.QuestionTemplateMasterID,
		GetPoints:                correctionResult.GetPoints,
		ExampleCorrection:        correctionResult.ExampleCorrection,
		CorrectRate:              correctionResult.CorrectRate,
		Advice:                   correctionResult.Advice,
		Status:                   correctionResult.Status,
		ChallengeCount:           correctionResult.ChallengeCount,
	}, nil
}

func (r *correctResultsRepository) GetCorrectionResultById(id string) (*model.CorrectionResults, error) {
	var correctionResult model.CorrectionResults
	if err := r.db.Model(&model.CorrectionResults{}).Where("id = ?", id).First(&correctionResult).Error; err != nil {
		return nil, fmt.Errorf("添削結果の取得に失敗しました: %w", err)
	}

	return &correctionResult, nil
}
