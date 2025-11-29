package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/config"
	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/repository"
	"github.com/Takanpon2512/english-app/internal/utils"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// LLMからのカテゴリ分析レスポンス用構造体
type CategoryAnalysisResult struct {
	IsWeakness bool     `json:"is_weakness"`
	IsStrength bool     `json:"is_strength"`
	Issues     []string `json:"issues"`
	Strengths  []string `json:"strengths"`
	Examples   []string `json:"examples"`
}

type WeaknessAnalysisService interface {
	CreateWeaknessAnalysis(userId string, req *model.CreateWeaknessAnalysisRequest) (*model.CreateWeaknessAnalysisResponse, error)
	GetWeaknessAnalysis(userId string, req *model.GetWeaknessAnalysisRequest) (*model.GetWeaknessAnalysisResponse, error)

	// LLMによる分析処理
	WeaknessCategoryAnalysis(userId string, projectId string) (map[string]*CategoryAnalysisResult, error)
	WeaknessDetailedAnalysis(userId string, projectId string) (*model.DetailedAnalysisResult, error)
	WeaknessLearningAdvice(userId string, projectId string, detailedAnalysis *model.DetailedAnalysisResult) (*model.PersonalizedAdvice, error)

	GetWeaknessAnalysisAllSummary(userId string, projectId string) (*model.WeaknessAnalysisAllSummary, error)
}

type weaknessAnalysisService struct {
	db                           *gorm.DB
	repo                         repository.WeaknessAnalysisRepository
	correctResultsRepo           repository.CorrectResultsRepository
	questionAnswersRepo          repository.QuestionAnswersRepository
	questionTemplateMastersRepo  repository.QuestionTemplateMastersRepository
	categoryMastersRepo          repository.CategoryMastersRepository
	weaknessCategoryAnalysisRepo repository.WeaknessCategoryAnalysisRepository
	weaknessDetailedAnalysisRepo repository.WeaknessDetailedAnalysisRepository
	weaknessLearningAdviceRepo   repository.WeaknessLearningAdviceRepository
	claudeClient                 anthropic.Client
	prompts                      *config.WeaknessAnalysisPrompts
}

func NewWeaknessAnalysisService(
	db *gorm.DB,
	repo repository.WeaknessAnalysisRepository,
	correctResultsRepo repository.CorrectResultsRepository,
	questionAnswersRepo repository.QuestionAnswersRepository,
	questionTemplateMastersRepo repository.QuestionTemplateMastersRepository,
	categoryMastersRepo repository.CategoryMastersRepository,
	weaknessCategoryAnalysisRepo repository.WeaknessCategoryAnalysisRepository,
	weaknessDetailedAnalysisRepo repository.WeaknessDetailedAnalysisRepository,
	weaknessLearningAdviceRepo repository.WeaknessLearningAdviceRepository,
) WeaknessAnalysisService {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		log.Fatal("CLAUDE_API_KEY environment variable is not set")
	}
	claudeClient := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &weaknessAnalysisService{
		db:                           db,
		repo:                         repo,
		correctResultsRepo:           correctResultsRepo,
		questionAnswersRepo:          questionAnswersRepo,
		questionTemplateMastersRepo:  questionTemplateMastersRepo,
		categoryMastersRepo:          categoryMastersRepo,
		weaknessCategoryAnalysisRepo: weaknessCategoryAnalysisRepo,
		weaknessDetailedAnalysisRepo: weaknessDetailedAnalysisRepo,
		weaknessLearningAdviceRepo:   weaknessLearningAdviceRepo,
		claudeClient:                 claudeClient,
		prompts:                      config.NewWeaknessAnalysisPrompts(),
	}
}

// CreateWeaknessAnalysis 学習弱点分析を作成する
// この時、学習カテゴリ分析、詳細分析、学習アドバイスを作成する
func (s *weaknessAnalysisService) CreateWeaknessAnalysis(userId string, req *model.CreateWeaknessAnalysisRequest) (*model.CreateWeaknessAnalysisResponse, error) {
	// トランザクション開始
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("トランザクション開始に失敗しました: %w", tx.Error)
	}

	// エラー時のロールバック処理
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 作成前に同じプロジェクトの分析が存在するか確認
	existingAnalysis, err := s.repo.GetWeaknessAnalysis(userId, &model.GetWeaknessAnalysisRequest{ProjectID: req.ProjectID})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if existingAnalysis != nil {
		tx.Rollback()
		return nil, fmt.Errorf("同じプロジェクトの分析が既に存在します")
	}

	// 弱点分析レコードを作成
	weaknessAnalysis, err := s.repo.CreateWeaknessAnalysis(userId, req)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 学習カテゴリ分析を作成する
	categoryAnalysisResults, err := s.WeaknessCategoryAnalysis(userId, req.ProjectID)
	if err != nil {
		// エラー時はステータスをFAILEDに更新
		s.repo.UpdateAnalysisStatus(weaknessAnalysis.ID, "FAILED")
		tx.Rollback()
		return nil, fmt.Errorf("カテゴリ分析の実行に失敗しました: %w", err)
	}

	// カテゴリ分析結果をパースしてデータベースに保存
	err = s.saveCategoryAnalysisResults(userId, weaknessAnalysis.ID, req.ProjectID, categoryAnalysisResults)
	if err != nil {
		// エラー時はステータスをFAILEDに更新
		s.repo.UpdateAnalysisStatus(weaknessAnalysis.ID, "FAILED")
		tx.Rollback()
		return nil, fmt.Errorf("カテゴリ分析結果の保存に失敗しました: %w", err)
	}

	// 詳細分析を作成する
	detailedAnalysisResult, err := s.WeaknessDetailedAnalysis(userId, req.ProjectID)
	if err != nil {
		// エラー時はステータスをFAILEDに更新
		s.repo.UpdateAnalysisStatus(weaknessAnalysis.ID, "FAILED")
		tx.Rollback()
		return nil, fmt.Errorf("詳細分析の実行に失敗しました: %w", err)
	}

	// 詳細分析結果をデータベースに保存
	err = s.saveDetailedAnalysisResult(userId, weaknessAnalysis.ID, detailedAnalysisResult)
	if err != nil {
		// エラー時はステータスをFAILEDに更新
		s.repo.UpdateAnalysisStatus(weaknessAnalysis.ID, "FAILED")
		tx.Rollback()
		return nil, fmt.Errorf("詳細分析結果の保存に失敗しました: %w", err)
	}

	// 学習アドバイスを作成する
	learningAdviceResult, err := s.WeaknessLearningAdvice(userId, req.ProjectID, detailedAnalysisResult)
	if err != nil {
		// エラー時はステータスをFAILEDに更新
		s.repo.UpdateAnalysisStatus(weaknessAnalysis.ID, "FAILED")
		tx.Rollback()
		return nil, fmt.Errorf("学習アドバイスの実行に失敗しました: %w", err)
	}

	// 学習アドバイス結果をデータベースに保存
	err = s.saveLearningAdviceResult(userId, weaknessAnalysis.ID, learningAdviceResult)
	if err != nil {
		// エラー時はステータスをFAILEDに更新
		s.repo.UpdateAnalysisStatus(weaknessAnalysis.ID, "FAILED")
		tx.Rollback()
		return nil, fmt.Errorf("学習アドバイス結果の保存に失敗しました: %w", err)
	}

	// 詳細分析結果から総合スコアを計算して更新
	overallScore := s.calculateOverallScore(detailedAnalysisResult)
	err = s.repo.UpdateOverallScore(weaknessAnalysis.ID, overallScore)
	if err != nil {
		// エラー時はステータスをFAILEDに更新
		s.repo.UpdateAnalysisStatus(weaknessAnalysis.ID, "FAILED")
		tx.Rollback()
		return nil, fmt.Errorf("総合スコアの更新に失敗しました: %w", err)
	}

	// 全ての分析が完了したので、ステータスをCOMPLETEDに更新
	err = s.repo.UpdateAnalysisStatus(weaknessAnalysis.ID, "COMPLETED")
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("分析ステータスの更新に失敗しました: %w", err)
	}

	// トランザクションをコミット
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("トランザクションのコミットに失敗しました: %w", err)
	}

	// レスポンスのステータスと総合スコアも更新
	weaknessAnalysis.AnalysisStatus = "COMPLETED"
	weaknessAnalysis.OverallScore = overallScore

	return weaknessAnalysis, nil
}

// GetWeaknessAnalysis 学習弱点分析を取得する
func (s *weaknessAnalysisService) GetWeaknessAnalysis(userId string, req *model.GetWeaknessAnalysisRequest) (*model.GetWeaknessAnalysisResponse, error) {
	return s.repo.GetWeaknessAnalysis(userId, req)
}

// weaknessCategoryの分析をLLMにて行う
func (s *weaknessAnalysisService) WeaknessCategoryAnalysis(userId string, projectId string) (map[string]*CategoryAnalysisResult, error) {
	// 解答データを取得
	correctResults, err := s.correctResultsRepo.GetCorrectResults(&model.GetCorrectResultsRequest{ProjectID: projectId})
	if err != nil {
		return nil, err
	}

	llmRequests := []model.LLMWeaknessCategoryAnalysisRequest{}
	categoryMap := make(map[string]string) // categoryName -> categoryID のマップ

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
		categoryMap[categoryMaster.CategoryMasters.Name] = categoryMaster.CategoryMasters.ID
	}

	// 学習データをJSON形式に変換
	jsonData, err := json.MarshalIndent(llmRequests, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal learning data: %w", err)
	}

	fmt.Println("jsonData", string(jsonData))

	// カテゴリごとにグループ化
	categoryGroups := make(map[string][]model.LLMWeaknessCategoryAnalysisRequest)
	for _, req := range llmRequests {
		categoryGroups[req.CategoryName] = append(categoryGroups[req.CategoryName], req)
	}

	// 各カテゴリごとに分析を実行
	results := make(map[string]*CategoryAnalysisResult)

	for categoryName, categoryRequests := range categoryGroups {
		// カテゴリ別の学習データをJSON形式に変換
		categoryJsonData, err := json.MarshalIndent(categoryRequests, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal category learning data: %w", err)
		}

		fmt.Printf("Category %s jsonData %s\n", categoryName, string(categoryJsonData))

		// プロンプト整形
		prompt := s.prompts.GetCategoryAnalysisPrompt(categoryName, string(categoryJsonData))

		// Claudeに分析リクエストを送信
		msg, err := s.claudeClient.Messages.New(
			context.Background(),
			anthropic.MessageNewParams{
				Model:     anthropic.ModelClaude3_7Sonnet20250219,
				MaxTokens: 6000,
				Messages: []anthropic.MessageParam{
					anthropic.NewUserMessage(
						anthropic.NewTextBlock(prompt),
					),
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("Claudeによる分析に失敗しました（%s）: %w", categoryName, err)
		}

		// レスポンスをパース
		var output string
		for _, block := range msg.Content {
			output += block.Text
		}

		fmt.Printf("Category %s raw Claude output: %s\n", categoryName, output)

		// JSONオブジェクトを抽出
		jsonOutput, err := utils.ExtractFirstJSONObject(output)
		if err != nil {
			fmt.Printf("Failed to extract JSON from Claude output（%s）: %s\n", categoryName, output)
			return nil, fmt.Errorf("failed to extract JSON object from Claude output（%s）: %w", categoryName, err)
		}
		fmt.Printf("Category %s extracted jsonOutput: %s\n", categoryName, jsonOutput)

		// JSONをパース
		var analysisResult CategoryAnalysisResult
		if err := json.Unmarshal([]byte(jsonOutput), &analysisResult); err != nil {
			fmt.Printf("JSON unmarshal error for category %s: %v\n", categoryName, err)
			fmt.Printf("Problematic JSON: %s\n", jsonOutput)

			// フォールバック：デフォルト値を設定
			fmt.Printf("Using fallback default values for category %s\n", categoryName)
			analysisResult = CategoryAnalysisResult{
				IsWeakness: false,
				IsStrength: false,
				Issues:     []string{"分析中にエラーが発生しました"},
				Strengths:  []string{},
				Examples:   []string{},
			}
		}

		results[categoryName] = &analysisResult
	}

	return results, nil
}

// WeaknessDetailedAnalysis 詳細分析をLLMにて行う
func (s *weaknessAnalysisService) WeaknessDetailedAnalysis(userId string, projectId string) (*model.DetailedAnalysisResult, error) {
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

	// 学習データをJSON形式に変換
	jsonData, err := json.MarshalIndent(llmRequests, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal detailed analysis data: %w", err)
	}

	fmt.Println("Detailed Analysis jsonData:", string(jsonData))

	// プロンプト整形
	prompt := s.prompts.GetDetailedAnalysisPrompt(string(jsonData))

	// Claudeに分析リクエストを送信
	msg, err := s.claudeClient.Messages.New(
		context.Background(),
		anthropic.MessageNewParams{
			Model:     anthropic.ModelClaude3_7Sonnet20250219,
			MaxTokens: 8000,
			Messages: []anthropic.MessageParam{
				anthropic.NewUserMessage(
					anthropic.NewTextBlock(prompt),
				),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Claudeによる詳細分析に失敗しました: %w", err)
	}

	// レスポンスをパース
	var output string
	for _, block := range msg.Content {
		output += block.Text
	}

	fmt.Printf("Detailed Analysis raw Claude output: %s\n", output)

	// JSONオブジェクトを抽出
	jsonOutput, err := utils.ExtractFirstJSONObject(output)
	if err != nil {
		fmt.Printf("Failed to extract JSON from detailed analysis output: %s\n", output)
		return nil, fmt.Errorf("failed to extract JSON object from detailed analysis output: %w", err)
	}
	fmt.Printf("Detailed Analysis extracted jsonOutput: %s\n", jsonOutput)

	// JSONをパース
	var detailedAnalysisResult model.DetailedAnalysisResult
	if err := json.Unmarshal([]byte(jsonOutput), &detailedAnalysisResult); err != nil {
		fmt.Printf("JSON unmarshal error for detailed analysis: %v\n", err)
		fmt.Printf("Problematic JSON: %s\n", jsonOutput)

		// フォールバック：デフォルト値を設定
		fmt.Println("Using fallback default values for detailed analysis")
		detailedAnalysisResult = model.DetailedAnalysisResult{
			Grammar: model.AnalysisDetail{
				Score:       50,
				Description: "分析中にエラーが発生しました",
				Examples:    []string{},
			},
			Vocabulary: model.AnalysisDetail{
				Score:       50,
				Description: "分析中にエラーが発生しました",
				Examples:    []string{},
			},
			Expression: model.AnalysisDetail{
				Score:       50,
				Description: "分析中にエラーが発生しました",
				Examples:    []string{},
			},
			Structure: model.AnalysisDetail{
				Score:       50,
				Description: "分析中にエラーが発生しました",
				Examples:    []string{},
			},
		}
	}

	return &detailedAnalysisResult, nil
}

// calculateOverallScore 詳細分析結果から総合スコアを計算する
func (s *weaknessAnalysisService) calculateOverallScore(detailedAnalysis *model.DetailedAnalysisResult) int {
	// 4つの領域（文法・語彙・表現・構成）のスコアの平均を計算
	totalScore := detailedAnalysis.Grammar.Score +
		detailedAnalysis.Vocabulary.Score +
		detailedAnalysis.Expression.Score +
		detailedAnalysis.Structure.Score

	overallScore := totalScore / 4

	// 0-100の範囲に正規化
	if overallScore < 0 {
		overallScore = 0
	} else if overallScore > 100 {
		overallScore = 100
	}

	return overallScore
}

// WeaknessLearningAdvice 学習アドバイスをLLMにて行う
func (s *weaknessAnalysisService) WeaknessLearningAdvice(userId string, projectId string, detailedAnalysis *model.DetailedAnalysisResult) (*model.PersonalizedAdvice, error) {
	// 詳細分析結果をJSON形式に変換
	jsonData, err := json.MarshalIndent(detailedAnalysis, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal detailed analysis data: %w", err)
	}

	fmt.Println("Learning Advice jsonData:", string(jsonData))

	// プロンプト整形
	prompt := s.prompts.GetLearningAdvicePrompt(string(jsonData))

	// Claudeに分析リクエストを送信
	msg, err := s.claudeClient.Messages.New(
		context.Background(),
		anthropic.MessageNewParams{
			Model:     anthropic.ModelClaude3_7Sonnet20250219,
			MaxTokens: 8000,
			Messages: []anthropic.MessageParam{
				anthropic.NewUserMessage(
					anthropic.NewTextBlock(prompt),
				),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Claudeによる学習アドバイス生成に失敗しました: %w", err)
	}

	// レスポンスをパース
	var output string
	for _, block := range msg.Content {
		output += block.Text
	}

	fmt.Printf("Learning Advice raw Claude output: %s\n", output)

	// JSONオブジェクトを抽出
	jsonOutput, err := utils.ExtractFirstJSONObject(output)
	if err != nil {
		fmt.Printf("Failed to extract JSON from learning advice output: %s\n", output)
		return nil, fmt.Errorf("failed to extract JSON object from learning advice output: %w", err)
	}
	fmt.Printf("Learning Advice extracted jsonOutput: %s\n", jsonOutput)

	// JSONをパース
	var learningAdviceResult model.PersonalizedAdvice
	if err := json.Unmarshal([]byte(jsonOutput), &learningAdviceResult); err != nil {
		fmt.Printf("JSON unmarshal error for learning advice: %v\n", err)
		fmt.Printf("Problematic JSON: %s\n", jsonOutput)

		// フォールバック：デフォルト値を設定
		fmt.Println("Using fallback default values for learning advice")
		learningAdviceResult = model.PersonalizedAdvice{
			LearningAdvice:      "分析中にエラーが発生しました。基本的な英語学習を継続してください。",
			RecommendedActions:  []string{"基本的な文法学習", "語彙力向上", "リーディング練習"},
			NextGoals:           []string{"基礎力向上", "継続的な学習習慣の確立"},
			StudyPlan:           "毎日30分の英語学習を継続し、基礎力を向上させましょう。",
			MotivationalMessage: "継続は力なり。一歩ずつ着実に進歩していきましょう。",
		}
	}

	return &learningAdviceResult, nil
}

// saveLearningAdviceResult 学習アドバイス結果をデータベースに保存する
func (s *weaknessAnalysisService) saveLearningAdviceResult(userId string, analysisId string, result *model.PersonalizedAdvice) error {
	// JSON配列を文字列に変換
	recommendedActionsJSON, _ := json.Marshal(result.RecommendedActions)
	nextGoalsJSON, _ := json.Marshal(result.NextGoals)

	// データベースに保存
	req := &model.CreateWeaknessLearningAdviceRequest{
		AnalysisID:          analysisId,
		LearningAdvice:      result.LearningAdvice,
		RecommendedActions:  string(recommendedActionsJSON),
		NextGoals:           string(nextGoalsJSON),
		StudyPlan:           result.StudyPlan,
		MotivationalMessage: result.MotivationalMessage,
	}

	_, err := s.weaknessLearningAdviceRepo.CreateWeaknessLearningAdvice(userId, req)
	if err != nil {
		return fmt.Errorf("学習アドバイス結果保存エラー: %w", err)
	}

	return nil
}

// saveCategoryAnalysisResults カテゴリ分析結果をデータベースに保存する
func (s *weaknessAnalysisService) saveCategoryAnalysisResults(userId string, analysisId string, projectId string, results map[string]*CategoryAnalysisResult) error {
	for categoryName, result := range results {
		// カテゴリIDを取得（カテゴリ名から検索）
		categoryMaster, err := s.categoryMastersRepo.GetCategoryMastersByName(categoryName)
		if err != nil {
			return fmt.Errorf("カテゴリマスター取得エラー（%s）: %w", categoryName, err)
		}
		if categoryMaster == nil {
			return fmt.Errorf("カテゴリマスターが見つかりません: %s", categoryName)
		}

		// カテゴリごとのスコアを計算（correction_resultsテーブルのデータに基づく）
		score, err := s.calculateCategoryScore(projectId, categoryMaster.CategoryMasters.ID)
		if err != nil {
			fmt.Printf("カテゴリ %s のスコア計算でエラー: %v\n", categoryName, err)
			// エラーの場合はデフォルト値を使用
			score = 50
		}

		// スコアに基づいて強み・弱みの判定を調整
		// LLMの分析結果とスコアが矛盾する場合はスコアを優先
		if score >= 80 && !result.IsStrength {
			fmt.Printf("カテゴリ %s: スコア %d に基づいて強みに調整\n", categoryName, score)
			result.IsStrength = true
			result.IsWeakness = false
		} else if score <= 40 && !result.IsWeakness {
			fmt.Printf("カテゴリ %s: スコア %d に基づいて弱みに調整\n", categoryName, score)
			result.IsWeakness = true
			result.IsStrength = false
		}

		// JSON配列を文字列に変換
		issuesJSON, _ := json.Marshal(result.Issues)
		strengthsJSON, _ := json.Marshal(result.Strengths)
		examplesJSON, _ := json.Marshal(result.Examples)

		// データベースに保存
		req := &model.CreateWeaknessCategoryAnalysisRequest{
			AnalysisID:   analysisId,
			CategoryID:   categoryMaster.CategoryMasters.ID,
			CategoryName: categoryName,
			Score:        score,
			IsWeakness:   result.IsWeakness,
			IsStrength:   result.IsStrength,
			Issues:       string(issuesJSON),
			Strengths:    string(strengthsJSON),
			Examples:     string(examplesJSON),
		}

		_, err = s.weaknessCategoryAnalysisRepo.CreateWeaknessCategoryAnalysis(userId, req)
		if err != nil {
			return fmt.Errorf("カテゴリ分析結果保存エラー（%s）: %w", categoryName, err)
		}
	}

	return nil
}

// calculateCategoryScore カテゴリごとのスコアを計算する（correction_resultsテーブルのデータに基づく）
func (s *weaknessAnalysisService) calculateCategoryScore(projectId string, categoryId string) (int, error) {
	// プロジェクトの全correction_resultsを取得
	correctResults, err := s.correctResultsRepo.GetCorrectResults(&model.GetCorrectResultsRequest{ProjectID: projectId})
	if err != nil {
		return 0, fmt.Errorf("correction_results取得エラー: %w", err)
	}

	if len(correctResults.CorrectResults) == 0 {
		return 50, nil // データがない場合はデフォルト値
	}

	// カテゴリに該当する問題の正答率を集計
	var totalCorrectRate int
	var categoryQuestionCount int

	for _, correctResult := range correctResults.CorrectResults {
		// 問題テンプレートマスターを取得してカテゴリIDを確認
		questionTemplate, err := s.questionTemplateMastersRepo.GetQuestionTemplateMasterByID(correctResult.QuestionTemplateMasterID)
		if err != nil {
			continue // エラーの場合はスキップ
		}

		if questionTemplate != nil && questionTemplate.CategoryID == categoryId {
			totalCorrectRate += correctResult.CorrectRate
			categoryQuestionCount++
		}
	}

	if categoryQuestionCount == 0 {
		return 50, nil // 該当カテゴリの問題がない場合はデフォルト値
	}

	// 平均正答率を計算
	averageScore := totalCorrectRate / categoryQuestionCount

	// 0-100の範囲に正規化
	if averageScore < 0 {
		averageScore = 0
	} else if averageScore > 100 {
		averageScore = 100
	}

	return averageScore, nil
}

// saveDetailedAnalysisResult 詳細分析結果をデータベースに保存する
func (s *weaknessAnalysisService) saveDetailedAnalysisResult(userId string, analysisId string, result *model.DetailedAnalysisResult) error {
	// JSON配列を文字列に変換
	grammarExamplesJSON, _ := json.Marshal(result.Grammar.Examples)
	vocabularyExamplesJSON, _ := json.Marshal(result.Vocabulary.Examples)
	expressionExamplesJSON, _ := json.Marshal(result.Expression.Examples)
	structureExamplesJSON, _ := json.Marshal(result.Structure.Examples)

	// データベースに保存
	req := &model.CreateWeaknessDetailedAnalysisRequest{
		AnalysisID:            analysisId,
		GrammarScore:          result.Grammar.Score,
		GrammarDescription:    result.Grammar.Description,
		GrammarExamples:       string(grammarExamplesJSON),
		VocabularyScore:       result.Vocabulary.Score,
		VocabularyDescription: result.Vocabulary.Description,
		VocabularyExamples:    string(vocabularyExamplesJSON),
		ExpressionScore:       result.Expression.Score,
		ExpressionDescription: result.Expression.Description,
		ExpressionExamples:    string(expressionExamplesJSON),
		StructureScore:        result.Structure.Score,
		StructureDescription:  result.Structure.Description,
		StructureExamples:     string(structureExamplesJSON),
	}

	_, err := s.weaknessDetailedAnalysisRepo.CreateWeaknessDetailedAnalysis(userId, req)
	if err != nil {
		return fmt.Errorf("詳細分析結果保存エラー: %w", err)
	}

	return nil
}

// 分析結果を全て取得
func (s *weaknessAnalysisService) GetWeaknessAnalysisAllSummary(userId string, projectId string) (*model.WeaknessAnalysisAllSummary, error) {
	weaknessAnalysisSummary, err := s.repo.GetWeaknessAnalysis(userId, &model.GetWeaknessAnalysisRequest{ProjectID: projectId})
	if err != nil {
		return nil, fmt.Errorf("分析結果取得エラー: %w", err)
	}
	if weaknessAnalysisSummary == nil {
		return nil, fmt.Errorf("分析結果が見つかりません")
	}

	weaknessCategoryAnalysisSummary, err := s.weaknessCategoryAnalysisRepo.GetWeaknessCategoryAnalysis(weaknessAnalysisSummary.Analysis.ID)
	if err != nil {
		return nil, fmt.Errorf("カテゴリ分析結果取得エラー: %w", err)
	}
	if weaknessCategoryAnalysisSummary == nil {
		return nil, fmt.Errorf("カテゴリ分析結果が見つかりません")
	}

	weaknessDetailedAnalysisSummary, err := s.weaknessDetailedAnalysisRepo.GetWeaknessDetailedAnalysis(weaknessAnalysisSummary.Analysis.ID)
	if err != nil {
		return nil, fmt.Errorf("詳細分析結果取得エラー: %w", err)
	}
	if weaknessDetailedAnalysisSummary == nil {
		return nil, fmt.Errorf("詳細分析結果が見つかりません")
	}

	weaknessLearningAdviceSummary, err := s.weaknessLearningAdviceRepo.GetWeaknessLearningAdvice(weaknessAnalysisSummary.Analysis.ID)
	if err != nil {
		return nil, fmt.Errorf("学習アドバイス結果取得エラー: %w", err)
	}
	if weaknessLearningAdviceSummary == nil {
		return nil, fmt.Errorf("学習アドバイス結果が見つかりません")
	}

	return &model.WeaknessAnalysisAllSummary{
		WeaknessAnalysisSummary:         weaknessAnalysisSummary.Analysis,
		WeaknessCategoryAnalysisSummary: []model.WeaknessCategoryAnalysisResponse{*weaknessCategoryAnalysisSummary},
		WeaknessDetailedAnalysisSummary: *weaknessDetailedAnalysisSummary,
		WeaknessLearningAdviceSummary:   *weaknessLearningAdviceSummary,
	}, nil
}
