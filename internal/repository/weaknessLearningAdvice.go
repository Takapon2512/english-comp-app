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
}

type weaknessLearningAdviceRepository struct {
  db * gorm.DB
}

func NewWeaknessLearningAdviceRepository(db *gorm.DB) WeaknessLearningAdviceRepository {
  return &weaknessLearningAdviceRepository{db: db}
}

// CreateWeaknessLearningAdvice 学習アドバイスを作成する
func (r *weaknessLearningAdviceRepository) CreateWeaknessLearningAdvice(userId string, req *model.CreateWeaknessLearningAdviceRequest) (*model.CreateWeaknessLearningAdviceResponse, error) {
  now := time.Now()

  weaknessLearningAdvice := &model.WeaknessLearningAdvice{
    ID: uuid.New().String(),
    AnalysisID: req.AnalysisID,
    LearningAdvice: req.LearningAdvice,
    RecommendedActions: req.RecommendedActions,
    NextGoals: req.NextGoals,
    StudyPlan: req.StudyPlan,
    MotivationalMessage: req.MotivationalMessage,
    CreatedAt: now,
    UpdatedAt: now,
    CreatedBy: userId,
    UpdatedBy: userId,
  }

  if err := r.db.Create(weaknessLearningAdvice).Error; err != nil {
    return nil, fmt.Errorf("failed to create weakness learning advice: %w", err)
  }

  return &model.CreateWeaknessLearningAdviceResponse{
    ID: weaknessLearningAdvice.ID,
    AnalysisID: weaknessLearningAdvice.AnalysisID,
    LearningAdvice: weaknessLearningAdvice.LearningAdvice,
    RecommendedActions: weaknessLearningAdvice.RecommendedActions,
    NextGoals: weaknessLearningAdvice.NextGoals,
    StudyPlan: weaknessLearningAdvice.StudyPlan,
    MotivationalMessage: weaknessLearningAdvice.MotivationalMessage,
  }, nil
}