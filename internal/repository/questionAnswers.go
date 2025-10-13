package repository

import (
	"fmt"
	"time"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionAnswersRepository interface {
	CreateQuestionAnswers(userID string, req *model.CreateQuestionAnswersRequest) (*model.CreateQuestionAnswersResponse, error)
	GetQuestionNoCorrectionAnswers(req *model.GetQuestionAnswersRequest) (*model.GetQuestionAnswersResponse, error)
	GetQuestionAnswerById(id string) (*model.QuestionAnswers, error)
}

type questionAnswersRepository struct {
	db *gorm.DB
}

func NewQuestionAnswersRepository(db *gorm.DB) QuestionAnswersRepository {
	return &questionAnswersRepository{db: db}
}

func (r *questionAnswersRepository) CreateQuestionAnswers(userID string, req *model.CreateQuestionAnswersRequest) (*model.CreateQuestionAnswersResponse, error) {
	now := time.Now()
	questionAnswer := &model.QuestionAnswers{
		ID:                       uuid.New().String(),
		UserID:                   userID,
		ProjectID:                req.ProjectID,
		QuestionTemplateMasterID: req.QuestionTemplateMasterID,
		UserAnswer:               req.UserAnswer,
		CreatedAt:                now,
		UpdatedAt:                now,
		CreatedBy:                userID,
		UpdatedBy:                userID,
	}

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(questionAnswer).Error; err != nil {
			return fmt.Errorf("回答データの作成に失敗しました: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &model.CreateQuestionAnswersResponse{
		ID:                       questionAnswer.ID,
		UserID:                   questionAnswer.UserID,
		ProjectID:                questionAnswer.ProjectID,
		QuestionTemplateMasterID: questionAnswer.QuestionTemplateMasterID,
		UserAnswer:               questionAnswer.UserAnswer,
		IsCorrection:             questionAnswer.IsCorrection,
	}, nil
}

func (r *questionAnswersRepository) GetQuestionNoCorrectionAnswers(req *model.GetQuestionAnswersRequest) (*model.GetQuestionAnswersResponse, error) {
	var questionAnswers []model.QuestionAnswers
	if err := r.db.Where("project_id = ? AND is_correction = ?", req.ProjectID, false).Find(&questionAnswers).Error; err != nil {
		return nil, fmt.Errorf("回答データの取得に失敗しました: %w", err)
	}

	return &model.GetQuestionAnswersResponse{
		QuestionAnswers: questionAnswers,
	}, nil
}

func (r *questionAnswersRepository) GetQuestionAnswerById(id string) (*model.QuestionAnswers, error) {
	var questionAnswer model.QuestionAnswers
	if err := r.db.Model(&model.QuestionAnswers{}).Where("id = ?", id).First(&questionAnswer).Error; err != nil {
		return nil, fmt.Errorf("回答データの取得に失敗しました: %w", err)
	}

	return &questionAnswer, nil
}
