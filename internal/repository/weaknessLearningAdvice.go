package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
)

type WeaknessLearningAdviceRepository interface {
	CreateWeaknessLearningAdvice(userId string, req *model.CreateWeaknessLearningAdviceRequest) (*model.CreateWeaknessLearningAdviceResponse, error)
	GetWeaknessLearningAdvice(analysisId string) (*model.WeaknessLearningAdviceSummary, error)
	UpdateWeaknessLearningAdvice(userId string, req *model.UpdateWeaknessLearningAdviceRequest) (*model.UpdateWeaknessLearningAdviceResponse, error)
}

type weaknessLearningAdviceRepository struct {
	db *gorm.DB
}

func NewWeaknessLearningAdviceRepository(db *gorm.DB) WeaknessLearningAdviceRepository {
	return &weaknessLearningAdviceRepository{db: db}
}

// CreateWeaknessLearningAdvice 学習アドバイスを作成する
func (r *weaknessLearningAdviceRepository) CreateWeaknessLearningAdvice(userId string, req *model.CreateWeaknessLearningAdviceRequest) (*model.CreateWeaknessLearningAdviceResponse, error) {
	now := time.Now()

	weaknessLearningAdvice := &model.WeaknessLearningAdvice{
		ID:                  uuid.New().String(),
		AnalysisID:          req.AnalysisID,
		LearningAdvice:      req.LearningAdvice,
		RecommendedActions:  req.RecommendedActions,
		NextGoals:           req.NextGoals,
		StudyPlan:           req.StudyPlan,
		MotivationalMessage: req.MotivationalMessage,
		CreatedAt:           now,
		UpdatedAt:           now,
		CreatedBy:           userId,
		UpdatedBy:           userId,
	}

	if err := r.db.Create(weaknessLearningAdvice).Error; err != nil {
		return nil, fmt.Errorf("failed to create weakness learning advice: %w", err)
	}

	return &model.CreateWeaknessLearningAdviceResponse{
		ID:                  weaknessLearningAdvice.ID,
		AnalysisID:          weaknessLearningAdvice.AnalysisID,
		LearningAdvice:      weaknessLearningAdvice.LearningAdvice,
		RecommendedActions:  weaknessLearningAdvice.RecommendedActions,
		NextGoals:           weaknessLearningAdvice.NextGoals,
		StudyPlan:           weaknessLearningAdvice.StudyPlan,
		MotivationalMessage: weaknessLearningAdvice.MotivationalMessage,
	}, nil
}

// GetWeaknessLearningAdvice 学習アドバイスを取得する
func (r *weaknessLearningAdviceRepository) GetWeaknessLearningAdvice(analysisId string) (*model.WeaknessLearningAdviceSummary, error) {
	var weaknessLearningAdvice model.WeaknessLearningAdviceSummary

	// 正しいテーブル名を指定してクエリを実行
	if err := r.db.Table("weakness_learning_advice").Where("analysis_id = ?", analysisId).First(&weaknessLearningAdvice).Error; err != nil {
		return nil, fmt.Errorf("failed to get weakness learning advice: %w", err)
	}

	return &model.WeaknessLearningAdviceSummary{
		ID:                  weaknessLearningAdvice.ID,
		AnalysisID:          weaknessLearningAdvice.AnalysisID,
		LearningAdvice:      weaknessLearningAdvice.LearningAdvice,
		RecommendedActions:  weaknessLearningAdvice.RecommendedActions,
		NextGoals:           weaknessLearningAdvice.NextGoals,
		StudyPlan:           weaknessLearningAdvice.StudyPlan,
		MotivationalMessage: weaknessLearningAdvice.MotivationalMessage,
	}, nil
}

// UpdateWeaknessLearningAdvice 学習アドバイスを更新する
func (r *weaknessLearningAdviceRepository) UpdateWeaknessLearningAdvice(userId string, req *model.UpdateWeaknessLearningAdviceRequest) (*model.UpdateWeaknessLearningAdviceResponse, error) {
	now := time.Now()

	weaknessLearningAdvice := &model.WeaknessLearningAdvice{
		ID:                    req.ID,
		AnalysisID:            req.AnalysisID,
		LearningAdvice:      req.LearningAdvice,
		RecommendedActions:  req.RecommendedActions,
		NextGoals:           req.NextGoals,
		StudyPlan:           req.StudyPlan,
		MotivationalMessage: req.MotivationalMessage,
		UpdatedAt:           now,
		UpdatedBy:           userId,
	}

	if err := r.db.Save(weaknessLearningAdvice).Error; err != nil {
		return nil, fmt.Errorf("failed to update weakness learning advice: %w", err)
	}

	return &model.UpdateWeaknessLearningAdviceResponse{
		ID:                  weaknessLearningAdvice.ID,
		AnalysisID:          weaknessLearningAdvice.AnalysisID,
		LearningAdvice:      weaknessLearningAdvice.LearningAdvice,
		RecommendedActions:  weaknessLearningAdvice.RecommendedActions,
		NextGoals:           weaknessLearningAdvice.NextGoals,
		StudyPlan:           weaknessLearningAdvice.StudyPlan,
		MotivationalMessage: weaknessLearningAdvice.MotivationalMessage,
	}, nil
}