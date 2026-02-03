package handlers

import (
	"github.com/gin-gonic/gin"
	"go-gin-todo-api/database"
	"go-gin-todo-api/middleware"
	"go-gin-todo-api/models"
	"go-gin-todo-api/utils"
)

// GetMe 自分の情報を取得
func GetMe(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		utils.RespondUnauthorized(c, "Unauthorized")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		statusCode, message := utils.HandleDBError(err)
		if statusCode == 404 {
			utils.RespondNotFound(c, "User not found")
		} else {
			utils.RespondError(c, statusCode, utils.ErrorCodeInternal, message)
		}
		return
	}

	c.JSON(200, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}
