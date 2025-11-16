package service

import (
	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/repository"
)

type WeaknessAnalysisService interface {
  CreateWeaknessAnalysis(userId string, req *model.CreateWeaknessAnalysisRequest) (*model.CreateWeaknessAnalysisResponse)
}

type weaknessAnalysisService struct {
  db *gorm.DB
  repo repository.WeaknessAnalysisRepository
}

func NewWeaknessAnalysisService (db *gorm.DB, repo repository.WeaknessAnalysisRepository) WeaknessAnalysisService {
  return &weaknessAnalysisService{
    db: db,
    repo: repo,
  }
}

// CreateWeaknessAnalysis 学習弱点分析を作成する
func (s *weaknessAnalysisService) CreateWeaknessAnalysis(userId string, req *model.CreateWeaknessAnalysisRequest) (*model.CreateWeaknessAnalysisResponse) {
  return s.CreateWeaknessAnalysis(userId, req)
}
