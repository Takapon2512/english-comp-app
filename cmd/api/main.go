package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	// サービスの初期化
	authService := service.NewAuthService(userRepo)
	projectService := service.NewProjectService(db, projectRepo)
	userTagsService := service.NewUserTagsService(db, userTagsRepo)

	// ハンドラーの初期化
	authHandler := handler.NewAuthHandler(authService, secretKey)
	projectHandler := handler.NewProjectHandler(projectService)
	userTagsHandler := handler.NewUserTagsHandler(userTagsService)

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

		api.POST("/user-tags", userTagsHandler.CreateUserTags)
		api.GET("/user-tags", userTagsHandler.GetUserTags)
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
