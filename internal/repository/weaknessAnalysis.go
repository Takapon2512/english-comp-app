package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
)

type WeaknessAnalysisRepository interface {
	CreateWeaknessAnalysis(userId string, req *model.CreateWeaknessAnalysisRequest) (*model.CreateWeaknessAnalysisResponse, error)
	GetWeaknessAnalysis(userId string, req *model.GetWeaknessAnalysisRequest) (*model.GetWeaknessAnalysisResponse, error)
	GetWeaknessAnalysisStatusSummary(userId string, analysisId string) (*model.WeaknessAnalysisStatusSummary, error)
	UpdateAnalysisStatus(analysisId string, status string) error
	UpdateOverallScore(analysisId string, overallScore int) error
}

type weaknessAnalysisRepository struct {
	db *gorm.DB
}

func NewWeaknessAnalysisRepository(db *gorm.DB) WeaknessAnalysisRepository {
	return &weaknessAnalysisRepository{db: db}
}

// CreateWeaknessAnalysis 学習弱点分析を作成する
func (r *weaknessAnalysisRepository) CreateWeaknessAnalysis(userId string, req *model.CreateWeaknessAnalysisRequest) (*model.CreateWeaknessAnalysisResponse, error) {
	now := time.Now()

	// WeaknessAnalysisエンティティを作成
	weaknessAnalysis := &model.WeaknessAnalysis{
		ID:              uuid.New().String(),
		ProjectID:       req.ProjectID,
		UserID:          userId,       // 外部キー制約のため必須
		AnalysisStatus:  "PROCESSING", // 初期状態は処理中
		OverallScore:    0,
		ImprovementRate: 0,
		AnalysisDate:    now,
		AnalyzedAnswers: 0,
		DataPeriodStart: now, // TODO: 実際のデータ期間を設定
		DataPeriodEnd:   now, // TODO: 実際のデータ期間を設定
		LLMModel:        string(anthropic.ModelClaude3_7Sonnet20250219),
		AnalysisVersion: "1.0", // TODO: 設定から取得
		CreatedAt:       now,
		UpdatedAt:       now,
		CreatedBy:       userId,
		UpdatedBy:       userId,
	}

	// データベースに保存
	if err := r.db.Create(weaknessAnalysis).Error; err != nil {
		return nil, fmt.Errorf("failed to create weakness analysis: %w", err)
	}

	// レスポンスを作成
	response := &model.CreateWeaknessAnalysisResponse{
		ID:             weaknessAnalysis.ID,
		ProjectID:      weaknessAnalysis.ProjectID,
		AnalysisStatus: weaknessAnalysis.AnalysisStatus,
		OverallScore:   weaknessAnalysis.OverallScore,
	}

	return response, nil
}

// UpdateWeaknessAnalysis 学習弱点分析を更新する
func (r *weaknessAnalysisRepository) UpdateWeaknessAnalysis(userId string, req *model.UpdateWeaknessAnalysisRequest) (*model.UpdateWeaknessAnalysisResponse, error) {
	now := time.Now()

	weaknessAnalysis := &model.WeaknessAnalysis{
		ID:              req.ID,
		ProjectID:       req.ProjectID,
		AnalysisStatus:  req.AnalysisStatus,
		OverallScore:    req.OverallScore,
		ImprovementRate: req.ImprovementRate,
		AnalysisDate:    req.AnalysisDate,
		AnalyzedAnswers: req.AnalyzedAnswers,
		DataPeriodStart: req.DataPeriodStart,
		DataPeriodEnd:   req.DataPeriodEnd,
		LLMModel:        req.LLMModel,
		AnalysisVersion: req.AnalysisVersion,
		UpdatedAt:       now,
	}

	if err := r.db.Save(weaknessAnalysis).Error; err != nil {
		return nil, fmt.Errorf("failed to update weakness analysis: %w", err)
	}

	return &model.UpdateWeaknessAnalysisResponse{
		ID:              weaknessAnalysis.ID,
		ProjectID:       weaknessAnalysis.ProjectID,
		AnalysisStatus:  weaknessAnalysis.AnalysisStatus,
		OverallScore:    weaknessAnalysis.OverallScore,
		ImprovementRate: weaknessAnalysis.ImprovementRate,
		AnalysisDate:    weaknessAnalysis.AnalysisDate,
		AnalyzedAnswers: weaknessAnalysis.AnalyzedAnswers,
		DataPeriodStart: weaknessAnalysis.DataPeriodStart,
		DataPeriodEnd:   weaknessAnalysis.DataPeriodEnd,
		LLMModel:        weaknessAnalysis.LLMModel,
		AnalysisVersion: weaknessAnalysis.AnalysisVersion,
	}, nil
}

// GetWeaknessAnalysis 学習弱点分析を取得する
func (r *weaknessAnalysisRepository) GetWeaknessAnalysis(userId string, req *model.GetWeaknessAnalysisRequest) (*model.GetWeaknessAnalysisResponse, error) {
	var weaknessAnalysis model.WeaknessAnalysis

	if err := r.db.Where("user_id = ? AND project_id = ?", userId, req.ProjectID).First(&weaknessAnalysis).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // レコードが見つからない場合はnilを返す
		}
		return nil, fmt.Errorf("failed to get weakness analysis: %w", err)
	}

	return &model.GetWeaknessAnalysisResponse{
		Analysis: model.WeaknessAnalysisSummary{
			ID:              weaknessAnalysis.ID,
			ProjectID:       weaknessAnalysis.ProjectID,
			AnalysisStatus:  weaknessAnalysis.AnalysisStatus,
			OverallScore:    weaknessAnalysis.OverallScore,
			ImprovementRate: weaknessAnalysis.ImprovementRate,
			AnalysisDate:    weaknessAnalysis.AnalysisDate,
			AnalyzedAnswers: weaknessAnalysis.AnalyzedAnswers,
			DataPeriodStart: weaknessAnalysis.DataPeriodStart,
			DataPeriodEnd:   weaknessAnalysis.DataPeriodEnd,
		},
	}, nil
}

// UpdateAnalysisStatus 分析ステータスを更新する
func (r *weaknessAnalysisRepository) UpdateAnalysisStatus(analysisId string, status string) error {
	now := time.Now()

	result := r.db.Model(&model.WeaknessAnalysis{}).
		Where("id = ?", analysisId).
		Updates(map[string]interface{}{
			"analysis_status": status,
			"updated_at":      now,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update analysis status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("analysis not found with id: %s", analysisId)
	}

	return nil
}

// UpdateOverallScore 総合スコアを更新する
func (r *weaknessAnalysisRepository) UpdateOverallScore(analysisId string, overallScore int) error {
	now := time.Now()

	result := r.db.Model(&model.WeaknessAnalysis{}).
		Where("id = ?", analysisId).
		Updates(map[string]interface{}{
			"overall_score": overallScore,
			"updated_at":    now,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update overall score: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("analysis not found with id: %s", analysisId)
	}

	return nil
}

// GetWeaknessAnalysisStatusSummary 分析状況のサマリーを取得する
func (r *weaknessAnalysisRepository) GetWeaknessAnalysisStatusSummary(userId string, analysisId string) (*model.WeaknessAnalysisStatusSummary, error) {
	var weaknessAnalysis model.WeaknessAnalysis

	if err := r.db.Where("user_id = ? AND id = ?", userId, analysisId).First(&weaknessAnalysis).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // レコードが見つからない場合はnilを返す
		}
		return nil, fmt.Errorf("failed to get weakness analysis status summary: %w", err)
	}

	return &model.WeaknessAnalysisStatusSummary{
		ID:             weaknessAnalysis.ID,
		ProjectID:      weaknessAnalysis.ProjectID,
		UserID:         weaknessAnalysis.UserID,
		AnalysisStatus: weaknessAnalysis.AnalysisStatus,
	}, nil
}
