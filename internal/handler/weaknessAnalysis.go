package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/service"
)

type WeaknessAnalysisHandler struct {
  weaknessAnalysisService service.WeaknessAnalysisService
}

func NewWeaknessAnalysisHandler(weaknessAnalysisService service.WeaknessAnalysisService) *WeaknessAnalysisHandler {
  return &WeaknessAnalysisHandler{
    weaknessAnalysisService: weaknessAnalysisService,
  }
}

// CreateWeaknessAnalysis 学習弱点分析を作成するハンドラー
func (h *WeaknessAnalysisHandler) CreateWeaknessAnalysis(c *gin.Context) {
  // コンテキストからユーザーIDを取得
  userId, exists := c.Get("user_id")
  if !exists {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
    return
  }

  var req model.CreateWeaknessAnalysisRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
    return
  }

  response, err := h.weaknessAnalysisService.CreateWeaknessAnalysis(userId.(string), &req)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.JSON(http.StatusCreated, response)
}

// GetWeaknessAnalysisAllSummary 学習弱点分析の全ての結果を取得するハンドラー
func (h *WeaknessAnalysisHandler) GetWeaknessAnalysisAllSummary(c *gin.Context) {
  // コンテキストからユーザーIDを取得
  userId, exists := c.Get("user_id")

  if !exists {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
    return
  }

  projectID := c.Param("project_id")
  if projectID == "" {
    c.JSON(http.StatusBadRequest, gin.H{"error": "プロジェクトIDが必要です"})
    return
  }

  response, err := h.weaknessAnalysisService.GetWeaknessAnalysisAllSummary(userId.(string), projectID)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.JSON(http.StatusOK, response)
}
