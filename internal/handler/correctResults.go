package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/service"
)

type CorrectResultsHandler struct {
	correctResultsService service.CorrectResultsService
}

func NewCorrectResultsHandler(correctResultsService service.CorrectResultsService) *CorrectResultsHandler {
	return &CorrectResultsHandler{
		correctResultsService: correctResultsService,
	}
}

// CreateCorrectResult 添削結果を作成するハンドラー
func (h *CorrectResultsHandler) CreateCorrectResult(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var reqCreate model.CreateCorrectionResultRequest
	if err := c.ShouldBindJSON(&reqCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	resCreate, err := h.correctResultsService.CreateCorrectionResult(userID.(string), &reqCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var reqGrand model.GrandCorrectResultRequest
	reqGrand.ID = resCreate.ID

	resGrand, err := h.correctResultsService.GrandCorrectResult(userID.(string), &reqGrand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resGrand)
}

// GetCorrectResults 添削結果を取得するハンドラー
func (h *CorrectResultsHandler) GetCorrectResults(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var reqPost model.GetCorrectResultsRequest
	if err := c.ShouldBindJSON(&reqPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	resPost, err := h.correctResultsService.GetCorrectResults(userID.(string), &reqPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resPost)
}

// GetCorrectResultsVersionList 添削結果のバージョン一覧を取得するハンドラー
func (h *CorrectResultsHandler) GetCorrectResultsVersionList(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var reqVersionList model.GetCorrectResultsVersionRequest
	if err := c.ShouldBindJSON(&reqVersionList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	resVersionList, err := h.correctResultsService.GetCorrectResultsVersionList(userID.(string), &reqVersionList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resVersionList)
}