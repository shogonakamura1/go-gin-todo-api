package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go-gin-todo-api/utils"
)

const UserIDKey = "user_id"

// AuthMiddleware JWTトークンを検証し、user_idをcontextに設定するミドルウェア
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.RespondUnauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		// "Bearer " プレフィックスをチェック
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondUnauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		userID, err := utils.ValidateAccessToken(tokenString)
		if err != nil {
			utils.RespondUnauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// user_idをcontextに設定
		c.Set(UserIDKey, userID)
		c.Next()
	}
}
