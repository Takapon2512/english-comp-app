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
	GetQuestionAnswersByProjectID(projectID string) (*model.GetQuestionAnswersResponse, error)
	GetQuestionAnswersByProjectIDAndStatus(projectID, status string) ([]model.QuestionAnswers, error)
	UpdateQuestionAnswer(tx *gorm.DB, questionAnswer *model.QuestionAnswers) error
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
		ChallengeCount:           req.ChallengeCount,
		Status:                   "PROCESSING",
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
		ChallengeCount:           questionAnswer.ChallengeCount,
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

func (r *questionAnswersRepository) GetQuestionAnswersByProjectID(projectID string) (*model.GetQuestionAnswersResponse, error) {
	var questionAnswers []model.QuestionAnswers
	if err := r.db.Where("project_id = ?", projectID).Find(&questionAnswers).Error; err != nil {
		return nil, fmt.Errorf("回答データの取得に失敗しました: %w", err)
	}

	return &model.GetQuestionAnswersResponse{
		QuestionAnswers: questionAnswers,
	}, nil
}

// GetQuestionAnswersByProjectIDAndStatus プロジェクトIDとステータスで回答データを取得
func (r *questionAnswersRepository) GetQuestionAnswersByProjectIDAndStatus(projectID, status string) ([]model.QuestionAnswers, error) {
	var questionAnswers []model.QuestionAnswers
	if err := r.db.Where("project_id = ? AND status = ?", projectID, status).Find(&questionAnswers).Error; err != nil {
		return nil, fmt.Errorf("回答データの取得に失敗しました: %w", err)
	}
	return questionAnswers, nil
}

// UpdateQuestionAnswer 回答データを更新（トランザクション対応）
func (r *questionAnswersRepository) UpdateQuestionAnswer(tx *gorm.DB, questionAnswer *model.QuestionAnswers) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	if err := db.Save(questionAnswer).Error; err != nil {
		return fmt.Errorf("回答データの更新に失敗しました: %w", err)
	}
	return nil
}
