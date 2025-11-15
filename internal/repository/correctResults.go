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
	GetCorrectResults(req *model.GetCorrectResultsRequest) (*model.GetCorrectResultsResponse, error)
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
		ProjectID:                req.ProjectID,
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
		ProjectID:                correctionResult.ProjectID,
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
		ProjectID:                correctionResult.ProjectID,
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

// 添削結果の一覧取得
func (r *correctResultsRepository) GetCorrectResults(req *model.GetCorrectResultsRequest) (*model.GetCorrectResultsResponse, error) {
	var correctResults []model.CorrectionResults;

	// challenge_countが指定されていない場合は、最大値を取得する
	challengeCount := req.ChallengeCount
	if challengeCount == 0 {
		var maxChallengeCount int
		if err := r.db.Model(&model.CorrectionResults{}).Where("project_id = ?", req.ProjectID).Select("COALESCE(MAX(challenge_count), 1)").Scan(&maxChallengeCount).Error; err != nil {
			return nil, fmt.Errorf("最大挑戦回数の取得に失敗しました: %w", err)
		}
		challengeCount = maxChallengeCount
	}

	if err := r.db.Model(&model.CorrectionResults{}).Where("project_id = ? AND challenge_count = ?", req.ProjectID, challengeCount).Find(&correctResults).Error; err != nil {
		return nil, fmt.Errorf("添削結果の取得に失敗しました: %w", err)
	}

	var correctResultsSummary []model.CorrectionResultsSummary;
	for _, correctResult := range correctResults {
		correctResultsSummary = append(correctResultsSummary, model.CorrectionResultsSummary{
			ID: correctResult.ID,
			QuestionAnswerID: correctResult.QuestionAnswerID,
			QuestionTemplateMasterID: correctResult.QuestionTemplateMasterID,
			ProjectID: correctResult.ProjectID,
			GetPoints: correctResult.GetPoints,
			ExampleCorrection: correctResult.ExampleCorrection,
			CorrectRate: correctResult.CorrectRate,
			Advice: correctResult.Advice,
			Status: correctResult.Status,
			ChallengeCount: correctResult.ChallengeCount,
		})
	}

	return &model.GetCorrectResultsResponse{
		CorrectResults: correctResultsSummary,
	}, nil
}