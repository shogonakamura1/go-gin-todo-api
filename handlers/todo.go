package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-gin-todo-api/database"
	"go-gin-todo-api/middleware"
	"go-gin-todo-api/models"
	"go-gin-todo-api/utils"
)

// CreateTodoRequest Todo作成リクエスト
type CreateTodoRequest struct {
	Title string `json:"title" binding:"required"`
}

// UpdateTodoRequest Todo更新リクエスト
type UpdateTodoRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

// GetTodos 自分のtodo一覧を取得
func GetTodos(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		utils.RespondUnauthorized(c, "Unauthorized")
		return
	}

	var todos []models.Todo
	if err := database.DB.Where("user_id = ?", userID).Find(&todos).Error; err != nil {
		statusCode, message := utils.HandleDBError(err)
		utils.RespondError(c, statusCode, utils.ErrorCodeInternal, message)
		return
	}

	c.JSON(http.StatusOK, todos)
}

// CreateTodo Todoを作成
func CreateTodo(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		utils.RespondUnauthorized(c, "Unauthorized")
		return
	}

	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, err.Error())
		return
	}

	todo := models.Todo{
		UserID:    userID.(uint),
		Title:     req.Title,
		Completed: false,
	}

	if err := database.DB.Create(&todo).Error; err != nil {
		statusCode, message := utils.HandleDBError(err)
		utils.RespondError(c, statusCode, utils.ErrorCodeInternal, message)
		return
	}

	c.JSON(http.StatusCreated, todo)
}

// GetTodo 特定のtodoを取得（自分のものだけ）
func GetTodo(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		utils.RespondUnauthorized(c, "Unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondBadRequest(c, "Invalid todo ID")
		return
	}

	var todo models.Todo
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		statusCode, message := utils.HandleDBError(err)
		if statusCode == 404 {
			utils.RespondNotFound(c, "Todo not found")
		} else {
			utils.RespondError(c, statusCode, utils.ErrorCodeInternal, message)
		}
		return
	}

	c.JSON(http.StatusOK, todo)
}

// UpdateTodo Todoを更新
func UpdateTodo(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		utils.RespondUnauthorized(c, "Unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondBadRequest(c, "Invalid todo ID")
		return
	}

	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, err.Error())
		return
	}

	// 自分のtodoか確認
	var todo models.Todo
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		statusCode, message := utils.HandleDBError(err)
		if statusCode == 404 {
			utils.RespondNotFound(c, "Todo not found")
		} else {
			utils.RespondError(c, statusCode, utils.ErrorCodeInternal, message)
		}
		return
	}

	// 更新フィールドを設定
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Completed != nil {
		updates["completed"] = *req.Completed
	}

	if len(updates) == 0 {
		utils.RespondBadRequest(c, "No fields to update")
		return
	}

	// 更新実行
	if err := database.DB.Model(&todo).Updates(updates).Error; err != nil {
		statusCode, message := utils.HandleDBError(err)
		utils.RespondError(c, statusCode, utils.ErrorCodeInternal, message)
		return
	}

	// 更新後のデータを取得
	database.DB.First(&todo, id)
	c.JSON(http.StatusOK, todo)
}

// DeleteTodo Todoを削除
func DeleteTodo(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		utils.RespondUnauthorized(c, "Unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondBadRequest(c, "Invalid todo ID")
		return
	}

	// 自分のtodoか確認して削除
	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Todo{})
	if result.Error != nil {
		statusCode, message := utils.HandleDBError(result.Error)
		utils.RespondError(c, statusCode, utils.ErrorCodeInternal, message)
		return
	}

	if result.RowsAffected == 0 {
		utils.RespondNotFound(c, "Todo not found")
		return
	}

	c.Status(http.StatusNoContent)
}
