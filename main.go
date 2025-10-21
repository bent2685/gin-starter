package main

import (
	"fmt"
	"gin-starter/config"
	"gin-starter/internal/interfaces/routes"
	"gin-starter/internal/middleware"
	"gin-starter/pkg/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.Init("")

	// 设置Gin模式
	if config.AppConfig.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 初始化日志系统
	if err := utils.SetupLogging(config.AppConfig.Log); err != nil {
		utils.Log.Fatalf("日志系统初始化失败: %v", err)
	}

	// 记录启动日志
	utils.Log.Info("服务启动中...")

	// 创建Gin引擎
	r := gin.New()
	// 注册全局中间件
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.RequestIDMiddleware())

	routerManager := routes.NewRouterManager()
	routerManager.RegisterRouter(routes.NewTestRouter())
	routerManager.SetupRoutes(r)

	addr := fmt.Sprintf("%s:%s", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	// 启动服务器
	if err := r.Run(addr); err != nil {
		log.Fatal("服务器启动失败:", err)
	}

	utils.Log.WithFields(map[string]interface{}{
		"addr": addr,
	}).Info("服务器配置信息")

	utils.Log.Success("服务启动完成")
}
