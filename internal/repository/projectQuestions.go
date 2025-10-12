package repository

import (
	"fmt"
	"time"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectQuestionsRepository interface {
	CreateProjectQuestions(req *model.CreateProjectQuestionsRequest) (*model.CreateProjectQuestionsResponse, error)
}

type projectQuestionsRepository struct {
	db *gorm.DB
}

func NewProjectQuestionsRepository(db *gorm.DB) ProjectQuestionsRepository {
	return &projectQuestionsRepository{db: db}
}

func (r *projectQuestionsRepository) CreateProjectQuestions(req *model.CreateProjectQuestionsRequest) (*model.CreateProjectQuestionsResponse, error) {
	now := time.Now()
	var createdProjectQuestions []model.ProjectQuestions

	for _, questionTemplateMasterID := range req.QuestionTemplateMasterIDs {
		projectQuestion := &model.ProjectQuestions{
			ID:                       uuid.New().String(),
			ProjectID:                req.ProjectID,
			QuestionTemplateMasterID: questionTemplateMasterID,
			CreatedAt:                now,
			UpdatedAt:                now,
			CreatedBy:                req.UserID,
			UpdatedBy:                req.UserID,
		}

		if err := r.db.Create(projectQuestion).Error; err != nil {
			return nil, fmt.Errorf("プロジェクト質問の作成に失敗しました: %w", err)
		}
		createdProjectQuestions = append(createdProjectQuestions, *projectQuestion)
	}

	// 作成されたプロジェクト質問のサマリーを作成
	projectQuestionsSummary := make([]model.ProjectQuestionsSummary, len(createdProjectQuestions))
	for i, projectQuestion := range createdProjectQuestions {
		projectQuestionsSummary[i] = model.ProjectQuestionsSummary{
			ID:                       projectQuestion.ID,
			ProjectID:                projectQuestion.ProjectID,
			QuestionTemplateMasterID: projectQuestion.QuestionTemplateMasterID,
		}
	}

	return &model.CreateProjectQuestionsResponse{
		ProjectQuestions: projectQuestionsSummary,
	}, nil
}
