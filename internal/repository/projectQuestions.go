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
	GetProjectQuestions(req *model.GetProjectQuestionsRequest) (*model.GetProjectQuestionsResponse, error)
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

func (r *projectQuestionsRepository) GetProjectQuestions(req *model.GetProjectQuestionsRequest) (*model.GetProjectQuestionsResponse, error) {
	// プロジェクトに紐づく質問テンプレートIDを取得
	var projectQuestions []model.ProjectQuestions
	if err := r.db.Where("project_id = ? AND deleted_at IS NULL", req.ProjectID).Find(&projectQuestions).Error; err != nil {
		return nil, fmt.Errorf("プロジェクト質問の取得に失敗しました: %w", err)
	}

	// 質問テンプレートIDのリストを作成
	var questionTemplateIDs []string
	for _, pq := range projectQuestions {
		questionTemplateIDs = append(questionTemplateIDs, pq.QuestionTemplateMasterID)
	}

	// 質問テンプレートの詳細情報を取得
	var questions []model.QuestionTemplateMastersSummary
	if len(questionTemplateIDs) > 0 {
		if err := r.db.Table("question_template_masters").
			Select("id, category_id, question_type, english, japanese, status, level, estimated_time, points").
			Where("id IN ? AND deleted_at IS NULL", questionTemplateIDs).
			Find(&questions).Error; err != nil {
			return nil, fmt.Errorf("質問テンプレートの取得に失敗しました: %w", err)
		}
	}

	return &model.GetProjectQuestionsResponse{
		Questions: questions,
	}, nil
}
