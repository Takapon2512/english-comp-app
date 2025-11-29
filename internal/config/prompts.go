package config

import "fmt"

// WeaknessAnalysisPrompts 弱点分析用のプロンプト設定
type WeaknessAnalysisPrompts struct {
	CategoryAnalysis CategoryAnalysisPrompt
	DetailedAnalysis DetailedAnalysisPrompt
	LearningAdvice   LearningAdvicePrompt
}

// CategoryAnalysisPrompt カテゴリ分析用プロンプト
type CategoryAnalysisPrompt struct {
	Template string
}

// DetailedAnalysisPrompt 詳細分析用プロンプト
type DetailedAnalysisPrompt struct {
	Template string
}

// LearningAdvicePrompt 学習アドバイス用プロンプト
type LearningAdvicePrompt struct {
	Template string
}

// NewWeaknessAnalysisPrompts プロンプト設定を初期化
func NewWeaknessAnalysisPrompts() *WeaknessAnalysisPrompts {
	return &WeaknessAnalysisPrompts{
		CategoryAnalysis: CategoryAnalysisPrompt{
			Template: `あなたはプロの英語教師です。
以下の「%s」カテゴリの学習データを分析し、このカテゴリでの学習者の強み・弱みを分析してJSON形式で出力してください。

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
%s

上記データを分析し、有効なJSONのみを出力してください：`,
		},
		DetailedAnalysis: DetailedAnalysisPrompt{
			Template: `あなたはプロの英語教師です。
以下の学習データを分析し、学習者の英語力を4つの領域（文法・語彙・表現・構成）で詳細に分析してJSON形式で出力してください。

【重要】以下の要件を厳密に守ってください：
1. 出力は有効なJSON形式のみにしてください
2. 説明文やマークダウン記法は一切含めないでください
3. JSONの前後に余計な文字を入れないでください
4. 配列が空の場合は空配列[]を使用してください
5. スコアは0-100の整数で設定してください

出力JSON形式：
{
  "grammar": {
    "score": 整数(0-100),
    "description": "文法面の詳細分析説明",
    "examples": ["具体例1", "具体例2"]
  },
  "vocabulary": {
    "score": 整数(0-100),
    "description": "語彙面の詳細分析説明",
    "examples": ["具体例1", "具体例2"]
  },
  "expression": {
    "score": 整数(0-100),
    "description": "表現面の詳細分析説明",
    "examples": ["具体例1", "具体例2"]
  },
  "structure": {
    "score": 整数(0-100),
    "description": "構成面の詳細分析説明",
    "examples": ["具体例1", "具体例2"]
  }
}

分析対象データ:
%s

上記データを分析し、有効なJSONのみを出力してください：`,
		},
		LearningAdvice: LearningAdvicePrompt{
			Template: `あなたはプロの英語学習コーチです。
以下の詳細分析結果に基づいて、学習者に個別化された学習アドバイスを作成してJSON形式で出力してください。

【重要】以下の要件を厳密に守ってください：
1. 出力は有効なJSON形式のみにしてください
2. 説明文やマークダウン記法は一切含めないでください
3. JSONの前後に余計な文字を入れないでください
4. 配列が空の場合は空配列[]を使用してください
5. 学習者のレベルに応じた具体的で実践的なアドバイスを提供してください

出力JSON形式：
{
  "learning_advice": "個別学習アドバイス（具体的な学習方法や注意点）",
  "recommended_actions": ["推奨アクション1", "推奨アクション2", "推奨アクション3"],
  "next_goals": ["短期目標1", "中期目標1", "長期目標1"],
  "study_plan": "詳細な個別学習プラン（期間・内容・方法を含む）",
  "motivational_message": "学習者を励ますパーソナライズされたメッセージ"
}

詳細分析結果:
%s

上記分析結果に基づいて、学習者に最適化された学習アドバイスを有効なJSONのみで出力してください：`,
		},
	}
}

// GetCategoryAnalysisPrompt カテゴリ分析用プロンプトを取得
func (p *WeaknessAnalysisPrompts) GetCategoryAnalysisPrompt(categoryName, jsonData string) string {
	return fmt.Sprintf(p.CategoryAnalysis.Template, categoryName, jsonData)
}

// GetDetailedAnalysisPrompt 詳細分析用プロンプトを取得
func (p *WeaknessAnalysisPrompts) GetDetailedAnalysisPrompt(jsonData string) string {
	return fmt.Sprintf(p.DetailedAnalysis.Template, jsonData)
}

// GetLearningAdvicePrompt 学習アドバイス用プロンプトを取得
func (p *WeaknessAnalysisPrompts) GetLearningAdvicePrompt(jsonData string) string {
	return fmt.Sprintf(p.LearningAdvice.Template, jsonData)
}
