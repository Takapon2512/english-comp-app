package model

import (
	"time"

	"gorm.io/gorm"
)

// WeaknessCategoryAnalysis はカテゴリ別の分析結果を管理するテーブル
// 各カテゴリ（文法項目）ごとの苦手・得意分析を個別レコードで管理
type WeaknessCategoryAnalysis struct {
	// 基本識別情報
	ID         string `json:"id" gorm:"primaryKey;type:char(36)"`        // レコードの一意識別子
	AnalysisID string `json:"analysis_id" gorm:"type:char(36);not null"` // 親分析レコードのID（外部キー）
	CategoryID string `json:"category_id" gorm:"type:char(36);not null"` // カテゴリマスターのID

	// カテゴリ分析結果
	CategoryName string `json:"category_name" gorm:"type:varchar(100);not null"` // カテゴリ名（キャッシュ用）
	Score        int    `json:"score" gorm:"type:int;not null"`                  // このカテゴリでのスコア（0-100）
	IsWeakness   bool   `json:"is_weakness" gorm:"type:boolean;not null"`        // 苦手カテゴリかどうか
	IsStrength   bool   `json:"is_strength" gorm:"type:boolean;not null"`        // 得意カテゴリかどうか

	// 分析詳細
	Issues    string `json:"issues" gorm:"type:json"`    // 具体的な問題点のJSON配列
	Strengths string `json:"strengths" gorm:"type:json"` // 具体的な強みのJSON配列
	Examples  string `json:"examples" gorm:"type:json"`  // 具体例のJSON配列

	// 標準的なデータベース管理フィールド
	CreatedAt time.Time      `json:"created_at" gorm:"not null"`               // レコード作成日時
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null"`               // レコード最終更新日時
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`                  // 論理削除日時
	CreatedBy string         `json:"created_by" gorm:"type:char(36);not null"` // レコード作成者のユーザーID
	UpdatedBy string         `json:"updated_by" gorm:"type:char(36);not null"` // レコード最終更新者のユーザーID
}

// weaknessCategoryAnalysisSummary はカテゴリ別の分析結果のサマリーを表す構造体
type WeaknessCategoryAnalysisSummary struct {
	ID           string `json:"id"`            // レコードの一意識別子
	AnalysisID   string `json:"analysis_id"`   // 親分析レコードのID
	CategoryID   string `json:"category_id"`   // カテゴリマスターのID
	CategoryName string `json:"category_name"` // カテゴリの表示名
	Score        int    `json:"score"`         // このカテゴリでのスコア（0-100）
	IsWeakness   bool   `json:"is_weakness"`   // 苦手カテゴリかどうか
	IsStrength   bool   `json:"is_strength"`   // 得意カテゴリかどうか
	Issues       string `json:"issues"`        // 具体的な問題点のJSON文字列
	Strengths    string `json:"strengths"`     // 具体的な強みのJSON文字列
	Examples     string `json:"examples"`      // 具体例のJSON文字列
}

// WeaknessCategoryAnalysisResponse はカテゴリ別の分析結果のレスポンス用構造体
type WeaknessCategoryAnalysisResponse struct {
	ID           string   `json:"id"`            // レコードの一意識別子
	AnalysisID   string   `json:"analysis_id"`   // 親分析レコードのID
	CategoryID   string   `json:"category_id"`   // カテゴリマスターのID
	CategoryName string   `json:"category_name"` // カテゴリの表示名
	Score        int      `json:"score"`         // このカテゴリでのスコア（0-100）
	IsWeakness   bool     `json:"is_weakness"`   // 苦手カテゴリかどうか
	IsStrength   bool     `json:"is_strength"`   // 得意カテゴリかどうか
	Issues       []string `json:"issues"`        // 具体的な問題点の配列（LLMが特定した課題）
	Strengths    []string `json:"strengths"`     // 具体的な強みの配列（LLMが評価した優秀な点）
	Examples     []string `json:"examples"`      // 具体例の配列（LLMが提供した改善例）
}

// ===== カテゴリ分析用のサマリー構造体 =====

// CategoryWeakness は苦手カテゴリの詳細情報を表す構造体
// LLMが特定した学習者の弱点領域とその具体的な問題点を含む
type CategoryWeakness struct {
	CategoryID   string   `json:"category_id"`   // カテゴリマスターのID（例: 時制、仮定法など）
	CategoryName string   `json:"category_name"` // カテゴリの表示名（日本語）
	Score        int      `json:"score"`         // このカテゴリでのスコア（0-100）
	Issues       []string `json:"issues"`        // 具体的な問題点の配列（LLMが特定した課題）
}

// CategoryStrength は得意カテゴリの詳細情報を表す構造体
// LLMが認識した学習者の強みと優れている点を含む
type CategoryStrength struct {
	CategoryID   string   `json:"category_id"`   // カテゴリマスターのID
	CategoryName string   `json:"category_name"` // カテゴリの表示名（日本語）
	Score        int      `json:"score"`         // このカテゴリでのスコア（0-100）
	Strengths    []string `json:"strengths"`     // 具体的な強みの配列（LLMが評価した優秀な点）
}

// CreateWeaknessCategoryAnalysisRequest はカテゴリ別の分析結果を作成するリクエスト用構造体
type CreateWeaknessCategoryAnalysisRequest struct {
	AnalysisID   string `json:"analysis_id" binding:"required"`   // 親分析レコードのID（必須）
	CategoryID   string `json:"category_id" binding:"required"`   // カテゴリマスターのID（必須）
	CategoryName string `json:"category_name" binding:"required"` // カテゴリの表示名（必須）
	Score        int    `json:"score" binding:"required"`         // このカテゴリでのスコア（必須）
	IsWeakness   bool   `json:"is_weakness" binding:"required"`   // 苦手カテゴリかどうか（必須）
	IsStrength   bool   `json:"is_strength" binding:"required"`   // 得意カテゴリかどうか（必須）
	Issues       string `json:"issues" binding:"required"`        // 具体的な問題点のJSON配列（必須）
	Strengths    string `json:"strengths" binding:"required"`     // 具体的な強みのJSON配列（必須）
	Examples     string `json:"examples" binding:"required"`      // 具体例のJSON配列（必須）
}

// UpdateWeaknessCategoryAnalysisRequest はカテゴリ別の分析結果を更新するリクエスト用構造体
type UpdateWeaknessCategoryAnalysisRequest struct {
	ID           string `json:"id" binding:"required"`           // 更新対象のカテゴリ別分析レコードのID（必須）
	AnalysisID   string `json:"analysis_id" binding:"required"`   // 親分析レコードのID（必須）
	CategoryID   string `json:"category_id" binding:"required"`   // カテゴリマスターのID（必須）
	CategoryName string `json:"category_name" binding:"required"` // カテゴリの表示名（必須）
	Score        int    `json:"score" binding:"required"`         // このカテゴリでのスコア（必須）
	IsWeakness   bool   `json:"is_weakness" binding:"required"`   // 苦手カテゴリかどうか（必須）
	IsStrength   bool   `json:"is_strength" binding:"required"`   // 得意カテゴリかどうか（必須）
	Issues       string `json:"issues" binding:"required"`        // 具体的な問題点のJSON配列（必須）
	Strengths    string `json:"strengths" binding:"required"`     // 具体的な強みのJSON配列（必須）
	Examples     string `json:"examples" binding:"required"`      // 具体例のJSON配列（必須）
}

// CreateWeaknessCategoryAnalysisResponse はカテゴリ別の分析結果を作成するレスポンス用構造体
type CreateWeaknessCategoryAnalysisResponse struct {
	ID           string `json:"id"`            // 作成されたカテゴリ別分析レコードのID
	AnalysisID   string `json:"analysis_id"`   // 親分析レコードのID
	CategoryID   string `json:"category_id"`   // カテゴリマスターのID
	CategoryName string `json:"category_name"` // カテゴリの表示名
	Score        int    `json:"score"`         // このカテゴリでのスコア
	IsWeakness   bool   `json:"is_weakness"`   // 苦手カテゴリかどうか
	IsStrength   bool   `json:"is_strength"`   // 得意カテゴリかどうか
	Issues       string `json:"issues"`        // 具体的な問題点のJSON配列
	Strengths    string `json:"strengths"`     // 具体的な強みのJSON配列
	Examples     string `json:"examples"`      // 具体例のJSON配列
}

// UpdateWeaknessCategoryAnalysisResponse はカテゴリ別の分析結果を更新するレスポンス用構造体
type UpdateWeaknessCategoryAnalysisResponse struct {
	ID           string `json:"id"`            // 更新されたカテゴリ別分析レコードのID
	AnalysisID   string `json:"analysis_id"`   // 親分析レコードのID
	CategoryID   string `json:"category_id"`   // カテゴリマスターのID
	CategoryName string `json:"category_name"` // カテゴリの表示名
	Score        int    `json:"score"`         // このカテゴリでのスコア
	IsWeakness   bool   `json:"is_weakness"`   // 苦手カテゴリかどうか
	IsStrength   bool   `json:"is_strength"`   // 得意カテゴリかどうか
	Issues       string `json:"issues"`        // 具体的な問題点のJSON配列
	Strengths    string `json:"strengths"`     // 具体的な強みのJSON配列
	Examples     string `json:"examples"`      // 具体例のJSON配列
}

// LLMに分析を依頼するリクエスト用構造体
type LLMWeaknessCategoryAnalysisRequest struct {
	CategoryName  string `json:"category_name"`  // カテゴリの表示名
	Question      string `json:"question"`       // 問題文
	UserAnswer    string `json:"user_answer"`    // 解答文
	CorrectAnswer string `json:"correct_answer"` // 正解文
}
