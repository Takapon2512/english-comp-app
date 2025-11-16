package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
)

type WeaknessCategoryAnalysisRepository interface {
	CreateWeaknessCategoryAnalysis(userId string, req *model.CreateWeaknessCategoryAnalysisRequest) (*model.CreateWeaknessCategoryAnalysisResponse, error)
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
