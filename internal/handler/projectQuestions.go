package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/service"
)

type ProjectQuestionsHandler struct {
	projectQuestionsService service.ProjectQuestionsService
}

func NewProjectQuestionsHandler(projectQuestionsService service.ProjectQuestionsService) *ProjectQuestionsHandler {
	return &ProjectQuestionsHandler{
		projectQuestionsService: projectQuestionsService,
	}
}

// CreateProjectQuestions プロジェクト質問を作成するハンドラー
func (h *ProjectQuestionsHandler) CreateProjectQuestions(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var req model.CreateProjectQuestionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	// プロジェクト質問作成
	response, err := h.projectQuestionsService.CreateProjectQuestions(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetProjectQuestions プロジェクト質問を取得するハンドラー
func (h *ProjectQuestionsHandler) GetProjectQuestions(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}
	

	var req model.GetProjectQuestionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	response, err := h.projectQuestionsService.GetProjectQuestions(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}