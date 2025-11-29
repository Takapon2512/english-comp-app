package model

import (
	"time"

	"gorm.io/gorm"
)

// WeaknessLearningAdvice は個別化された学習アドバイスと提案を管理するテーブル
// LLMが生成したパーソナライズされた学習支援情報を格納
type WeaknessLearningAdvice struct {
	// 基本識別情報
	ID         string `json:"id" gorm:"primaryKey;type:char(36)"`        // レコードの一意識別子
	AnalysisID string `json:"analysis_id" gorm:"type:char(36);not null"` // 親分析レコードのID（外部キー）

	// 学習アドバイス情報
	LearningAdvice      string `json:"learning_advice" gorm:"type:text"`      // 個別学習アドバイス（具体的な学習方法や注意点）
	RecommendedActions  string `json:"recommended_actions" gorm:"type:json"`  // 推奨アクションのJSON配列（具体的な学習行動の提案）
	NextGoals           string `json:"next_goals" gorm:"type:json"`           // 次の学習目標のJSON配列（短期・中期目標の設定）
	StudyPlan           string `json:"study_plan" gorm:"type:text"`           // 個別学習プラン（期間、内容、方法を含む詳細プラン）
	MotivationalMessage string `json:"motivational_message" gorm:"type:text"` // モチベーション向上メッセージ（励ましや成長の認識）

	// 標準的なデータベース管理フィールド
	CreatedAt time.Time      `json:"created_at" gorm:"not null"`               // レコード作成日時
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null"`               // レコード最終更新日時
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`                  // 論理削除日時
	CreatedBy string         `json:"created_by" gorm:"type:char(36);not null"` // レコード作成者のユーザーID
	UpdatedBy string         `json:"updated_by" gorm:"type:char(36);not null"` // レコード最終更新者のユーザーID
}

// TableName GORMのテーブル名を明示的に指定
func (WeaknessLearningAdvice) TableName() string {
	return "weakness_learning_advice"
}

// ===== 学習アドバイス用のサマリー構造体 =====

// PersonalizedAdvice は個別化された学習支援情報を表す構造体
// LLMが学習者の特性に合わせて生成したカスタマイズされたアドバイス
type PersonalizedAdvice struct {
	LearningAdvice      string   `json:"learning_advice"`      // 個別学習アドバイス（具体的な学習方法）
	RecommendedActions  []string `json:"recommended_actions"`  // 推奨する具体的な学習行動の配列
	NextGoals           []string `json:"next_goals"`           // 次に設定すべき学習目標の配列
	StudyPlan           string   `json:"study_plan"`           // 詳細な個別学習プラン（期間・内容・方法）
	MotivationalMessage string   `json:"motivational_message"` // 学習者を励ますパーソナライズされたメッセージ
}

// CreateWeaknessLearningAdviceRequest は学習アドバイスを作成するリクエスト用構造体
type CreateWeaknessLearningAdviceRequest struct {
	AnalysisID            string `json:"analysis_id" binding:"required"`   // 親分析レコードのID（必須）
	LearningAdvice        string `json:"learning_advice"` // 個別学習アドバイス（具体的な学習方法や注意点）
	RecommendedActions    string `json:"recommended_actions"` // 推奨アクションのJSON配列（具体的な学習行動の提案）
	NextGoals             string `json:"next_goals"` // 次の学習目標のJSON配列（短期・中期目標の設定）
	StudyPlan             string `json:"study_plan"` // 個別学習プラン（期間、内容、方法を含む詳細プラン）
	MotivationalMessage   string `json:"motivational_message"` // モチベーション向上メッセージ（励ましや成長の認識）
}

// CreateWeaknessLearningAdviceResponse は学習アドバイスを作成するレスポンス用構造体
type CreateWeaknessLearningAdviceResponse struct {
	ID 						string `json:"id"` // 作成された学習アドバイスレコードのID
	AnalysisID 				string `json:"analysis_id"` // 親分析レコードのID
	LearningAdvice 			string `json:"learning_advice"` // 個別学習アドバイス（具体的な学習方法や注意点）
	RecommendedActions 		string `json:"recommended_actions"` // 推奨アクションのJSON配列（具体的な学習行動の提案）
	NextGoals 				string `json:"next_goals"` // 次の学習目標のJSON配列（短期・中期目標の設定）
	StudyPlan 				string `json:"study_plan"` // 個別学習プラン（期間、内容、方法を含む詳細プラン）
	MotivationalMessage 	string `json:"motivational_message"` // モチベーション向上メッセージ（励ましや成長の認識）
}

// LLMが生成した個別学習アドバイスを表す構造体
type LLMWeaknessLearningAdviceRequest struct {
	GrammarScore          int    `json:"grammar_score"`          // 文法スコア（0-100）
	GrammarDescription    string `json:"grammar_description"`    // 文法面の詳細分析説明
	VocabularyScore       int    `json:"vocabulary_score"`       // 語彙スコア（0-100）
	VocabularyDescription string `json:"vocabulary_description"` // 語彙面の詳細分析説明
	ExpressionScore       int    `json:"expression_score"`       // 表現スコア（0-100）
	ExpressionDescription string `json:"expression_description"` // 表現面の詳細分析説明
}
