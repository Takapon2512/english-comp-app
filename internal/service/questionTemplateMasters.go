package service

import (
	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/repository"
)

type QuestionTemplateMastersService interface {
	GetQuestionTemplateMasters(userID string, req *model.GetQuestionTemplateMastersSearchRequest) (*model.GetQuestionTemplateMastersSearchResponse, error)
	GetQuestionTemplateMasterByID(userID string, id string) (*model.QuestionTemplateMastersSummary, error)
}

type questionTemplateMastersService struct {
	db   *gorm.DB
	repo repository.QuestionTemplateMastersRepository
}

func NewQuestionTemplateMastersService(db *gorm.DB, repo repository.QuestionTemplateMastersRepository) QuestionTemplateMastersService {
	return &questionTemplateMastersService{db: db, repo: repo}
}

func (s *questionTemplateMastersService) GetQuestionTemplateMasters(userID string, req *model.GetQuestionTemplateMastersSearchRequest) (*model.GetQuestionTemplateMastersSearchResponse, error) {
	return s.repo.GetQuestionTemplateMasters(req)
}

func (s *questionTemplateMastersService) GetQuestionTemplateMasterByID(userID string, id string) (*model.QuestionTemplateMastersSummary, error) {
	return s.repo.GetQuestionTemplateMasterByID(id)
}