package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/service"
)

type UserTagsHandler struct {
	userTagsService service.UserTagsService
}

func NewUserTagsHandler(userTagsService service.UserTagsService) *UserTagsHandler {
	return &UserTagsHandler{
		userTagsService: userTagsService,
	}
}

// createUserTags ユーザータグを作成するハンドラー
func (h *UserTagsHandler) CreateUserTags(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var req model.CreateUserTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	// ユーザータグ作成
	response, err := h.userTagsService.CreateUserTags(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// getUserTags ユーザータグを取得するハンドラー
func (h *UserTagsHandler) GetUserTags(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	response, err := h.userTagsService.GetUserTags(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}