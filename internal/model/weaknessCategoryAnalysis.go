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
