package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/service"
)

type QuestionMastersHandler struct {
	questionMastersService service.QuestionTemplateMastersService
}

func NewQuestionMastersHandler(questionMastersService service.QuestionTemplateMastersService) *QuestionMastersHandler {
	return &QuestionMastersHandler{
		questionMastersService: questionMastersService,
	}
}

// GetQuestionMasters 質問マスター一覧を取得するハンドラー
func (h *QuestionMastersHandler) GetQuestionMasters(c *gin.Context) {

	// コンテキストからユーザーIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	// JSONボディからリクエストを取得
	var req model.GetQuestionTemplateMastersSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	// クエリパラメータからページネーション情報を取得（優先）
	if perPageStr := c.Query("per_page"); perPageStr != "" {
		if n, err := strconv.Atoi(perPageStr); err == nil && n > 0 {
			req.PerPage = n
		}
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if n, err := strconv.Atoi(pageStr); err == nil && n > 0 {
			req.Page = n
		}
	}

	// デフォルト値の設定
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PerPage <= 0 {
		req.PerPage = 10
	}

	response, err := h.questionMastersService.GetQuestionTemplateMasters(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetQuestionMasterByID 質問マスターをIDで取得するハンドラー
func (h *QuestionMastersHandler) GetQuestionMasterByID(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "質問マスターIDが必要です"})
		return
	}

	response, err := h.questionMastersService.GetQuestionTemplateMasterByID(userID.(string), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
