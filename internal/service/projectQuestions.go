package service

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/repository"
)

type ProjectQuestionsService interface {
	CreateProjectQuestions(userID string, req *model.CreateProjectQuestionsRequest) (*model.CreateProjectQuestionsResponse, error)
}

type projectQuestionsService struct {
	db                          *gorm.DB
	repo                        repository.ProjectQuestionsRepository
	questionTemplateMastersRepo repository.QuestionTemplateMastersRepository
}

func NewProjectQuestionsService(db *gorm.DB, repo repository.ProjectQuestionsRepository, questionTemplateMastersRepo repository.QuestionTemplateMastersRepository) ProjectQuestionsService {
	return &projectQuestionsService{
		db:                          db,
		repo:                        repo,
		questionTemplateMastersRepo: questionTemplateMastersRepo,
	}
}

func (s *projectQuestionsService) CreateProjectQuestions(userID string, req *model.CreateProjectQuestionsRequest) (*model.CreateProjectQuestionsResponse, error) {
	// ユーザーIDをリクエストに設定
	req.UserID = userID

	// プロジェクト質問を作成
	response, err := s.repo.CreateProjectQuestions(req)
	if err != nil {
		return nil, err
	}

	// 質問テンプレートマスターの詳細情報を取得してレスポンスに追加
	for i, projectQuestion := range response.ProjectQuestions {
		master, err := s.questionTemplateMastersRepo.GetQuestionTemplateMasterByID(projectQuestion.QuestionTemplateMasterID)
		if err != nil {
			return nil, fmt.Errorf("質問テンプレートマスターの取得に失敗しました: %w", err)
		}
		response.ProjectQuestions[i].Questions = *master
	}

	return response, nil
}
