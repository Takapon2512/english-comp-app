package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/handler"
	"github.com/Takanpon2512/english-app/internal/middleware"
	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/repository"
	"github.com/Takanpon2512/english-app/internal/service"
)

func main() {
	// .envファイルの読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}
	// データベース接続設定
	dbUser := getEnvOrDefault("DB_USER", "root")
	dbPass := getEnvOrDefault("DB_PASSWORD", "password")
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "3306")
	dbName := getEnvOrDefault("DB_NAME", "english_app")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// データベース接続
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("データベース接続に失敗しました:", err)
	}

	// マイグレーション
	err = db.AutoMigrate(&model.User{}, &model.RefreshToken{})
	if err != nil {
		log.Fatal("マイグレーションに失敗しました:", err)
	}

	// Ginの初期化
	r := gin.Default()

	// CORS設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // NextJSのデフォルトポート
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ヘルスチェックエンドポイント
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// 環境変数から秘密鍵を取得
	secretKey := getEnvOrDefault("JWT_SECRET_KEY", "your-secret-key")

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	userTagsRepo := repository.NewUserTagsRepository(db)
	categoryMastersRepo := repository.NewCategoryMastersRepository(db)
	questionTemplateMastersRepo := repository.NewQuestionTemplateMastersRepository(db)
	projectQuestionsRepo := repository.NewProjectQuestionsRepository(db)
	questionAnswersRepo := repository.NewQuestionAnswersRepository(db)
	correctResultsRepo := repository.NewCorrectResultsRepository(db)

	// サービスの初期化
	authService := service.NewAuthService(userRepo)
	projectService := service.NewProjectService(db, projectRepo)
	userTagsService := service.NewUserTagsService(db, userTagsRepo)
	categoryMastersService := service.NewCategoryMastersService(db, categoryMastersRepo)
	questionTemplateMastersService := service.NewQuestionTemplateMastersService(db, questionTemplateMastersRepo)
	projectQuestionsService := service.NewProjectQuestionsService(db, projectQuestionsRepo, questionTemplateMastersRepo)
	questionAnswersService := service.NewQuestionAnswersService(db, questionAnswersRepo, projectQuestionsRepo, questionTemplateMastersRepo)
	correctResultsService := service.NewCorrectResultsService(db, correctResultsRepo, questionTemplateMastersRepo, questionAnswersRepo, categoryMastersRepo)

	// ハンドラーの初期化
	authHandler := handler.NewAuthHandler(authService, secretKey)
	projectHandler := handler.NewProjectHandler(projectService)
	userTagsHandler := handler.NewUserTagsHandler(userTagsService)
	categoryMastersHandler := handler.NewCategoryMastersHandler(categoryMastersService)
	questionTemplateMastersHandler := handler.NewQuestionMastersHandler(questionTemplateMastersService)
	projectQuestionsHandler := handler.NewProjectQuestionsHandler(projectQuestionsService)
	questionAnswersHandler := handler.NewQuestionAnswersHandler(questionAnswersService)
	correctResultsHandler := handler.NewCorrectResultsHandler(correctResultsService)

	// 認証ミドルウェアの初期化
	authMiddleware := middleware.NewAuthMiddleware(middleware.AuthConfig{
		SecretKey: secretKey,
	})

	// 認証不要のエンドポイント
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/signup", authHandler.Signup)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
	}

	// 認証が必要なエンドポイント
	api := r.Group("/api/v1")
	api.Use(authMiddleware)
	{
		// 認証が必要なエンドポイントをここに追加
		api.GET("/user", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			email, _ := c.Get("email")
			c.JSON(200, gin.H{
				"user_id": userID,
				"email":   email,
			})
		})

		api.POST("/projects", projectHandler.CreateProject)
		api.GET("/projects", projectHandler.GetProjects)
		api.GET("/projects/:id", projectHandler.GetProjectDetail)
		api.POST("/projects/create-questions", projectQuestionsHandler.CreateProjectQuestions)
		api.POST("/projects/questions", projectQuestionsHandler.GetProjectQuestions)

		api.POST("/user-tags", userTagsHandler.CreateUserTags)
		api.GET("/user-tags", userTagsHandler.GetUserTags)
		api.PUT("/user-tags/update", userTagsHandler.UpdateUserTags)
		api.PUT("/user-tags/delete", userTagsHandler.DeleteUserTags)

		api.GET("/category-masters", categoryMastersHandler.GetCategoryMasters)
		// api.GET("/category-master", categoryMastersHandler.GetCategoryMastersByID)

		api.POST("/question-masters", questionTemplateMastersHandler.GetQuestionMasters)
		api.GET("/question-masters/:id", questionTemplateMastersHandler.GetQuestionMasterByID)

		api.POST("/question-answers", questionAnswersHandler.CreateQuestionAnswers)
		api.GET("/question-answers/:project_id", questionAnswersHandler.GetQuestionAnswersByProjectID)
		api.PUT("/question-answers/finish/:project_id", questionAnswersHandler.UpdateQuestionAnswersFinish)
		api.POST("/question-answers/question-to-answer/:project_id", questionAnswersHandler.GetProjectQuestionToAnswer)

		api.POST("/correct-results", correctResultsHandler.CreateCorrectResult)
		api.POST("/correct-results/get", correctResultsHandler.GetCorrectResults)
		api.POST("/correct-results/version-list", correctResultsHandler.GetCorrectResultsVersionList)
	}

	// サーバーの起動
	port := getEnvOrDefault("PORT", "8080")
	if err := r.Run(":" + port); err != nil {
		log.Fatal("サーバーの起動に失敗しました:", err)
	}
}

// 環境変数を取得、未設定の場合はデフォルト値を返す
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
