package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// データベース接続
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/english_comp_app?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("データベース接続に失敗しました:", err)
	}

	// マイグレーション実行
	if err := db.AutoMigrate(&model.CategoryMasters{}, &model.QuestionTemplateMasters{}); err != nil {
		log.Fatal("マイグレーションに失敗しました:", err)
	}

	// Seeder実行
	categoryMap, err := seedCategoryMasters(db)
	if err != nil {
		log.Fatal("CategoryMasters seeder実行に失敗しました:", err)
	}

	if err := seedQuestionTemplateMasters(db, categoryMap); err != nil {
		log.Fatal("QuestionTemplateMasters seeder実行に失敗しました:", err)
	}

	fmt.Println("Seederが正常に完了しました")
}

func seedCategoryMasters(db *gorm.DB) (map[string]string, error) {
	// 既存データをクリア
	if err := db.Where("1 = 1").Delete(&model.CategoryMasters{}).Error; err != nil {
		return nil, fmt.Errorf("既存データの削除に失敗しました: %w", err)
	}

	// サンプルデータ
	categories := []model.CategoryMasters{
		{
			ID:        uuid.New().String(),
			Name:      "日常会話",
			CreatedBy: "system",
			UpdatedBy: "system",
		},
		{
			ID:        uuid.New().String(),
			Name:      "翻訳練習",
			CreatedBy: "system",
			UpdatedBy: "system",
		},
		{
			ID:        uuid.New().String(),
			Name:      "文法問題",
			CreatedBy: "system",
			UpdatedBy: "system",
		},
		{
			ID:        uuid.New().String(),
			Name:      "ビジネス英語",
			CreatedBy: "system",
			UpdatedBy: "system",
		},
		{
			ID:        uuid.New().String(),
			Name:      "ディスカッション",
			CreatedBy: "system",
			UpdatedBy: "system",
		},
	}

	// データベースに挿入
	if err := db.Create(&categories).Error; err != nil {
		return nil, fmt.Errorf("データの挿入に失敗しました: %w", err)
	}

	// カテゴリ名とIDのマッピングを作成
	categoryMap := make(map[string]string)
	for _, category := range categories {
		categoryMap[category.Name] = category.ID
	}

	fmt.Printf("CategoryMasters seeder: %d件のデータを挿入しました\n", len(categories))
	return categoryMap, nil
}

func seedQuestionTemplateMasters(db *gorm.DB, categoryMap map[string]string) error {
	// 既存データをクリア
	if err := db.Where("1 = 1").Delete(&model.QuestionTemplateMasters{}).Error; err != nil {
		return fmt.Errorf("既存データの削除に失敗しました: %w", err)
	}

	// サンプルデータ
	questionTemplates := []model.QuestionTemplateMasters{
		{
			ID:            uuid.New().String(),
			CategoryID:    categoryMap["日常会話"],
			QuestionType:  "essay",
			English:       "Write about your favorite hobby and explain why you enjoy it.",
			Japanese:      "あなたの好きな趣味について書き、なぜそれを楽しんでいるのか説明してください。",
			Status:        "ACTIVE",
			Level:         "basic",
			EstimatedTime: 15,
			Points:        10,
			CreatedBy:     "system",
			UpdatedBy:     "system",
		},
		{
			ID:            uuid.New().String(),
			CategoryID:    categoryMap["日常会話"],
			QuestionType:  "essay",
			English:       "Describe your ideal vacation destination and explain what activities you would do there.",
			Japanese:      "理想的な休暇の目的地を説明し、そこでどのような活動をするか説明してください。",
			Status:        "ACTIVE",
			Level:         "inter",
			EstimatedTime: 20,
			Points:        15,
			CreatedBy:     "system",
			UpdatedBy:     "system",
		},
		{
			ID:            uuid.New().String(),
			CategoryID:    categoryMap["翻訳練習"],
			QuestionType:  "translate",
			English:       "Translate the following Japanese sentence to English: 私は昨日、友達と映画を見に行きました。",
			Japanese:      "次の日本語の文を英語に翻訳してください：私は昨日、友達と映画を見に行きました。",
			Status:        "ACTIVE",
			Level:         "basic",
			EstimatedTime: 10,
			Points:        8,
			CreatedBy:     "system",
			UpdatedBy:     "system",
		},
		{
			ID:            uuid.New().String(),
			CategoryID:    categoryMap["翻訳練習"],
			QuestionType:  "translate",
			English:       "Translate the following English sentence to Japanese: The weather is beautiful today, so I decided to go for a walk in the park.",
			Japanese:      "次の英語の文を日本語に翻訳してください：The weather is beautiful today, so I decided to go for a walk in the park.",
			Status:        "ACTIVE",
			Level:         "inter",
			EstimatedTime: 12,
			Points:        12,
			CreatedBy:     "system",
			UpdatedBy:     "system",
		},
		{
			ID:            uuid.New().String(),
			CategoryID:    categoryMap["文法問題"],
			QuestionType:  "fill",
			English:       "Complete the sentence: I have been studying English _____ three years.",
			Japanese:      "文を完成させてください：I have been studying English _____ three years.",
			Status:        "ACTIVE",
			Level:         "inter",
			EstimatedTime: 5,
			Points:        5,
			CreatedBy:     "system",
			UpdatedBy:     "system",
		},
		{
			ID:            uuid.New().String(),
			CategoryID:    categoryMap["文法問題"],
			QuestionType:  "fill",
			English:       "Complete the sentence: If I _____ (have) more time, I would travel around the world.",
			Japanese:      "文を完成させてください：If I _____ (have) more time, I would travel around the world.",
			Status:        "ACTIVE",
			Level:         "inter",
			EstimatedTime: 8,
			Points:        8,
			CreatedBy:     "system",
			UpdatedBy:     "system",
		},
		{
			ID:            uuid.New().String(),
			CategoryID:    categoryMap["文法問題"],
			QuestionType:  "fill",
			English:       "Complete the sentence: The book _____ (write) by a famous author last year.",
			Japanese:      "文を完成させてください：The book _____ (write) by a famous author last year.",
			Status:        "ACTIVE",
			Level:         "inter",
			EstimatedTime: 6,
			Points:        6,
			CreatedBy:     "system",
			UpdatedBy:     "system",
		},
		{
			ID:            uuid.New().String(),
			CategoryID:    categoryMap["文法問題"],
			QuestionType:  "fill",
			English:       "Complete the sentence: By the time I arrived, the meeting _____ (already/start).",
			Japanese:      "文を完成させてください：By the time I arrived, the meeting _____ (already/start).",
			Status:        "ACTIVE",
			Level:         "adv",
			EstimatedTime: 10,
			Points:        10,
			CreatedBy:     "system",
			UpdatedBy:     "system",
		},
		{
			ID:            uuid.New().String(),
			CategoryID:    categoryMap["翻訳練習"],
			QuestionType:  "translate",
			English:       "Translate the following Japanese sentence to English: もし時間があれば、世界中を旅行したいと思います。",
			Japanese:      "次の日本語の文を英語に翻訳してください：もし時間があれば、世界中を旅行したいと思います。",
			Status:        "ACTIVE",
			Level:         "inter",
			EstimatedTime: 15,
			Points:        15,
			CreatedBy:     "system",
			UpdatedBy:     "system",
		},
		{
			ID:            uuid.New().String(),
			CategoryID:    categoryMap["翻訳練習"],
			QuestionType:  "translate",
			English:       "Translate the following English sentence to Japanese: If I had studied harder, I would have passed the exam.",
			Japanese:      "次の英語の文を日本語に翻訳してください：If I had studied harder, I would have passed the exam.",
			Status:        "ACTIVE",
			Level:         "adv",
			EstimatedTime: 18,
			Points:        18,
			CreatedBy:     "system",
			UpdatedBy:     "system",
		},
		{
			ID:            uuid.New().String(),
			CategoryID:    categoryMap["ディスカッション"],
			QuestionType:  "essay",
			English:       "Discuss the advantages and disadvantages of social media in modern society.",
			Japanese:      "現代社会におけるソーシャルメディアの利点と欠点について議論してください。",
			Status:        "ACTIVE",
			Level:         "adv",
			EstimatedTime: 30,
			Points:        25,
			CreatedBy:     "system",
			UpdatedBy:     "system",
		},
		{
			ID:            uuid.New().String(),
			CategoryID:    categoryMap["ディスカッション"],
			QuestionType:  "essay",
			English:       "Explain your opinion on the impact of technology on education.",
			Japanese:      "教育におけるテクノロジーの影響についてあなたの意見を説明してください。",
			Status:        "ACTIVE",
			Level:         "adv",
			EstimatedTime: 25,
			Points:        20,
			CreatedBy:     "system",
			UpdatedBy:     "system",
		},
	}

	// データベースに挿入
	if err := db.Create(&questionTemplates).Error; err != nil {
		return fmt.Errorf("データの挿入に失敗しました: %w", err)
	}

	fmt.Printf("QuestionTemplateMasters seeder: %d件のデータを挿入しました\n", len(questionTemplates))
	return nil
}
