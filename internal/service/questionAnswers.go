package service

import (
	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/repository"
)

type QuestionAnswersService interface {
	CreateQuestionAnswers(userID string, req *model.CreateQuestionAnswersRequest) (*model.CreateQuestionAnswersResponse, error)
	GetQuestionNoCorrectionAnswers(req *model.GetQuestionAnswersRequest) (*model.GetQuestionAnswersResponse, error)
}

type questionAnswersService struct {
	db   *gorm.DB
	repo repository.QuestionAnswersRepository
}

func NewQuestionAnswersService(db *gorm.DB, repo repository.QuestionAnswersRepository) QuestionAnswersService {
	return &questionAnswersService{db: db, repo: repo}
}

func (s *questionAnswersService) CreateQuestionAnswers(userID string, req *model.CreateQuestionAnswersRequest) (*model.CreateQuestionAnswersResponse, error) {
	return s.repo.CreateQuestionAnswers(userID, req)
}

func (s *questionAnswersService) GetQuestionNoCorrectionAnswers(req *model.GetQuestionAnswersRequest) (*model.GetQuestionAnswersResponse, error) {
	return s.repo.GetQuestionNoCorrectionAnswers(req)
}