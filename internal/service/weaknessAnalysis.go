package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/repository"
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
	WeaknessDetailedAnalysis(userId string, projectId string) ([]model.LLMWeaknessDetailedAnalysisRequest, error)
	WeaknessLearningAdvice(userId string, projectId string) (model.LLMWeaknessLearningAdviceRequest, error)
}

type weaknessAnalysisService struct {
	db                           *gorm.DB
	repo                         repository.WeaknessAnalysisRepository
	correctResultsRepo           repository.CorrectResultsRepository
	questionAnswersRepo          repository.QuestionAnswersRepository
	questionTemplateMastersRepo  repository.QuestionTemplateMastersRepository
	categoryMastersRepo          repository.CategoryMastersRepository
	weaknessCategoryAnalysisRepo repository.WeaknessCategoryAnalysisRepository
	claudeClient                 anthropic.Client
}

func NewWeaknessAnalysisService(
	db *gorm.DB,
	repo repository.WeaknessAnalysisRepository,
	correctResultsRepo repository.CorrectResultsRepository,
	questionAnswersRepo repository.QuestionAnswersRepository,
	questionTemplateMastersRepo repository.QuestionTemplateMastersRepository,
	categoryMastersRepo repository.CategoryMastersRepository,
	weaknessCategoryAnalysisRepo repository.WeaknessCategoryAnalysisRepository,
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
		claudeClient:                 claudeClient,
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
		tx.Rollback()
		return nil, fmt.Errorf("カテゴリ分析の実行に失敗しました: %w", err)
	}

	// カテゴリ分析結果をパースしてデータベースに保存
	err = s.saveCategoryAnalysisResults(userId, weaknessAnalysis.ID, req.ProjectID, categoryAnalysisResults)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("カテゴリ分析結果の保存に失敗しました: %w", err)
	}

	// 詳細分析を作成する
	// llmRequestsDetailedAnalysis, err := s.WeaknessDetailedAnalysis(userId, req.ProjectID)
	// if err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }
	// b, _ = json.MarshalIndent(llmRequestsDetailedAnalysis, "", "  ")
	// fmt.Println(string(b))

	// 学習アドバイスを作成する

	// トランザクションをコミット
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("トランザクションのコミットに失敗しました: %w", err)
	}

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
		prompt := `あなたはプロの英語教師です。
以下の「` + categoryName + `」カテゴリの学習データを分析し、このカテゴリでの学習者の強み・弱みを分析してJSON形式で出力してください。

【重要】以下の要件を厳密に守ってください：
1. 出力は有効なJSON形式のみにしてください
2. 説明文やマークダウン記法は一切含めないでください
3. JSONの前後に余計な文字を入れないでください
4. 配列が空の場合は空配列[]を使用してください

出力JSON形式：
{
  "is_weakness": boolean,
  "is_strength": boolean,
  "issues": ["問題点1", "問題点2"],
  "strengths": ["強み1", "強み2"],
  "examples": ["具体例1", "具体例2"]
}

分析対象データ:
` + string(categoryJsonData) + `

上記データを分析し、有効なJSONのみを出力してください：`

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
		jsonOutput, err := extractFirstJSONObject(output)
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

// WeaknessLearningAdvice 学習アドバイスをLLMにて行う
func (s *weaknessAnalysisService) WeaknessLearningAdvice(userId string, projectId string) (model.LLMWeaknessLearningAdviceRequest, error) {
	// 分析データを取得
	var analysis model.WeaknessAnalysis
	if err := s.db.Where("user_id = ? AND project_id = ?", userId, projectId).Order("created_at desc").First(&analysis).Error; err != nil {
		return model.LLMWeaknessLearningAdviceRequest{}, fmt.Errorf("failed to get weakness analysis: %w", err)
	}

	// 詳細分析結果を取得
	var detailedAnalysis model.WeaknessDetailedAnalysis
	if err := s.db.Where("analysis_id = ?", analysis.ID).First(&detailedAnalysis).Error; err != nil {
		return model.LLMWeaknessLearningAdviceRequest{}, fmt.Errorf("failed to get weakness detailed analysis: %w", err)
	}

	// LLMリクエストを作成
	llmRequest := model.LLMWeaknessLearningAdviceRequest{
		GrammarScore:          detailedAnalysis.GrammarScore,
		GrammarDescription:    detailedAnalysis.GrammarDescription,
		VocabularyScore:       detailedAnalysis.VocabularyScore,
		VocabularyDescription: detailedAnalysis.VocabularyDescription,
		ExpressionScore:       detailedAnalysis.ExpressionScore,
		ExpressionDescription: detailedAnalysis.ExpressionDescription,
	}

	return llmRequest, nil
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
