package model

import "time"

// ===== リクエスト・レスポンス構造体 =====

// CreateWeaknessAnalysisRequest はLLM分析の開始リクエスト用構造体
// プロジェクトIDを指定して、そのプロジェクト内の回答データを分析対象とする
type CreateWeaknessAnalysisRequest struct {
	ProjectID string `json:"project_id" binding:"required"` // 分析対象プロジェクトのID（必須）
}

// CreateWeaknessAnalysisResponse はLLM分析開始のレスポンス用構造体
// 分析処理は非同期で実行されるため、処理状況を確認するためのIDを返す
type CreateWeaknessAnalysisResponse struct {
	ID             string `json:"id"`              // 作成された分析レコードのID
	ProjectID      string `json:"project_id"`      // 分析対象プロジェクトのID
	AnalysisStatus string `json:"analysis_status"` // 分析処理の初期状態（通常は"PROCESSING"）
}

// GetWeaknessAnalysisRequest は分析結果取得リクエスト用構造体
// プロジェクトIDを指定して、最新の分析結果を取得する
type GetWeaknessAnalysisRequest struct {
	ProjectID string `json:"project_id" binding:"required"` // 分析結果を取得したいプロジェクトのID（必須）
}

// GetWeaknessAnalysisResponse は分析結果取得のレスポンス用構造体
// クライアントに返す最終的な分析結果のラッパー
type GetWeaknessAnalysisResponse struct {
	Analysis WeaknessAnalysisSummary `json:"analysis"` // 構造化された分析結果のサマリー
}

// ===== サマリー構造体（クライアント向け） =====

// WeaknessAnalysisSummary はクライアントに返す分析結果のサマリー構造体
// データベースのJSON フィールドを構造化されたオブジェクトに変換して提供
type WeaknessAnalysisSummary struct {
	ID                 string                 `json:"id"`                  // 分析結果のID
	ProjectID          string                 `json:"project_id"`          // 分析対象プロジェクトのID
	AnalysisStatus     string                 `json:"analysis_status"`     // 分析処理状況
	OverallScore       int                    `json:"overall_score"`       // 総合学習スコア（0-100）
	ImprovementRate    int                    `json:"improvement_rate"`    // 前回からの改善率（-100〜+100）
	WeakCategories     []CategoryWeakness     `json:"weak_categories"`     // 苦手カテゴリの詳細配列
	StrengthCategories []CategoryStrength     `json:"strength_categories"` // 得意カテゴリの詳細配列
	DetailedAnalysis   DetailedAnalysisResult `json:"detailed_analysis"`   // 4領域（文法・語彙・表現・構成）の詳細分析
	PersonalizedAdvice PersonalizedAdvice     `json:"personalized_advice"` // 個別化された学習アドバイス
	AnalysisDate       time.Time              `json:"analysis_date"`       // 分析実行日時
	AnalyzedAnswers    int                    `json:"analyzed_answers"`    // 分析対象回答数
	DataPeriodStart    time.Time              `json:"data_period_start"`   // 分析対象期間の開始日
	DataPeriodEnd      time.Time              `json:"data_period_end"`     // 分析対象期間の終了日
}
