package service

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/repository"
)

type WeaknessAnalysisService interface {
	CreateWeaknessAnalysis(userId string, req *model.CreateWeaknessAnalysisRequest) (*model.CreateWeaknessAnalysisResponse, error)

	// LLMによる分析処理
	WeaknessCategoryAnalysis(userId string, projectId string) ([]model.LLMWeaknessCategoryAnalysisRequest, error)
}

type weaknessAnalysisService struct {
	db                          *gorm.DB
	repo                        repository.WeaknessAnalysisRepository
	correctResultsRepo          repository.CorrectResultsRepository
	questionAnswersRepo         repository.QuestionAnswersRepository
	questionTemplateMastersRepo repository.QuestionTemplateMastersRepository
	categoryMastersRepo         repository.CategoryMastersRepository
}

func NewWeaknessAnalysisService(
	db *gorm.DB,
	repo repository.WeaknessAnalysisRepository,
	correctResultsRepo repository.CorrectResultsRepository,
	questionAnswersRepo repository.QuestionAnswersRepository,
	questionTemplateMastersRepo repository.QuestionTemplateMastersRepository,
	categoryMastersRepo repository.CategoryMastersRepository,
) WeaknessAnalysisService {
	return &weaknessAnalysisService{
		db:                          db,
		repo:                        repo,
		correctResultsRepo:          correctResultsRepo,
		questionAnswersRepo:         questionAnswersRepo,
		questionTemplateMastersRepo: questionTemplateMastersRepo,
		categoryMastersRepo:         categoryMastersRepo,
	}
}

// CreateWeaknessAnalysis 学習弱点分析を作成する
// この時、学習カテゴリ分析、詳細分析、学習アドバイスを作成する
func (s *weaknessAnalysisService) CreateWeaknessAnalysis(userId string, req *model.CreateWeaknessAnalysisRequest) (*model.CreateWeaknessAnalysisResponse, error) {
	weaknessAnalysis, err := s.repo.CreateWeaknessAnalysis(userId, req)

	if err != nil {
		return nil, err
	}

	// 学習カテゴリ分析を作成する
  llmRequests, err := s.WeaknessCategoryAnalysis(userId, req.ProjectID)
  if err != nil {
    return nil, err
  }
  fmt.Println(llmRequests)

	// 詳細分析を作成する

	// 学習アドバイスを作成する

	return weaknessAnalysis, nil
}

// UpdateWeaknessAnalysis 学習弱点分析を更新する
func (s *weaknessAnalysisService) UpdateWeaknessAnalysis(userId string, req *model.UpdateWeaknessAnalysisRequest) (*model.UpdateWeaknessAnalysisResponse, error) {
	return s.repo.UpdateWeaknessAnalysis(userId, req)
}

// weaknessCategoryの分析をLLMにて行う
func (s *weaknessAnalysisService) WeaknessCategoryAnalysis(userId string, projectId string) ([]model.LLMWeaknessCategoryAnalysisRequest, error) {
	// 解答データを取得
	correctResults, err := s.correctResultsRepo.GetCorrectResults(&model.GetCorrectResultsRequest{ProjectID: projectId})
	if err != nil {
		return nil, err
	}

	llmRequests := []model.LLMWeaknessCategoryAnalysisRequest{}

	for _, correctResult := range correctResults.CorrectResults {
		// 解答データを取得
		questionAnswer, err := s.questionAnswersRepo.GetQuestionAnswerById(correctResult.QuestionAnswerID)
		if err != nil {
			return nil, err
		}

		if questionAnswer == nil {
			return nil, fmt.Errorf("解答データ（ID: %s）が見つかりません", correctResult.QuestionAnswerID)
		}

		// 問題データ取得
		questionTemplateMaster, err := s.questionTemplateMastersRepo.GetQuestionTemplateMasterByID(correctResult.QuestionAnswer.QuestionTemplateMasterID)
		if err != nil {
			return nil, err
		}

		if questionTemplateMaster == nil {
			return nil, fmt.Errorf("問題データ（ID: %s）が見つかりません", correctResult.QuestionAnswer.QuestionTemplateMasterID)
		}

		// カテゴリデータを取得
		categoryMaster, err := s.categoryMastersRepo.GetCategoryMastersByID(questionTemplateMaster.CategoryID)
		if err != nil {
			return nil, err
		}

		if categoryMaster == nil {
			return nil, fmt.Errorf("カテゴリデータ（ID: %s）が見つかりません", questionTemplateMaster.CategoryID)
		}

		// LLMによる分析を行う
		llmRequest := model.LLMWeaknessCategoryAnalysisRequest{
			CategoryName:  categoryMaster.CategoryMasters.Name,
			Question:      questionTemplateMaster.English,
			UserAnswer:    questionAnswer.UserAnswer,
			CorrectAnswer: correctResult.ExampleCorrection,
		}

		llmRequests = append(llmRequests, llmRequest)
	}

	return llmRequests, nil
}
