package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/service"
)

type CategoryMastersHandler struct {
	categoryMastersService service.CategoryMastersService
}

func NewCategoryMastersHandler(categoryMastersService service.CategoryMastersService) *CategoryMastersHandler {
	return &CategoryMastersHandler{
		categoryMastersService: categoryMastersService,
	}
}

// GetCategoryMasters カテゴリーマスター一覧を取得するハンドラー
func (h *CategoryMastersHandler) GetCategoryMasters(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	// クエリパラメータを直接取得
	perPage := c.DefaultQuery("per_page", "10")
	page := c.DefaultQuery("page", "1")

	var req model.GetCategoryMastersSearchRequest

	if perPage != "" {
		if n, err := strconv.Atoi(perPage); err == nil && n > 0 {
			req.PerPage = n
		}
	}

	if page != "" {
		if n, err := strconv.Atoi(page); err == nil && n > 0 {
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

	response, err := h.categoryMastersService.GetCategoryMasters(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}