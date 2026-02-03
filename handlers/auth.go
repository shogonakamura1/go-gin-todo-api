package handlers

import (
	"net/http"
	"time"

	"go-gin-todo-api/database"
	"go-gin-todo-api/models"
	"go-gin-todo-api/utils"

	"github.com/gin-gonic/gin"
)

// RegisterRequest ユーザー登録リクエスト
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5"`
}

// LoginRequest ログインリクエスト
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshRequest リフレッシュトークンリクエスト
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Register ユーザー登録ハンドラー
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, err.Error())
		return
	}

	// パスワードをハッシュ化
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.RespondInternalError(c, "Failed to hash password")
		return
	}

	// ユーザーを作成
	user := models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		statusCode, message := utils.HandleDBError(err)
		if statusCode == 409 {
			utils.RespondConflict(c, "Email already exists")
		} else {
			utils.RespondError(c, statusCode, utils.ErrorCodeInternal, message)
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}

// Login ログインハンドラー
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, err.Error())
		return
	}

	// ユーザーを検索
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.RespondUnauthorized(c, "Invalid email or password")
		return
	}

	// パスワードを検証
	if !utils.ComparePassword(user.PasswordHash, req.Password) {
		utils.RespondUnauthorized(c, "Invalid email or password")
		return
	}

	// アクセストークンを生成
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		utils.RespondInternalError(c, "Failed to generate access token")
		return
	}

	// リフレッシュトークンを生成
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.RespondInternalError(c, "Failed to generate refresh token")
		return
	}

	// リフレッシュトークンをハッシュ化してDBに保存
	tokenHash := utils.HashRefreshToken(refreshToken)
	expiresAt := time.Now().Add(utils.GetRefreshTokenTTL())

	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	if err := database.DB.Create(&refreshTokenModel).Error; err != nil {
		statusCode, message := utils.HandleDBError(err)
		utils.RespondError(c, statusCode, utils.ErrorCodeInternal, message)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
	})
}

// Refresh リフレッシュトークンハンドラー
func Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, err.Error())
		return
	}

	// リフレッシュトークンをハッシュ化
	tokenHash := utils.HashRefreshToken(req.RefreshToken)

	// DBでトークンを検索
	var refreshToken models.RefreshToken
	if err := database.DB.Where("token_hash = ?", tokenHash).
		Where("revoked_at IS NULL").
		Where("expires_at > ?", time.Now()).
		First(&refreshToken).Error; err != nil {
		utils.RespondUnauthorized(c, "Invalid or expired refresh token")
		return
	}

	// 新しいアクセストークンを生成
	accessToken, err := utils.GenerateAccessToken(refreshToken.UserID)
	if err != nil {
		utils.RespondInternalError(c, "Failed to generate access token")
		return
	}

	// ローテーション: 新しいリフレッシュトークンを生成
	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.RespondInternalError(c, "Failed to generate refresh token")
		return
	}

	// 古いトークンをrevoke
	now := time.Now()
	database.DB.Model(&refreshToken).Update("revoked_at", now)

	// 新しいリフレッシュトークンをDBに保存
	newTokenHash := utils.HashRefreshToken(newRefreshToken)
	newExpiresAt := time.Now().Add(utils.GetRefreshTokenTTL())

	newRefreshTokenModel := models.RefreshToken{
		UserID:    refreshToken.UserID,
		TokenHash: newTokenHash,
		ExpiresAt: newExpiresAt,
	}

	if err := database.DB.Create(&newRefreshTokenModel).Error; err != nil {
		statusCode, message := utils.HandleDBError(err)
		utils.RespondError(c, statusCode, utils.ErrorCodeInternal, message)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
		"token_type":    "Bearer",
	})
}

// Logout ログアウトハンドラー
func Logout(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, err.Error())
		return
	}

	// リフレッシュトークンをハッシュ化
	tokenHash := utils.HashRefreshToken(req.RefreshToken)

	// トークンをrevoke
	result := database.DB.Model(&models.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Where("revoked_at IS NULL").
		Update("revoked_at", time.Now())

	if result.Error != nil {
		statusCode, message := utils.HandleDBError(result.Error)
		utils.RespondError(c, statusCode, utils.ErrorCodeInternal, message)
		return
	}

	if result.RowsAffected == 0 {
		utils.RespondNotFound(c, "Token not found or already revoked")
		return
	}

	c.Status(http.StatusNoContent)
}
