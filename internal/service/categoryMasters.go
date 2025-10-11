package service

import (
	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/repository"
)

type CategoryMastersService interface {
	GetCategoryMasters(userID string, req *model.GetCategoryMastersSearchRequest) (*model.GetCategoryMastersSearchResponse, error)
	GetCategoryMastersByID(userID string, id string) (*model.GetCategoryMastersByIDResponse, error)
}

type categoryMastersService struct {
	db   *gorm.DB
	repo repository.CategoryMastersRepository
}

func NewCategoryMastersService(db *gorm.DB, repo repository.CategoryMastersRepository) CategoryMastersService {
	return &categoryMastersService{db: db, repo: repo}
}

func (s *categoryMastersService) GetCategoryMasters(userID string, req *model.GetCategoryMastersSearchRequest) (*model.GetCategoryMastersSearchResponse, error) {
	return s.repo.GetCategoryMasters(req)
}

func (s *categoryMastersService) GetCategoryMastersByID(userID string, id string) (*model.GetCategoryMastersByIDResponse, error) {
	return s.repo.GetCategoryMastersByID(id)
}
