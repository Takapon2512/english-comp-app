package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/service"
)

type ProjectHandler struct {
	projectService service.ProjectService
}

func NewProjectHandler(projectService service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// CreateProject プロジェクトを作成するハンドラー
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var req model.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	// プロジェクト作成
	response, err := h.projectService.CreateProject(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetProjects プロジェクト一覧を取得するハンドラー
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	// クエリパラメータを直接取得
	perPage := c.DefaultQuery("per_page", "10")
	page := c.DefaultQuery("page", "1")

	var req model.GetProjectsRequest
	// 手動でパラメータを設定
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
	req.UserID = userID.(string)

	response, err := h.projectService.GetProjects(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProjectDetail プロジェクト詳細を取得するハンドラー
func (h *ProjectHandler) GetProjectDetail(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var req model.GetProjectDetailRequest
	req.ID = c.Param("id")
	if req.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "プロジェクトIDが必要です"})
		return
	}

	response, err := h.projectService.GetProjectDetail(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
