package model

import (
	"time"

	"gorm.io/gorm"
)

// WeaknessAnalysis はプロジェクト内でのユーザーの学習弱点分析結果のメインエンティティ
// LLM（Claude API）を使用した分析の基本情報とメタデータを管理する
type WeaknessAnalysis struct {
	// 基本識別情報
	ID        string `json:"id" gorm:"primaryKey;type:char(36)"`       // 分析結果の一意識別子（UUID）
	ProjectID string `json:"project_id" gorm:"type:char(36);not null"` // 分析対象プロジェクトのID
	UserID    string `json:"user_id" gorm:"type:char(36);not null"`    // 分析対象ユーザーのID

	// LLM分析結果の基本情報
	AnalysisStatus  string `json:"analysis_status" gorm:"type:varchar(20);not null;default:'PROCESSING'"` // 分析処理状況（PROCESSING: 処理中, COMPLETED: 完了, FAILED: 失敗）
	OverallScore    int    `json:"overall_score" gorm:"type:int;default:0"`                               // LLMが算出した総合学習スコア（0-100点）
	ImprovementRate int    `json:"improvement_rate" gorm:"type:int;default:0"`                            // 前回分析からの改善率（パーセンテージ、-100〜+100）

	// 分析処理に関するメタデータ
	AnalysisDate    time.Time `json:"analysis_date" gorm:"not null"`                       // 分析実行日時
	AnalyzedAnswers int       `json:"analyzed_answers" gorm:"type:int;not null;default:0"` // 分析対象となった回答数
	DataPeriodStart time.Time `json:"data_period_start" gorm:"not null"`                   // 分析対象データの期間開始日
	DataPeriodEnd   time.Time `json:"data_period_end" gorm:"not null"`                     // 分析対象データの期間終了日
	LLMModel        string    `json:"llm_model" gorm:"type:varchar(50);not null"`          // 使用したLLMモデル名（例: claude-3-sonnet-20240229）
	AnalysisVersion string    `json:"analysis_version" gorm:"type:varchar(10);not null"`   // 分析ロジックのバージョン（プロンプトや処理の変更管理用）

	// 標準的なデータベース管理フィールド
	CreatedAt time.Time      `json:"created_at" gorm:"not null"`               // レコード作成日時
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null"`               // レコード最終更新日時
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`                  // 論理削除日時（ソフトデリート用）
	DeletedBy string         `json:"deleted_by" gorm:"type:char(36)"`          // 削除実行者のユーザーID
	CreatedBy string         `json:"created_by" gorm:"type:char(36);not null"` // レコード作成者のユーザーID
	UpdatedBy string         `json:"updated_by" gorm:"type:char(36);not null"` // レコード最終更新者のユーザーID

	// リレーション（関連テーブルとの関係）
	CategoryAnalyses []WeaknessCategoryAnalysis `json:"category_analyses,omitempty" gorm:"foreignKey:AnalysisID"` // カテゴリ別分析結果
	DetailedAnalysis *WeaknessDetailedAnalysis  `json:"detailed_analysis,omitempty" gorm:"foreignKey:AnalysisID"` // 詳細分析結果
	LearningAdvice   *WeaknessLearningAdvice    `json:"learning_advice,omitempty" gorm:"foreignKey:AnalysisID"`   // 学習アドバイス
}
