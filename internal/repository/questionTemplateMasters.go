package repository

import (
	"fmt"

	"github.com/Takanpon2512/english-app/internal/model"
	"gorm.io/gorm"
)

type QuestionTemplateMastersRepository interface {
	GetQuestionTemplateMasters(req *model.GetQuestionTemplateMastersSearchRequest) (*model.GetQuestionTemplateMastersSearchResponse, error)
	GetQuestionTemplateMasterByID(id string) (*model.QuestionTemplateMastersSummary, error)
}

type questionTemplateMastersRepository struct {
	db *gorm.DB
}

func NewQuestionTemplateMastersRepository(db *gorm.DB) QuestionTemplateMastersRepository {
	return &questionTemplateMastersRepository{db: db}
}

func (r *questionTemplateMastersRepository) GetQuestionTemplateMasters(req *model.GetQuestionTemplateMastersSearchRequest) (*model.GetQuestionTemplateMastersSearchResponse, error) {
	var questionTemplateMasters []model.QuestionTemplateMasters
	var total int64

	query := r.db.Model(&model.QuestionTemplateMasters{}).Where("status = ?", "ACTIVE")

	if req.CategoryID != "" {
		query = query.Where("category_id = ?", req.CategoryID)
	}

	if req.QuestionType != "" {
		query = query.Where("question_type = ?", req.QuestionType)
	}

	if req.Level != "" {
		query = query.Where("level = ?", req.Level)
	}

	if req.EstimatedTime != 0 {
		query = query.Where("estimated_time = ?", req.EstimatedTime)
	}

	if req.Points != 0 {
		query = query.Where("points = ?", req.Points)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("質問テンプレートマスターの取得に失敗しました: %w", err)
	}

	offset := (req.Page - 1) * req.PerPage

	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(req.PerPage).
		Find(&questionTemplateMasters).Error; err != nil {
		return nil, fmt.Errorf("質問テンプレートマスターの取得に失敗しました: %w", err)
	}

	questionTemplateMastersSummary := make([]model.QuestionTemplateMastersSummary, len(questionTemplateMasters))
	for i, questionTemplateMaster := range questionTemplateMasters {
		questionTemplateMastersSummary[i] = model.QuestionTemplateMastersSummary{
			ID:            questionTemplateMaster.ID,
			CategoryID:    questionTemplateMaster.CategoryID,
			QuestionType:  questionTemplateMaster.QuestionType,
			English:       questionTemplateMaster.English,
			Japanese:      questionTemplateMaster.Japanese,
			Status:        questionTemplateMaster.Status,
			Level:         questionTemplateMaster.Level,
			EstimatedTime: questionTemplateMaster.EstimatedTime,
			Points:        questionTemplateMaster.Points,
		}
	}

	return &model.GetQuestionTemplateMastersSearchResponse{
		QuestionTemplateMasters: questionTemplateMastersSummary,
		Total:                   int(total),
		Page:                    req.Page,
		PerPage:                 req.PerPage,
	}, nil
}

func (r *questionTemplateMastersRepository) GetQuestionTemplateMasterByID(id string) (*model.QuestionTemplateMastersSummary, error) {
	var questionTemplateMaster model.QuestionTemplateMastersSummary

	if err := r.db.Model(&model.QuestionTemplateMasters{}).Where("id = ?", id).First(&questionTemplateMaster).Error; err != nil {
		return nil, fmt.Errorf("質問テンプレートマスターの取得に失敗しました: %w", err)
	}

	return &questionTemplateMaster, nil
}
