package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"go-gin-todo-api/database"
	"go-gin-todo-api/handlers"
	"go-gin-todo-api/middleware"
)

func main() {
	// .env.localを優先的に読み込み（ローカル開発用）
	if err := godotenv.Load(".env.local"); err != nil {
		// .env.localがない場合は.envを読み込み
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env file not found, using system environment variables")
		}
	}

	database.InitDB()
	
	r := gin.Default()

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Status": http.StatusOK,
		})
	})

	// 認証エンドポイント（認証不要）
	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/refresh", handlers.Refresh)
		auth.POST("/logout", handlers.Logout)
	}

	// 認証必須エンドポイント
	api := r.Group("/")
	api.Use(middleware.AuthMiddleware())
	{
		// ユーザー確認
		api.GET("/me", handlers.GetMe)

		// Todoエンドポイント
		api.GET("/todos", handlers.GetTodos)
		api.POST("/todos", handlers.CreateTodo)
		api.GET("/todos/:id", handlers.GetTodo)
		api.PATCH("/todos/:id", handlers.UpdateTodo)
		api.DELETE("/todos/:id", handlers.DeleteTodo)
	}

	r.Run()
}
