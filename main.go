package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"jwt-practice/handlers"
	"jwt-practice/middleware"
	"jwt-practice/models"
)

func main() {
	// 連接資料庫
	db, err := gorm.Open(sqlite.Open("jwt.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("無法連接資料庫:", err)
	}

	// 透過AutoMigrate自動建立、更新資料表
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("資料庫遷移失敗:", err)
	}

	// 創建 Gin 路由引擎
	r := gin.Default()

	// 初始化處理器
	authHandler := handlers.NewAuthHandler(db)

	// 公開路由
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.GET("/allProfile", authHandler.GetAllProfile)
	r.DELETE("/deleteUser", authHandler.DeleteProfile)

	// 需要驗證的路由
	auth := r.Group("/user")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/profile", authHandler.GetUserProfile)
	}

	// 啟動服務器
	if err := r.Run(":8081"); err != nil {
		log.Fatal("服務器啟動失敗:", err)
	}
}
