//go:build js && wasm
// +build js,wasm

package main

import (
	"log"
	"net/http"

	"seiyuu-chat/database"
	"seiyuu-chat/handlers"
	"seiyuu-chat/middleware"
	"seiyuu-chat/services"

	"github.com/gin-gonic/gin"
	"github.com/syumai/workers"
)

func main() {
	// 创建Gin路由器
	router := gin.New()

	// 添加全局中间件
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorLoggingMiddleware())

	// 初始化数据库和缓存
	db, err := database.NewD1Client("SEIYUU_DB")
	if err != nil {
		log.Fatalf("Failed to connect to D1 database: %v", err)
	}
	defer db.Close()

	cache, err := database.NewKVClient("SEIYUU_KV")
	if err != nil {
		log.Fatalf("Failed to connect to KV namespace: %v", err)
	}

	// 初始化服务层
	seiyuuService := services.NewSeiyuuService(db, cache)
	aiService, err := services.NewAIService()
	if err != nil {
		log.Fatalf("Failed to initialize AI service: %v", err)
	}
	schedulerService := services.NewSchedulerService(seiyuuService, aiService)
	moegirlService := services.NewMoegirlService()

	// 初始化处理器层
	seiyuuHandler := handlers.NewSeiyuuHandler(seiyuuService)
	schedulerHandler := handlers.NewSchedulerHandler(schedulerService)
	adminHandler := handlers.NewAdminHandler(db)
	moegirlHandler := handlers.NewMoegirlHandler(moegirlService, aiService)

	// 注册路由
	setupRoutes(router, seiyuuHandler, schedulerHandler, adminHandler, moegirlHandler)

	// 使用syumai/workers启动Worker
	workers.Serve(router)
}

// setupRoutes 设置路由
func setupRoutes(
	router *gin.Engine,
	seiyuuHandler *handlers.SeiyuuHandler,
	schedulerHandler *handlers.SchedulerHandler,
	adminHandler *handlers.AdminHandler,
	moegirlHandler *handlers.MoegirlHandler,
) {
	// API版本组
	api := router.Group("/api")

	// 健康检查
	api.GET("/health", adminHandler.HealthCheck)

	// 公开接口 - 声优相关
	api.GET("/seiyuu", seiyuuHandler.GetAllSeiyuu)
	api.GET("/seiyuu/:id", seiyuuHandler.GetSeiyuuByID)

	// 公开接口 - 智能调度器
	api.POST("/scheduler/select", middleware.APIRateLimitMiddleware(), schedulerHandler.SelectSeiyuu)

	// 管理员登录（无需认证）
	api.POST("/admin/login", adminHandler.Login)

	// 管理员接口（需要认证）
	admin := api.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleware())
	// admin.Use(middleware.AdminRateLimitMiddleware())
	{
		// 管理员信息
		admin.GET("/profile", adminHandler.GetProfile)

		// 声优管理
		admin.GET("/seiyuu", seiyuuHandler.GetAllSeiyuuAdmin) // 获取所有声优（含待审核）
		admin.POST("/seiyuu", seiyuuHandler.CreateSeiyuu)
		admin.PUT("/seiyuu/:id", seiyuuHandler.UpdateSeiyuu)
		admin.DELETE("/seiyuu/:id", seiyuuHandler.DeleteSeiyuu)

		// 萌娘百科集成
		admin.GET("/moegirl/:name", moegirlHandler.GetRawData)
		admin.GET("/moegirl/search", moegirlHandler.SearchSeiyuu)
		admin.POST("/process-profile", moegirlHandler.ProcessProfile)
	}

	// 404处理
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "接口不存在",
		})
	})

}
