package repository

import (
	"fmt"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
)

type WeaknessAnalysisRepository interface {
	CreateWeaknessAnalysis(userId string, req *model.CreateWeaknessAnalysisRequest) (*model.CreateWeaknessAnalysisResponse, error)
  UpdateWeaknessAnalysis(userId string, req *model.UpdateWeaknessAnalysisRequest) (*model.UpdateWeaknessAnalysisResponse, error)
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
		AnalysisStatus:  "PROCESSING", // 初期状態は処理中
		OverallScore:    0,
		ImprovementRate: 0,
		AnalysisDate:    now,
		AnalyzedAnswers: 0,
		DataPeriodStart: now,                        // TODO: 実際のデータ期間を設定
		DataPeriodEnd:   now,                        // TODO: 実際のデータ期間を設定
		LLMModel:        string(anthropic.ModelClaude3_7Sonnet20250219),
		AnalysisVersion: "1.0",                      // TODO: 設定から取得
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
	}

	return response, nil
}

// UpdateWeaknessAnalysis 学習弱点分析を更新する
func (r *weaknessAnalysisRepository) UpdateWeaknessAnalysis(userId string, req *model.UpdateWeaknessAnalysisRequest) (*model.UpdateWeaknessAnalysisResponse, error) {
  now := time.Now()

  weaknessAnalysis := &model.WeaknessAnalysis{
    ID: req.ID,
    ProjectID: req.ProjectID,
    AnalysisStatus: req.AnalysisStatus,
    OverallScore: req.OverallScore,
    ImprovementRate: req.ImprovementRate,
    AnalysisDate: req.AnalysisDate,
    AnalyzedAnswers: req.AnalyzedAnswers,
    DataPeriodStart: req.DataPeriodStart,
    DataPeriodEnd: req.DataPeriodEnd,
    LLMModel: req.LLMModel,
    AnalysisVersion: req.AnalysisVersion,
    UpdatedAt: now,
  }

  if err := r.db.Save(weaknessAnalysis).Error; err != nil {
    return nil, fmt.Errorf("failed to update weakness analysis: %w", err)
  }

  return &model.UpdateWeaknessAnalysisResponse{
    ID: weaknessAnalysis.ID,
    ProjectID: weaknessAnalysis.ProjectID,
    AnalysisStatus: weaknessAnalysis.AnalysisStatus,
    OverallScore: weaknessAnalysis.OverallScore,
    ImprovementRate: weaknessAnalysis.ImprovementRate,
    AnalysisDate: weaknessAnalysis.AnalysisDate,
    AnalyzedAnswers: weaknessAnalysis.AnalyzedAnswers,
    DataPeriodStart: weaknessAnalysis.DataPeriodStart,
    DataPeriodEnd: weaknessAnalysis.DataPeriodEnd,
    LLMModel: weaknessAnalysis.LLMModel,
    AnalysisVersion: weaknessAnalysis.AnalysisVersion,
  }, nil
}