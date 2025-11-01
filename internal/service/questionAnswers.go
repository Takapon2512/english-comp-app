package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/repository"
)

type QuestionAnswersService interface {
	CreateQuestionAnswers(userID string, req *model.CreateQuestionAnswersRequest) (*model.CreateQuestionAnswersResponse, error)
	GetQuestionNoCorrectionAnswers(req *model.GetQuestionAnswersRequest) (*model.GetQuestionAnswersResponse, error)
	GetQuestionAnswersByProjectID(projectID string) (*model.GetQuestionAnswersResponse, error)
	UpdateQuestionAnswersFinish(projectID string) (*model.UpdateQuestionAnswersFinishResponse, error)
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

func (s *questionAnswersService) GetQuestionAnswersByProjectID(projectID string) (*model.GetQuestionAnswersResponse, error) {
	return s.repo.GetQuestionAnswersByProjectID(projectID)
}

func (s *questionAnswersService) UpdateQuestionAnswersFinish(projectID string) (*model.UpdateQuestionAnswersFinishResponse, error) {
	// プロジェクトのPROCESSING状態の回答データを取得
	questionAnswers, err := s.repo.GetQuestionAnswersByProjectIDAndStatus(projectID, "PROCESSING")
	if err != nil {
		return nil, fmt.Errorf("回答データの取得に失敗しました: %w", err)
	}

	now := time.Now()
	var updatedQuestionAnswers []model.QuestionAnswers

	// トランザクション内でビジネスロジックを実行
	err = s.db.Transaction(func(tx *gorm.DB) error {
		for _, questionAnswer := range questionAnswers {
			
			// ChallengeCountが0より大きい場合は更新日時のみ更新
			if questionAnswer.ChallengeCount > 0 {
				questionAnswer.UpdatedAt = now
				questionAnswer.Status = "FINISHED"

				if err := s.repo.UpdateQuestionAnswer(tx, &questionAnswer); err != nil {
					return fmt.Errorf("回答データの更新に失敗しました: %w", err)
				}
			} else {
				// ChallengeCountが0の場合はステータスをFINISHEDに変更
				questionAnswer.Status = "FINISHED"
				questionAnswer.UpdatedAt = now
				if err := s.repo.UpdateQuestionAnswer(tx, &questionAnswer); err != nil {
					return fmt.Errorf("回答データの更新に失敗しました: %w", err)
				}
			}
			updatedQuestionAnswers = append(updatedQuestionAnswers, questionAnswer)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &model.UpdateQuestionAnswersFinishResponse{
		QuestionAnswers: updatedQuestionAnswers,
	}, nil
}
