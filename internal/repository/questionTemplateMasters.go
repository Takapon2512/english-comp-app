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

	if req.MinPoints != 0 {
		query = query.Where("points >= ?", req.MinPoints)
	}

	if req.MaxPoints != 0 {
		query = query.Where("points <= ?", req.MaxPoints)
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
		categoryInfo, err := r.GetCategoryInfo(questionTemplateMaster.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("カテゴリーの取得に失敗しました: %w", err)
		}

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
			Category:      *categoryInfo,
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
	var questionTemplateMaster model.QuestionTemplateMasters

	if err := r.db.Model(&model.QuestionTemplateMasters{}).Where("id = ?", id).First(&questionTemplateMaster).Error; err != nil {
		return nil, fmt.Errorf("質問テンプレートマスターの取得に失敗しました: %w", err)
	}

	// カテゴリ情報を取得
	categoryInfo, err := r.GetCategoryInfo(questionTemplateMaster.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("カテゴリーの取得に失敗しました: %w", err)
	}

	// サマリーを作成
	summary := &model.QuestionTemplateMastersSummary{
		ID:            questionTemplateMaster.ID,
		CategoryID:    questionTemplateMaster.CategoryID,
		QuestionType:  questionTemplateMaster.QuestionType,
		English:       questionTemplateMaster.English,
		Japanese:      questionTemplateMaster.Japanese,
		Status:        questionTemplateMaster.Status,
		Level:         questionTemplateMaster.Level,
		EstimatedTime: questionTemplateMaster.EstimatedTime,
		Points:        questionTemplateMaster.Points,
		Category:      *categoryInfo,
	}

	return summary, nil
}

func (r *questionTemplateMastersRepository) GetCategoryInfo(categoryID string) (*model.CategoryInfo, error) {
	var categoryInfo model.CategoryInfo

	if err := r.db.Model(&model.CategoryMasters{}).Where("id = ?", categoryID).First(&categoryInfo).Error; err != nil {
		return nil, fmt.Errorf("カテゴリーの取得に失敗しました: %w", err)
	}

	return &model.CategoryInfo{
		ID:   categoryInfo.ID,
		Name: categoryInfo.Name,
	}, nil
}
