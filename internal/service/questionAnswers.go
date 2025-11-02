package service

import (
	"fmt"
	"math/rand"
	"slices"
	"time"

	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/repository"
)

type QuestionAnswersService interface {
	CreateQuestionAnswers(userID string, req *model.CreateQuestionAnswersRequest) (*model.CreateQuestionAnswersResponse, error)
	GetQuestionAnswersByProjectID(projectID string) (*model.GetQuestionAnswersResponse, error)
	UpdateQuestionAnswersFinish(projectID string) (*model.UpdateQuestionAnswersFinishResponse, error)
	GetProjectQuestionToAnswer(projectID string) (*model.GetProjectQuestionToAnswerResponse, error)
}

type questionAnswersService struct {
	db                          *gorm.DB
	repo                        repository.QuestionAnswersRepository
	projectQuestionsRepo        repository.ProjectQuestionsRepository
	questionTemplateMastersRepo repository.QuestionTemplateMastersRepository
}

func NewQuestionAnswersService(db *gorm.DB, repo repository.QuestionAnswersRepository, projectQuestionsRepo repository.ProjectQuestionsRepository, questionTemplateMastersRepo repository.QuestionTemplateMastersRepository) QuestionAnswersService {
	return &questionAnswersService{
		db:                          db,
		repo:                        repo,
		projectQuestionsRepo:        projectQuestionsRepo,
		questionTemplateMastersRepo: questionTemplateMastersRepo,
	}
}

// 解答作成
func (s *questionAnswersService) CreateQuestionAnswers(userID string, req *model.CreateQuestionAnswersRequest) (*model.CreateQuestionAnswersResponse, error) {
	return s.repo.CreateQuestionAnswers(userID, req)
}

// プロジェクトに紐づく解答を取得
func (s *questionAnswersService) GetQuestionAnswersByProjectID(projectID string) (*model.GetQuestionAnswersResponse, error) {
	return s.repo.GetQuestionAnswersByProjectID(projectID)
}

// 解答を完了にする
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

// プロジェクトに紐づいている問題のうち未回答の問題からランダムに1題取得
func (s *questionAnswersService) GetProjectQuestionToAnswer(projectID string) (*model.GetProjectQuestionToAnswerResponse, error) {
	// プロジェクトの解答を全て取得
	processingQuestionAnswers, err := s.repo.GetQuestionAnswersByProjectIDAndStatus(projectID, "PROCESSING")
	if err != nil {
		return nil, fmt.Errorf("解答の取得に失敗しました: %w", err)
	}

	nowQuestionNumber := len(processingQuestionAnswers) + 1

	// プロジェクトに紐づく問題を全て取得
	projectQuestions, err := s.projectQuestionsRepo.GetProjectQuestions(&model.GetProjectQuestionsRequest{ProjectID: projectID})
	if err != nil {
		return nil, fmt.Errorf("問題の取得に失敗しました: %w", err)
	}

	if len(processingQuestionAnswers) > 0 {
		// 解答済みの問題のテンプレートIDリストを作成
		answeredQuestionTemplateMasterIDs := make([]string, len(processingQuestionAnswers))
		for i, qa := range processingQuestionAnswers {
			answeredQuestionTemplateMasterIDs[i] = qa.QuestionTemplateMasterID
		}

		// 未解答の問題のテンプレートIDリストを作成
		var unansweredQuestionTemplateMasterIDs []string
		for _, pq := range projectQuestions.Questions {
			if !slices.Contains(answeredQuestionTemplateMasterIDs, pq.ID) {
				unansweredQuestionTemplateMasterIDs = append(unansweredQuestionTemplateMasterIDs, pq.ID)
			}
		}

		// 未解答の問題がない場合はnilを返す
		if len(unansweredQuestionTemplateMasterIDs) == 0 {
			return nil, nil
		}

		// 未解答の問題のうちランダムに1題取得
		randomQuestion := unansweredQuestionTemplateMasterIDs[rand.Intn(len(unansweredQuestionTemplateMasterIDs))]
		question, err := s.questionTemplateMastersRepo.GetQuestionTemplateMasterByID(randomQuestion)
		if err != nil {
			return nil, fmt.Errorf("問題の取得に失敗しました: %w", err)
		}

		return &model.GetProjectQuestionToAnswerResponse{
			Question: *question,
			NowQuestionNumber: nowQuestionNumber,
		}, nil
	} else {
		// プロジェクトに問題がない場合のチェック
		if len(projectQuestions.Questions) == 0 {
			return nil, fmt.Errorf("プロジェクトに問題が登録されていません")
		}

		// プロジェクトに紐づく問題の中から1題ランダムに取得
		randomQuestion := projectQuestions.Questions[rand.Intn(len(projectQuestions.Questions))]
		question, err := s.questionTemplateMastersRepo.GetQuestionTemplateMasterByID(randomQuestion.ID)
		if err != nil {
			return nil, fmt.Errorf("問題の取得に失敗しました: %w", err)
		}
		return &model.GetProjectQuestionToAnswerResponse{
			Question: *question,
			NowQuestionNumber: nowQuestionNumber,
		}, nil
	}
}
