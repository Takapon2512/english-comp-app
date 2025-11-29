package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/repository"
)

type CorrectResultsService interface {
	CreateCorrectionResult(userID string, req *model.CreateCorrectionResultRequest) (*model.CreateCorrectionResultResponse, error)
	GrandCorrectResult(userID string, req *model.GrandCorrectResultRequest) (*model.GrandCorrectResultResponse, error)
	GetCorrectResults(userID string, req *model.GetCorrectResultsRequest) (*model.GetCorrectResultsResponse, error)
	GetCorrectResultsVersionList(userID string, req *model.GetCorrectResultsVersionRequest) (*model.GetCorrectResultsVersionListResponse, error)
}

type correctResultsService struct {
	db                          *gorm.DB
	repo                        repository.CorrectResultsRepository
	questionTemplateMastersRepo repository.QuestionTemplateMastersRepository
	questionAnswersRepo         repository.QuestionAnswersRepository
	categoryMastersRepo         repository.CategoryMastersRepository
	claudeClient                anthropic.Client
}

func NewCorrectResultsService(
	db *gorm.DB,
	repo repository.CorrectResultsRepository,
	questionTemplateMastersRepo repository.QuestionTemplateMastersRepository,
	questionAnswersRepo repository.QuestionAnswersRepository,
	categoryMastersRepo repository.CategoryMastersRepository,
) CorrectResultsService {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		log.Fatal("CLAUDE_API_KEY environment variable is not set")
	}
	claudeClient := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &correctResultsService{
		db:                          db,
		repo:                        repo,
		questionTemplateMastersRepo: questionTemplateMastersRepo,
		questionAnswersRepo:         questionAnswersRepo,
		categoryMastersRepo:         categoryMastersRepo,
		claudeClient:                claudeClient,
	}
}

// 添削結果のデータを作成
func (s *correctResultsService) CreateCorrectionResult(userID string, req *model.CreateCorrectionResultRequest) (*model.CreateCorrectionResultResponse, error) {
	// 解答データ取得
	userAnswer, err := s.questionAnswersRepo.GetQuestionAnswerById(req.QuestionAnswerID)
	if err != nil {
		return nil, fmt.Errorf("解答データの取得に失敗しました: %w", err)
	}

	req.ChallengeCount = userAnswer.ChallengeCount

	return s.repo.CreateCorrectionResult(req)
}

// 添削結果のデータを更新
func (s *correctResultsService) UpdateCorrectionResult(userID string, req *model.UpdateCorrectionResultRequest) (*model.UpdateCorrectionResultResponse, error) {
	return s.repo.UpdateCorrectionResult(req)
}

// LLMで作成した添削結果のデータを取得し、反映
func (s *correctResultsService) GrandCorrectResult(userID string, req *model.GrandCorrectResultRequest) (*model.GrandCorrectResultResponse, error) {
	// 問題・解答データ取得のためのデータを取得
	correctionResult, err := s.repo.GetCorrectionResultById(req.ID)
	if err != nil {
		return nil, fmt.Errorf("添削結果の取得に失敗しました: %w", err)
	}

	// 解答データを取得
	userAnswer, err := s.questionAnswersRepo.GetQuestionAnswerById(correctionResult.QuestionAnswerID)
	if err != nil {
		return nil, fmt.Errorf("解答データの取得に失敗しました: %w", err)
	}

	// 質問テンプレートマスターを取得
	questionTemplateMaster, err := s.questionTemplateMastersRepo.GetQuestionTemplateMasterLLMById(correctionResult.QuestionTemplateMasterID)
	if err != nil {
		return nil, fmt.Errorf("質問テンプレートマスターの取得に失敗しました: %w", err)
	}

	// LLMで添削を行うプロンプトを作成
	prompt := fmt.Sprintf(`
		あなたは英語の作文を採点する教師です。以下の条件で採点を行ってください：
		
		問題：
		%s

		日本語での説明：
		%s

		学習者の解答：
		%s

		出力要件：
		- 次の厳密なJSONオブジェクト「のみ」を返してください。
		- コードブロック( バッククォート3つ )や前後の説明文、余計な文字は一切出力しないでください。
		- 値は有効なJSONとし、数値は整数で出力してください。
		- キーは英語のまま使用してください。
		- アドバイスは日本語で出力してください。
		
		出力フォーマット（参考）：
		{
			"points": 採点結果（%d点満点の整数）, 
			"correct_rate": 正答率（0-100の整数）, 
			"example_correction": 模範解答の文字列, 
			"advice": 改善のためのアドバイスの文字列
		}
	`, questionTemplateMaster.English, questionTemplateMaster.Japanese, userAnswer.UserAnswer, questionTemplateMaster.Points)

	log.Println("prompt", prompt)

	// Claudeに採点リクエストを送信（Anthropic Go SDK v1）
	msg, err := s.claudeClient.Messages.New(
		context.Background(),
		anthropic.MessageNewParams{
			Model:     anthropic.ModelClaude3_7Sonnet20250219,
			MaxTokens: 5000,
			Messages: []anthropic.MessageParam{
				anthropic.NewUserMessage(
					anthropic.NewTextBlock(prompt),
				),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Claudeによる採点に失敗しました: %w", err)
	}

	// レスポンスをパース（Contentのテキスト結合 → JSONとして解釈）
	var output string
	for _, block := range msg.Content {
		if block.Type == "text" {
			output += block.Text
		}
	}
	// Claudeの出力にコードフェンスや説明が混ざる場合があるため、
	// 最初の完全なJSONオブジェクトのみを抽出してからUnmarshalする
	jsonStr, err := extractFirstJSONObject(output)
	if err != nil {
		return nil, fmt.Errorf("ClaudeのレスポンスからJSON抽出に失敗しました: %w", err)
	}
	var llmResponse struct {
		Points            int    `json:"points"`
		CorrectRate       int    `json:"correct_rate"`
		ExampleCorrection string `json:"example_correction"`
		Advice            string `json:"advice"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &llmResponse); err != nil {
		return nil, fmt.Errorf("Claudeのレスポンスのパースに失敗しました: %w", err)
	}

	s.repo.UpdateCorrectionResult(&model.UpdateCorrectionResultRequest{
		ID:                correctionResult.ID,
		GetPoints:         llmResponse.Points,
		ExampleCorrection: llmResponse.ExampleCorrection,
		CorrectRate:       llmResponse.CorrectRate,
		Advice:            llmResponse.Advice,
		Status:            "COMPLETED",
	})

	return &model.GrandCorrectResultResponse{
		ID:                       correctionResult.ID,
		QuestionAnswerID:         correctionResult.QuestionAnswerID,
		QuestionTemplateMasterID: correctionResult.QuestionTemplateMasterID,
		GetPoints:                llmResponse.Points,
		ExampleCorrection:        llmResponse.ExampleCorrection,
		CorrectRate:              llmResponse.CorrectRate,
		Advice:                   llmResponse.Advice,
		Status:                   "COMPLETED",
		ChallengeCount:           correctionResult.ChallengeCount,
	}, nil
}

// 添削結果の取得
func (s *correctResultsService) GetCorrectResults(userID string, req *model.GetCorrectResultsRequest) (*model.GetCorrectResultsResponse, error) {
	// 添削結果の取得
	correctResults, err := s.repo.GetCorrectResults(req)
	if err != nil {
		return nil, fmt.Errorf("添削結果の取得に失敗しました: %w", err)
	}

	var correctResultsSummary []model.CorrectionResultsSummary
	for _, correctResult := range correctResults.CorrectResults {
		questionAnswer, err := s.questionAnswersRepo.GetQuestionAnswerById(correctResult.QuestionAnswerID)
		if err != nil {
			return nil, fmt.Errorf("質問回答の取得に失敗しました: %w", err)
		}
		if questionAnswer == nil {
			return nil, fmt.Errorf("質問回答（ID: %s）が見つかりません", correctResult.QuestionAnswerID)
		}

		questionTemplateMaster, err := s.questionTemplateMastersRepo.GetQuestionTemplateMasterByID(correctResult.QuestionTemplateMasterID)
		if err != nil {
			return nil, fmt.Errorf("質問テンプレートの取得に失敗しました: %w", err)
		}
		if questionTemplateMaster == nil {
			return nil, fmt.Errorf("質問テンプレート（ID: %s）が見つかりません", correctResult.QuestionTemplateMasterID)
		}

		categoryMaster, err := s.categoryMastersRepo.GetCategoryMastersByID(questionTemplateMaster.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("カテゴリマスターの取得に失敗しました: %w", err)
		}
		if categoryMaster == nil {
			return nil, fmt.Errorf("カテゴリマスター（ID: %s）が見つかりません", questionTemplateMaster.CategoryID)
		}

		correctResultsSummary = append(correctResultsSummary, model.CorrectionResultsSummary{
			ID:                       correctResult.ID,
			ProjectID:                correctResult.ProjectID,
			QuestionAnswerID:         correctResult.QuestionAnswerID,
			QuestionTemplateMasterID: correctResult.QuestionTemplateMasterID,
			GetPoints:                correctResult.GetPoints,
			ExampleCorrection:        correctResult.ExampleCorrection,
			CorrectRate:              correctResult.CorrectRate,
			Advice:                   correctResult.Advice,
			Status:                   correctResult.Status,
			ChallengeCount:           correctResult.ChallengeCount,
			QuestionAnswer: model.QuestionAnswersSummary{
				ID:                       correctResult.QuestionAnswerID,
				ProjectID:                correctResult.ProjectID,
				UserID:                   questionAnswer.UserID,
				UserAnswer:               questionAnswer.UserAnswer,
				QuestionTemplateMasterID: questionAnswer.QuestionTemplateMasterID,
			},
			QuestionTemplateMaster: model.QuestionTemplateMastersSummary{
				ID:            correctResult.QuestionTemplateMasterID,
				CategoryID:    questionTemplateMaster.CategoryID,
				QuestionType:  questionTemplateMaster.QuestionType,
				English:       questionTemplateMaster.English,
				Japanese:      questionTemplateMaster.Japanese,
				Points:        questionTemplateMaster.Points,
				Level:         questionTemplateMaster.Level,
				EstimatedTime: questionTemplateMaster.EstimatedTime,
				Status:        questionTemplateMaster.Status,
				Category: model.CategoryInfo{
					ID:   questionTemplateMaster.CategoryID,
					Name: categoryMaster.CategoryMasters.Name,
				},
			},
		})
	}

	return &model.GetCorrectResultsResponse{
		CorrectResults: correctResultsSummary,
	}, nil
}

// 添削結果のバージョン一覧を取得
func (s *correctResultsService) GetCorrectResultsVersionList(userID string, req *model.GetCorrectResultsVersionRequest) (*model.GetCorrectResultsVersionListResponse, error) {
	// リポジトリ層からデータを取得
	versionList, err := s.repo.GetCorrectResultsVersionList(req)
	if err != nil {
		return nil, err
	}

	// ビジネスロジック：空の場合のデフォルト処理
	if len(versionList) == 0 {
		// 添削結果が存在しない場合は空のリストを返す
		return &model.GetCorrectResultsVersionListResponse{
			VersionList: []model.VersionList{},
		}, nil
	}

	return &model.GetCorrectResultsVersionListResponse{
		VersionList: versionList,
	}, nil
}

// Claudeの出力テキストから最初の完全なJSONオブジェクトを抽出する
// - コードブロック( ```json ... ``` )が含まれていても取り除いて抽出する
func extractFirstJSONObject(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	fmt.Printf("extractFirstJSONObject input: %s\n", trimmed)

	// コードフェンスがある場合は除去
	if strings.HasPrefix(trimmed, "```") {
		fmt.Println("Removing code fences...")
		// 先頭のフェンスとオプションの言語指定を取り除く
		trimmed = strings.TrimPrefix(trimmed, "```json")
		trimmed = strings.TrimPrefix(trimmed, "```JSON")
		trimmed = strings.TrimPrefix(trimmed, "```")
		trimmed = strings.TrimSpace(trimmed)
		// 末尾のフェンスを取り除く
		if idx := strings.LastIndex(trimmed, "```"); idx != -1 {
			trimmed = trimmed[:idx]
		}
		trimmed = strings.TrimSpace(trimmed)
		fmt.Printf("After removing code fences: %s\n", trimmed)
	}

	// 最初の完全なJSONオブジェクトを括弧のネストで検出
	inString := false
	escape := false
	depth := 0
	start := -1
	for i, r := range trimmed {
		if inString {
			if escape {
				escape = false
				continue
			}
			if r == '\\' {
				escape = true
				continue
			}
			if r == '"' {
				inString = false
			}
			continue
		}
		if r == '"' {
			inString = true
			continue
		}
		switch r {
		case '{':
			if depth == 0 {
				start = i
			}
			depth++
		case '}':
			if depth > 0 {
				depth--
				if depth == 0 && start != -1 {
					result := trimmed[start : i+1]
					fmt.Printf("Extracted JSON: %s\n", result)
					return result, nil
				}
			}
		}
	}

	fmt.Printf("Failed to extract JSON from: %s\n", trimmed)
	return "", fmt.Errorf("valid JSON object not found in LLM output")
}
