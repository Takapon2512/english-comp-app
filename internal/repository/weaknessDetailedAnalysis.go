package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
)

type WeaknessDetailedAnalysisRepository interface {
	CreateWeaknessDetailedAnalysis(userId string, req *model.CreateWeaknessDetailedAnalysisRequest) (*model.CreateWeaknessDetailedAnalysisResponse, error)
	GetWeaknessDetailedAnalysis(analysisId string) (*model.WeaknessDetailedAnalysisSummary, error)
	UpdateWeaknessDetailedAnalysis(userId string, req *model.UpdateWeaknessDetailedAnalysisRequest) (*model.UpdateWeaknessDetailedAnalysisResponse, error)
}

type weaknessDetailedAnalysisRepository struct {
	db *gorm.DB
}

func NewWeaknessDetailedAnalysisRepository(db *gorm.DB) WeaknessDetailedAnalysisRepository {
	return &weaknessDetailedAnalysisRepository{db: db}
}

// CreateWeaknessDetailedAnalysis 詳細分析結果を作成する
func (r *weaknessDetailedAnalysisRepository) CreateWeaknessDetailedAnalysis(userId string, req *model.CreateWeaknessDetailedAnalysisRequest) (*model.CreateWeaknessDetailedAnalysisResponse, error) {
	now := time.Now()

	weaknessDetailedAnalysis := &model.WeaknessDetailedAnalysis{
		ID:                    uuid.New().String(),
		AnalysisID:            req.AnalysisID,
		GrammarScore:          req.GrammarScore,
		GrammarDescription:    req.GrammarDescription,
		GrammarExamples:       req.GrammarExamples,
		VocabularyScore:       req.VocabularyScore,
		VocabularyDescription: req.VocabularyDescription,
		VocabularyExamples:    req.VocabularyExamples,
		ExpressionScore:       req.ExpressionScore,
		ExpressionDescription: req.ExpressionDescription,
		ExpressionExamples:    req.ExpressionExamples,
		StructureScore:        req.StructureScore,
		StructureDescription:  req.StructureDescription,
		StructureExamples:     req.StructureExamples,
		CreatedAt:             now,
		UpdatedAt:             now,
		CreatedBy:             userId,
		UpdatedBy:             userId,
	}

	if err := r.db.Create(weaknessDetailedAnalysis).Error; err != nil {
		return nil, fmt.Errorf("failed to create weakness detailed analysis: %w", err)
	}

	return &model.CreateWeaknessDetailedAnalysisResponse{
		ID:                    weaknessDetailedAnalysis.ID,
		AnalysisID:            weaknessDetailedAnalysis.AnalysisID,
		GrammarScore:          weaknessDetailedAnalysis.GrammarScore,
		GrammarDescription:    weaknessDetailedAnalysis.GrammarDescription,
		GrammarExamples:       weaknessDetailedAnalysis.GrammarExamples,
		VocabularyScore:       weaknessDetailedAnalysis.VocabularyScore,
		VocabularyDescription: weaknessDetailedAnalysis.VocabularyDescription,
		VocabularyExamples:    weaknessDetailedAnalysis.VocabularyExamples,
		ExpressionScore:       weaknessDetailedAnalysis.ExpressionScore,
		ExpressionDescription: weaknessDetailedAnalysis.ExpressionDescription,
		ExpressionExamples:    weaknessDetailedAnalysis.ExpressionExamples,
	}, nil
}

// UpdateWeaknessDetailedAnalysis 詳細分析結果を更新する
func (r *weaknessDetailedAnalysisRepository) UpdateWeaknessDetailedAnalysis(userId string, req *model.UpdateWeaknessDetailedAnalysisRequest) (*model.UpdateWeaknessDetailedAnalysisResponse, error) {
	now := time.Now()

	weaknessDetailedAnalysis := &model.WeaknessDetailedAnalysis{
		ID:                    req.ID,
		AnalysisID:            req.AnalysisID,
		GrammarScore:          req.GrammarScore,
		GrammarDescription:    req.GrammarDescription,
		GrammarExamples:       req.GrammarExamples,
		VocabularyScore:       req.VocabularyScore,
		VocabularyDescription: req.VocabularyDescription,
		VocabularyExamples:    req.VocabularyExamples,
		ExpressionScore:       req.ExpressionScore,
		ExpressionDescription: req.ExpressionDescription,
		ExpressionExamples:    req.ExpressionExamples,
		StructureScore:        req.StructureScore,
		StructureDescription:  req.StructureDescription,
		StructureExamples:     req.StructureExamples,
		UpdatedAt:             now,
		UpdatedBy:             userId,
	}

	if err := r.db.Save(weaknessDetailedAnalysis).Error; err != nil {
		return nil, fmt.Errorf("failed to update weakness detailed analysis: %w", err)
	}

	return &model.UpdateWeaknessDetailedAnalysisResponse{
		ID:                    weaknessDetailedAnalysis.ID,
		AnalysisID:            weaknessDetailedAnalysis.AnalysisID,
		GrammarScore:          weaknessDetailedAnalysis.GrammarScore,
		GrammarDescription:    weaknessDetailedAnalysis.GrammarDescription,
		GrammarExamples:       weaknessDetailedAnalysis.GrammarExamples,
		VocabularyScore:       weaknessDetailedAnalysis.VocabularyScore,
		VocabularyDescription: weaknessDetailedAnalysis.VocabularyDescription,
		VocabularyExamples:    weaknessDetailedAnalysis.VocabularyExamples,
		ExpressionScore:       weaknessDetailedAnalysis.ExpressionScore,
		ExpressionDescription: weaknessDetailedAnalysis.ExpressionDescription,
		ExpressionExamples:    weaknessDetailedAnalysis.ExpressionExamples,
		StructureScore:        weaknessDetailedAnalysis.StructureScore,
		StructureDescription:  weaknessDetailedAnalysis.StructureDescription,
		StructureExamples:     weaknessDetailedAnalysis.StructureExamples,
	}, nil
}

// GetWeaknessDetailedAnalysis 詳細分析結果を取得する
func (r *weaknessDetailedAnalysisRepository) GetWeaknessDetailedAnalysis(analysisId string) (*model.WeaknessDetailedAnalysisSummary, error) {
	var weaknessDetailedAnalysis model.WeaknessDetailedAnalysisSummary

	// 正しいテーブル名を指定してクエリを実行
	if err := r.db.Table("weakness_detailed_analyses").Where("analysis_id = ?", analysisId).First(&weaknessDetailedAnalysis).Error; err != nil {
		return nil, fmt.Errorf("failed to get weakness detailed analysis: %w", err)
	}

	return &model.WeaknessDetailedAnalysisSummary{
		ID:                    weaknessDetailedAnalysis.ID,
		AnalysisID:            weaknessDetailedAnalysis.AnalysisID,
		GrammarScore:          weaknessDetailedAnalysis.GrammarScore,
		GrammarDescription:    weaknessDetailedAnalysis.GrammarDescription,
		GrammarExamples:       weaknessDetailedAnalysis.GrammarExamples,
		VocabularyScore:       weaknessDetailedAnalysis.VocabularyScore,
		VocabularyDescription: weaknessDetailedAnalysis.VocabularyDescription,
		VocabularyExamples:    weaknessDetailedAnalysis.VocabularyExamples,
		ExpressionScore:       weaknessDetailedAnalysis.ExpressionScore,
		ExpressionDescription: weaknessDetailedAnalysis.ExpressionDescription,
		ExpressionExamples:    weaknessDetailedAnalysis.ExpressionExamples,
		StructureScore:        weaknessDetailedAnalysis.StructureScore,
		StructureDescription:  weaknessDetailedAnalysis.StructureDescription,
		StructureExamples:     weaknessDetailedAnalysis.StructureExamples,
	}, nil
}
