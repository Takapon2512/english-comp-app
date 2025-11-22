package service

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/repository"
)

type WeaknessAnalysisService interface {
	CreateWeaknessAnalysis(userId string, req *model.CreateWeaknessAnalysisRequest) (*model.CreateWeaknessAnalysisResponse, error)
	GetWeaknessAnalysis(userId string, req *model.GetWeaknessAnalysisRequest) (*model.GetWeaknessAnalysisResponse, error)

	// LLMによる分析処理
	WeaknessCategoryAnalysis(userId string, projectId string) ([]model.LLMWeaknessCategoryAnalysisRequest, error)
	WeaknessDetailedAnalysis(userId string, projectId string) ([]model.LLMWeaknessDetailedAnalysisRequest, error)
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
	// 作成前に同じプロジェクトの分析が存在するか確認
	existingAnalysis, err := s.repo.GetWeaknessAnalysis(userId, &model.GetWeaknessAnalysisRequest{ProjectID: req.ProjectID})
	if err != nil {
		return nil, err
	}
	if existingAnalysis != nil {
		return nil, fmt.Errorf("同じプロジェクトの分析が既に存在します")
	}

	weaknessAnalysis, err := s.repo.CreateWeaknessAnalysis(userId, req)

	if err != nil {
		return nil, err
	}

	// 学習カテゴリ分析を作成する
	llmRequestsCategoryAnalysis, err := s.WeaknessCategoryAnalysis(userId, req.ProjectID)
	if err != nil {
		return nil, err
	}
	fmt.Println(llmRequestsCategoryAnalysis)

	// 詳細分析を作成する
	llmRequestsDetailedAnalysis, err := s.WeaknessDetailedAnalysis(userId, req.ProjectID)
	if err != nil {
		return nil, err
	}
	fmt.Println(llmRequestsDetailedAnalysis)

	// 学習アドバイスを作成する

	return weaknessAnalysis, nil
}

// UpdateWeaknessAnalysis 学習弱点分析を更新する
func (s *weaknessAnalysisService) UpdateWeaknessAnalysis(userId string, req *model.UpdateWeaknessAnalysisRequest) (*model.UpdateWeaknessAnalysisResponse, error) {
	// 更新前に該当のプロジェクトが存在することを確認
	existingAnalysis, err := s.repo.GetWeaknessAnalysis(userId, &model.GetWeaknessAnalysisRequest{ProjectID: req.ProjectID})

	if err != nil {
		return nil, err
	}

	if existingAnalysis == nil {
		return nil, fmt.Errorf("プロジェクトが存在しません")
	}
}

// GetWeaknessAnalysis 学習弱点分析を取得する
func (s *weaknessAnalysisService) GetWeaknessAnalysis(userId string, req *model.GetWeaknessAnalysisRequest) (*model.GetWeaknessAnalysisResponse, error) {
	return s.repo.GetWeaknessAnalysis(userId, req)
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
		questionTemplateMaster, err := s.questionTemplateMastersRepo.GetQuestionTemplateMasterByID(questionAnswer.QuestionTemplateMasterID)
		if err != nil {
			return nil, err
		}

		if questionTemplateMaster == nil {
			return nil, fmt.Errorf("問題データ（ID: %s）が見つかりません", questionAnswer.QuestionTemplateMasterID)
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

// WeaknessDetailedAnalysis 詳細分析をLLMにて行う
func (s *weaknessAnalysisService) WeaknessDetailedAnalysis(userId string, projectId string) ([]model.LLMWeaknessDetailedAnalysisRequest, error) {
	// 解答データを取得
	correctResults, err := s.correctResultsRepo.GetCorrectResults(&model.GetCorrectResultsRequest{ProjectID: projectId})
	if err != nil {
		return nil, err
	}

	llmRequests := []model.LLMWeaknessDetailedAnalysisRequest{}

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
		questionTemplateMaster, err := s.questionTemplateMastersRepo.GetQuestionTemplateMasterByID(questionAnswer.QuestionTemplateMasterID)
		if err != nil {
			return nil, err
		}

		if questionTemplateMaster == nil {
			return nil, fmt.Errorf("問題データ（ID: %s）が見つかりません", questionAnswer.QuestionTemplateMasterID)
		}

		// LLMによる分析を行う
		llmRequest := model.LLMWeaknessDetailedAnalysisRequest{
			Question:      questionTemplateMaster.English,
			UserAnswer:    questionAnswer.UserAnswer,
			CorrectAnswer: correctResult.ExampleCorrection,
		}

		llmRequests = append(llmRequests, llmRequest)
	}

	return llmRequests, nil
}
