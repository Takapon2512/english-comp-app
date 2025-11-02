package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/service"
)

type QuestionAnswersHandler struct {
	questionAnswersService service.QuestionAnswersService
}

func NewQuestionAnswersHandler(questionAnswersService service.QuestionAnswersService) *QuestionAnswersHandler {
	return &QuestionAnswersHandler{
		questionAnswersService: questionAnswersService,
	}
}

// CreateQuestionAnswers 質問回答を作成するハンドラー
func (h *QuestionAnswersHandler) CreateQuestionAnswers(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}
	
	var req model.CreateQuestionAnswersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}
	
	
	response, err := h.questionAnswersService.CreateQuestionAnswers(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	
	c.JSON(http.StatusCreated, response)
}

// GetQuestionAnswersByProjectID プロジェクトIDに紐づく質問回答を取得するハンドラー
func (h *QuestionAnswersHandler) GetQuestionAnswersByProjectID(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	projectID := c.Param("project_id")
	response, err := h.questionAnswersService.GetQuestionAnswersByProjectID(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateQuestionAnswersFinish プロジェクトIDに紐づく質問回答を更新するハンドラー
func (h *QuestionAnswersHandler) UpdateQuestionAnswersFinish(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	projectID := c.Param("project_id")
	response, err := h.questionAnswersService.UpdateQuestionAnswersFinish(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProjectQuestionToAnswer プロジェクトに紐づく問題の中から1題ランダムに取得するハンドラー
func (h *QuestionAnswersHandler) GetProjectQuestionToAnswer(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	projectID := c.Param("project_id")
	response, err := h.questionAnswersService.GetProjectQuestionToAnswer(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}