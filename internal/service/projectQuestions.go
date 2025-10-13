package service

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/repository"
)

type ProjectQuestionsService interface {
	CreateProjectQuestions(userID string, req *model.CreateProjectQuestionsRequest) (*model.CreateProjectQuestionsResponse, error)
	GetProjectQuestions(userID string, req *model.GetProjectQuestionsRequest) (*model.GetProjectQuestionsResponse, error)
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

func (s *projectQuestionsService) GetProjectQuestions(userID string, req *model.GetProjectQuestionsRequest) (*model.GetProjectQuestionsResponse, error) {
	req.UserID = userID

	// プロジェクト質問を取得
	response, err := s.repo.GetProjectQuestions(req)
	if err != nil {
		return nil, err
	}

	// カテゴリ情報を取得してマッピング
	if len(response.Questions) > 0 {
		// カテゴリIDのリストを作成
		var categoryIDs []string
		for _, question := range response.Questions {
			categoryIDs = append(categoryIDs, question.CategoryID)
		}

		// カテゴリ情報を取得
		var categories []model.CategoryInfo
		if err := s.db.Table("category_masters").
			Select("id, name").
			Where("id IN ? AND deleted_at IS NULL", categoryIDs).
			Find(&categories).Error; err != nil {
			return nil, fmt.Errorf("カテゴリ情報の取得に失敗しました: %w", err)
		}

		// カテゴリ情報をマップに変換
		categoryMap := make(map[string]model.CategoryInfo)
		for _, category := range categories {
			categoryMap[category.ID] = category
		}

		// 各質問にカテゴリ情報をマッピング
		for i, question := range response.Questions {
			response.Questions[i].Category = categoryMap[question.CategoryID]
		}
	}

	return response, nil
}
