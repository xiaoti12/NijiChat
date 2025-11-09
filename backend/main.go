//go:build js && wasm
// +build js,wasm

package main

import (
	"log"
	"net/http"

	"seiyuu-chat/database"
	"seiyuu-chat/handlers"
	"seiyuu-chat/middleware"
	"seiyuu-chat/router"
	"seiyuu-chat/services"

	"github.com/syumai/workers"
)

// Recovery 恢复中间件
func RecoveryMiddleware() router.HandlerFunc {
	return func(c *router.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"success": false,
					"error":   "内部服务器错误",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func main() {
	// 创建自定义路由器
	engine := router.New()

	// 添加全局中间件
	engine.Use(RecoveryMiddleware())
	engine.Use(middleware.LoggingMiddleware())
	engine.Use(middleware.CORSMiddleware())
	engine.Use(middleware.ErrorLoggingMiddleware())

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
	relationshipService := services.NewRelationshipService(db, cache, seiyuuService)

	// 初始化处理器层
	seiyuuHandler := handlers.NewSeiyuuHandler(seiyuuService)
	adminHandler := handlers.NewAdminHandler()
	relationshipsHandler := handlers.NewRelationshipsHandler(relationshipService)

	// 注册路由
	setupRoutes(engine, seiyuuHandler, adminHandler, relationshipsHandler)

	// 使用syumai/workers启动Worker
	workers.Serve(engine)
}

// setupRoutes 设置路由
func setupRoutes(
	engine *router.Engine,
	seiyuuHandler *handlers.SeiyuuHandler,
	adminHandler *handlers.AdminHandler,
	relationshipsHandler *handlers.RelationshipsHandler,
) {
	// API版本组
	api := engine.Group("/api")

	// 健康检查
	api.GET("/health", adminHandler.HealthCheck)

	// 公开接口 - 声优相关
	api.GET("/seiyuu", seiyuuHandler.GetAllSeiyuu)
	api.GET("/seiyuu/:id", seiyuuHandler.GetSeiyuuByID)

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
		admin.GET("/seiyuu/moegirl/:name", seiyuuHandler.GetMoegirlRawData)

		// 声优关系管理
		admin.POST("/relationships", relationshipsHandler.CreateRelationship)               // 创建关系
		admin.GET("/relationships", relationshipsHandler.GetRelationship)                   // 获取特定关系或所有关系
		admin.PUT("/relationships/:id", relationshipsHandler.UpdateRelationship)           // 更新关系
		admin.DELETE("/relationships/:id", relationshipsHandler.DeleteRelationship)        // 删除关系
		admin.GET("/seiyuu/:id/relationships", relationshipsHandler.GetSeiyuuRelationships) // 获取声优的所有关系
	}

	// 404处理
	engine.NoRoute(func(c *router.Context) {
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "接口不存在",
		})
	})

}
