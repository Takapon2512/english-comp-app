package model

import (
	"time"

	"gorm.io/gorm"
)

// WeaknessDetailedAnalysis は4つの主要学習領域の詳細分析結果を管理するテーブル
// 文法・語彙・表現・構成の各領域での詳細な分析結果を格納
type WeaknessDetailedAnalysis struct {
	// 基本識別情報
	ID         string `json:"id" gorm:"primaryKey;type:char(36)"`        // レコードの一意識別子
	AnalysisID string `json:"analysis_id" gorm:"type:char(36);not null"` // 親分析レコードのID（外部キー）

	// 文法領域の分析結果
	GrammarScore       int    `json:"grammar_score" gorm:"type:int;not null"` // 文法スコア（0-100）
	GrammarDescription string `json:"grammar_description" gorm:"type:text"`   // 文法面の詳細分析説明
	GrammarExamples    string `json:"grammar_examples" gorm:"type:json"`      // 文法の具体例JSON配列

	// 語彙領域の分析結果
	VocabularyScore       int    `json:"vocabulary_score" gorm:"type:int;not null"` // 語彙スコア（0-100）
	VocabularyDescription string `json:"vocabulary_description" gorm:"type:text"`   // 語彙面の詳細分析説明
	VocabularyExamples    string `json:"vocabulary_examples" gorm:"type:json"`      // 語彙の具体例JSON配列

	// 表現領域の分析結果
	ExpressionScore       int    `json:"expression_score" gorm:"type:int;not null"` // 表現スコア（0-100）
	ExpressionDescription string `json:"expression_description" gorm:"type:text"`   // 表現面の詳細分析説明
	ExpressionExamples    string `json:"expression_examples" gorm:"type:json"`      // 表現の具体例JSON配列

	// 構成領域の分析結果
	StructureScore       int    `json:"structure_score" gorm:"type:int;not null"` // 構成スコア（0-100）
	StructureDescription string `json:"structure_description" gorm:"type:text"`   // 構成面の詳細分析説明
	StructureExamples    string `json:"structure_examples" gorm:"type:json"`      // 構成の具体例JSON配列

	// 標準的なデータベース管理フィールド
	CreatedAt time.Time      `json:"created_at" gorm:"not null"`               // レコード作成日時
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null"`               // レコード最終更新日時
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`                  // 論理削除日時
	CreatedBy string         `json:"created_by" gorm:"type:char(36);not null"` // レコード作成者のユーザーID
	UpdatedBy string         `json:"updated_by" gorm:"type:char(36);not null"` // レコード最終更新者のユーザーID
}

// weaknessDetailedAnalysisSummary は詳細分析結果のサマリーを表す構造体
type WeaknessDetailedAnalysisSummary struct {
	ID         string `json:"id"`         // レコードの一意識別子
	AnalysisID string `json:"analysis_id"` // 親分析レコードのID
	GrammarScore int    `json:"grammar_score"` // 文法スコア（0-100）
	GrammarDescription string `json:"grammar_description"` // 文法面の詳細分析説明
	GrammarExamples string `json:"grammar_examples"` // 文法の具体例JSON配列
	VocabularyScore int    `json:"vocabulary_score"` // 語彙スコア（0-100）
	VocabularyDescription string `json:"vocabulary_description"` // 語彙面の詳細分析説明
	VocabularyExamples string `json:"vocabulary_examples"` // 語彙の具体例JSON配列
	ExpressionScore int    `json:"expression_score"` // 表現スコア（0-100）
	ExpressionDescription string `json:"expression_description"` // 表現面の詳細分析説明
	ExpressionExamples string `json:"expression_examples"` // 表現の具体例JSON配列
	StructureScore int    `json:"structure_score"` // 構成スコア（0-100）
	StructureDescription string `json:"structure_description"` // 構成面の詳細分析説明
	StructureExamples string `json:"structure_examples"` // 構成の具体例JSON配列
}

// ===== 詳細分析用のサマリー構造体 =====

// DetailedAnalysisResult は4つの主要学習領域の詳細分析結果を表す構造体
// LLMが各領域を個別に評価し、具体的なフィードバックを提供
type DetailedAnalysisResult struct {
	Grammar    AnalysisDetail `json:"grammar"`    // 文法領域の詳細分析（時制、語順、構文など）
	Vocabulary AnalysisDetail `json:"vocabulary"` // 語彙領域の詳細分析（語彙力、使い分けなど）
	Expression AnalysisDetail `json:"expression"` // 表現領域の詳細分析（自然さ、流暢さなど）
	Structure  AnalysisDetail `json:"structure"`  // 構成領域の詳細分析（論理性、一貫性など）
}

// AnalysisDetail は各学習領域の詳細分析情報を表す構造体
// スコア、説明、具体例を含む包括的な評価結果
type AnalysisDetail struct {
	Score       int      `json:"score"`       // この領域でのスコア（0-100）
	Description string   `json:"description"` // LLMによる詳細な分析説明
	Examples    []string `json:"examples"`    // 具体的な例文や改善例の配列
}

// CreateWeaknessDetailedAnalysisRequest は詳細分析結果を作成するリクエスト用構造体
type CreateWeaknessDetailedAnalysisRequest struct {
	AnalysisID            string `json:"analysis_id" binding:"required"`            // 親分析レコードのID（必須）
	GrammarScore          int    `json:"grammar_score" binding:"required"`          // 文法スコア（0-100）
	GrammarDescription    string `json:"grammar_description" binding:"required"`    // 文法面の詳細分析説明
	GrammarExamples       string `json:"grammar_examples" binding:"required"`       // 文法の具体例JSON配列
	VocabularyScore       int    `json:"vocabulary_score" binding:"required"`       // 語彙スコア（0-100）
	VocabularyDescription string `json:"vocabulary_description" binding:"required"` // 語彙面の詳細分析説明
	VocabularyExamples    string `json:"vocabulary_examples" binding:"required"`    // 語彙の具体例JSON配列
	ExpressionScore       int    `json:"expression_score" binding:"required"`       // 表現スコア（0-100）
	ExpressionDescription string `json:"expression_description" binding:"required"` // 表現面の詳細分析説明
	ExpressionExamples    string `json:"expression_examples" binding:"required"`    // 表現の具体例JSON配列
	StructureScore        int    `json:"structure_score" binding:"required"`        // 構成スコア（0-100）
	StructureDescription  string `json:"structure_description" binding:"required"`  // 構成面の詳細分析説明
	StructureExamples     string `json:"structure_examples" binding:"required"`     // 構成の具体例JSON配列
}

// CreateWeaknessDetailedAnalysisResponse は詳細分析結果を作成するレスポンス用構造体
type CreateWeaknessDetailedAnalysisResponse struct {
	ID                    string `json:"id"`                     // 作成された詳細分析レコードのID
	AnalysisID            string `json:"analysis_id"`            // 親分析レコードのID
	GrammarScore          int    `json:"grammar_score"`          // 文法スコア（0-100）
	GrammarDescription    string `json:"grammar_description"`    // 文法面の詳細分析説明
	GrammarExamples       string `json:"grammar_examples"`       // 文法の具体例JSON配列
	VocabularyScore       int    `json:"vocabulary_score"`       // 語彙スコア（0-100）
	VocabularyDescription string `json:"vocabulary_description"` // 語彙面の詳細分析説明
	VocabularyExamples    string `json:"vocabulary_examples"`    // 語彙の具体例JSON配列
	ExpressionScore       int    `json:"expression_score"`       // 表現スコア（0-100）
	ExpressionDescription string `json:"expression_description"` // 表現面の詳細分析説明
	ExpressionExamples    string `json:"expression_examples"`    // 表現の具体例JSON配列
}

// LLMに詳細分析を依頼するリクエスト用構造体
type LLMWeaknessDetailedAnalysisRequest struct {
	Question      string `json:"question"`       // 問題文
	UserAnswer    string `json:"user_answer"`    // 解答文
	CorrectAnswer string `json:"correct_answer"` // 正解文
}
