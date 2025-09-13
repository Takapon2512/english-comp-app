package handler

import (
	"net/http"
	"time"

	"github.com/Takanpon2512/english-app/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	authService service.AuthService
	secretKey   string
}

func NewAuthHandler(authService service.AuthService, secretKey string) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		secretKey:   secretKey,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
	TokenType    string       `json:"token_type"`
	User         ResponseUser `json:"user"`
}

// ResponseUser は、APIレスポンス用のユーザー情報を表す
type ResponseUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// JWTトークンの生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークンの生成に失敗しました"})
		return
	}

	// リフレッシュトークンの生成
	refreshToken, err := h.authService.CreateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "リフレッシュトークンの生成に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshToken.TokenHash,
		ExpiresIn:    24 * 60 * 60, // 24時間
		TokenType:    "Bearer",
		User: ResponseUser{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	})
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	user, err := h.authService.Signup(req.Email, req.Password, req.Name)
	if err != nil {
		if err == service.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー登録に失敗しました"})
		return
	}

	// JWTトークンの生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークンの生成に失敗しました"})
		return
	}

	// リフレッシュトークンの生成
	refreshToken, err := h.authService.CreateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "リフレッシュトークンの生成に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshToken.TokenHash,
		ExpiresIn:    24 * 60 * 60, // 24時間
		TokenType:    "Bearer",
		User: ResponseUser{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	// リフレッシュトークンを検証し、新しいものを生成
	user, newRefreshToken, err := h.authService.RotateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "無効なリフレッシュトークンです"})
		return
	}

	// 新しいアクセストークンの生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークンの生成に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  tokenString,
		RefreshToken: newRefreshToken.TokenHash, // 新しいリフレッシュトークンを返す
		ExpiresIn:    24 * 60 * 60,              // 24時間
		TokenType:    "Bearer",
		User: ResponseUser{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	err := h.authService.Logout(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ログアウトに失敗しました"})
		return
	}

	c.JSON(http.StatusOK, LogoutResponse{
		Message: "ログアウトしました",
	})
}