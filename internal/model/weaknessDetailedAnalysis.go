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
