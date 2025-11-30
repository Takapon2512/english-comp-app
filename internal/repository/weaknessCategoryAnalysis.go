package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/utils"
)

type WeaknessCategoryAnalysisRepository interface {
	CreateWeaknessCategoryAnalysis(userId string, req *model.CreateWeaknessCategoryAnalysisRequest) (*model.CreateWeaknessCategoryAnalysisResponse, error)
	GetWeaknessCategoryAnalysis(analysisId string) ([]*model.WeaknessCategoryAnalysisResponse, error)
	UpdateWeaknessCategoryAnalysis(userId string, req *model.UpdateWeaknessCategoryAnalysisRequest) (*model.UpdateWeaknessCategoryAnalysisResponse, error)
}

type weaknessCategoryAnalysisRepository struct {
	db *gorm.DB
}

func NewWeaknessCategoryAnalysisRepository(db *gorm.DB) WeaknessCategoryAnalysisRepository {
	return &weaknessCategoryAnalysisRepository{db: db}
}

// CreateWeaknessCategoryAnalysis カテゴリ別の分析結果を作成する
func (r *weaknessCategoryAnalysisRepository) CreateWeaknessCategoryAnalysis(userId string, req *model.CreateWeaknessCategoryAnalysisRequest) (*model.CreateWeaknessCategoryAnalysisResponse, error) {
	now := time.Now()

	weaknessCategoryAnalysis := &model.WeaknessCategoryAnalysis{
		ID:           uuid.New().String(),
		AnalysisID:   req.AnalysisID,
		CategoryID:   req.CategoryID,
		CategoryName: req.CategoryName,
		Score:        req.Score,
		IsWeakness:   req.IsWeakness,
		IsStrength:   req.IsStrength,
		Issues:       req.Issues,
		Strengths:    req.Strengths,
		Examples:     req.Examples,
		CreatedAt:    now,
		UpdatedAt:    now,
		CreatedBy:    userId,
		UpdatedBy:    userId,
	}

	// データベースに保存
	if err := r.db.Create(weaknessCategoryAnalysis).Error; err != nil {
		return nil, fmt.Errorf("failed to create weakness category analysis: %w", err)
	}

	return &model.CreateWeaknessCategoryAnalysisResponse{
		ID:           weaknessCategoryAnalysis.ID,
		AnalysisID:   weaknessCategoryAnalysis.AnalysisID,
		CategoryID:   weaknessCategoryAnalysis.CategoryID,
		CategoryName: weaknessCategoryAnalysis.CategoryName,
		Score:        weaknessCategoryAnalysis.Score,
		IsWeakness:   weaknessCategoryAnalysis.IsWeakness,
		IsStrength:   weaknessCategoryAnalysis.IsStrength,
		Issues:       weaknessCategoryAnalysis.Issues,
		Strengths:    weaknessCategoryAnalysis.Strengths,
		Examples:     weaknessCategoryAnalysis.Examples,
	}, nil
}

// GetWeaknessCategoryAnalysis カテゴリ別の分析結果を取得する
func (r *weaknessCategoryAnalysisRepository) GetWeaknessCategoryAnalysis(analysisId string) ([]*model.WeaknessCategoryAnalysisResponse, error) {
	var weaknessCategoryAnalyses []model.WeaknessCategoryAnalysisSummary

	// 正しいテーブル名を指定してクエリを実行
	if err := r.db.Table("weakness_category_analyses").Where("analysis_id = ?", analysisId).Find(&weaknessCategoryAnalyses).Error; err != nil {
		return nil, fmt.Errorf("failed to get weakness category analysis: %w", err)
	}

	// 複数のカテゴリ分析結果をレスポンス形式に変換
	var responses []*model.WeaknessCategoryAnalysisResponse
	for _, analysis := range weaknessCategoryAnalyses {
		// JSON文字列をスライスに変換
		issues, err := utils.ParseJSONStringArray(analysis.Issues)
		if err != nil {
			return nil, fmt.Errorf("failed to parse issues JSON: %w", err)
		}

		strengths, err := utils.ParseJSONStringArray(analysis.Strengths)
		if err != nil {
			return nil, fmt.Errorf("failed to parse strengths JSON: %w", err)
		}

		examples, err := utils.ParseJSONStringArray(analysis.Examples)
		if err != nil {
			return nil, fmt.Errorf("failed to parse examples JSON: %w", err)
		}

		response := &model.WeaknessCategoryAnalysisResponse{
			ID:           analysis.ID,
			AnalysisID:   analysis.AnalysisID,
			CategoryID:   analysis.CategoryID,
			CategoryName: analysis.CategoryName,
			Score:        analysis.Score,
			IsWeakness:   analysis.IsWeakness,
			IsStrength:   analysis.IsStrength,
			Issues:       issues,
			Strengths:    strengths,
			Examples:     examples,
		}
		responses = append(responses, response)
	}

	// 全ての結果を返す
	if len(responses) == 0 {
		return nil, fmt.Errorf("no weakness category analysis found for analysis_id: %s", analysisId)
	}

	return responses, nil
}

// UpdateWeaknessCategoryAnalysis カテゴリ別の分析結果を更新する
func (r *weaknessCategoryAnalysisRepository) UpdateWeaknessCategoryAnalysis(userId string, req *model.UpdateWeaknessCategoryAnalysisRequest) (*model.UpdateWeaknessCategoryAnalysisResponse, error) {
	now := time.Now()

	weaknessCategoryAnalysis := &model.WeaknessCategoryAnalysis{
		ID:           req.ID,
		AnalysisID:   req.AnalysisID,
		CategoryID:   req.CategoryID,
		CategoryName: req.CategoryName,
		Score:        req.Score,
		IsWeakness:   req.IsWeakness,
		IsStrength:   req.IsStrength,
		Issues:       req.Issues,
		Strengths:    req.Strengths,
		Examples:     req.Examples,
		UpdatedAt:    now,
		UpdatedBy:    userId,
	}

	// データベースに保存
	if err := r.db.Model(&model.WeaknessCategoryAnalysis{}).Where("id = ?", req.ID).Updates(weaknessCategoryAnalysis).Error; err != nil {
		return nil, fmt.Errorf("failed to update weakness category analysis: %w", err)
	}

	return &model.UpdateWeaknessCategoryAnalysisResponse{
		ID:           weaknessCategoryAnalysis.ID,
		AnalysisID:   weaknessCategoryAnalysis.AnalysisID,
		CategoryID:   weaknessCategoryAnalysis.CategoryID,
		CategoryName: weaknessCategoryAnalysis.CategoryName,
		Score:        weaknessCategoryAnalysis.Score,
		IsWeakness:   weaknessCategoryAnalysis.IsWeakness,
		IsStrength:   weaknessCategoryAnalysis.IsStrength,
		Issues:       weaknessCategoryAnalysis.Issues,
		Strengths:    weaknessCategoryAnalysis.Strengths,
		Examples:     weaknessCategoryAnalysis.Examples,
	}, nil
}
