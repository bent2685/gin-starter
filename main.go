package main

import (
	"fmt"
	"gin-starter/config"
	"gin-starter/internal/application/services/rbac"
	"gin-starter/internal/infra/database"
	"gin-starter/internal/interfaces/routes"
	"gin-starter/internal/interfaces/validators"
	"gin-starter/internal/middleware"
	"gin-starter/pkg/utils"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 检查是否有迁移参数
	args := os.Args[1:]
	runMigration := false
	for _, arg := range args {
		if arg == "migrate" {
			runMigration = true
			break
		}
	}

	// 初始化配置
	config.Init("")

	// 设置Gin模式
	if config.AppConfig.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 注册自定义验证器
	validators.RegisterCustomValidators()

	// 初始化日志系统
	if err := utils.SetupLogging(config.AppConfig.Log); err != nil {
		utils.Log.Fatalf("日志系统初始化失败: %v", err)
	}

	// 记录启动日志
	utils.Log.Info("服务启动中...")

	// 初始化数据库
	database.InitDatabase()
	defer database.Close()

	// 如果是迁移模式，执行迁移后退出
	if runMigration {
		utils.Log.Info("执行数据库迁移...")
		database.AutoMigrate()

		// 初始化RBAC服务（仅用于添加示例数据）
		rbacService, err := rbac.NewRBACService()
		if err != nil {
			utils.Log.Fatalf("RBAC服务初始化失败: %v", err)
		}

		// 添加一些默认策略示例
		rbacService.AddPolicy("admin", "/users/*", "*")
		rbacService.AddPolicy("user", "/users/:id", "GET")
		rbacService.AddRoleForUser("1", "admin")
		rbacService.SavePolicy()

		utils.Log.Info("数据库迁移完成")
		return
	}

	// 初始化RBAC服务
	rbacService, err := rbac.NewRBACService()
	if err != nil {
		utils.Log.Fatalf("RBAC服务初始化失败: %v", err)
	}

	// 创建Gin引擎
	r := gin.New()
	// 注册全局中间件
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.RequestIDMiddleware())

	routerManager := routes.NewRouterManager()
	routerManager.RegisterRouter(routes.NewTestRouter())
	routerManager.RegisterRouter(routes.NewUserRouter())                 // 注册用户路由
	routerManager.RegisterRouter(routes.NewRBACRouter(rbacService))      // 注册RBAC路由
	routerManager.RegisterRouter(routes.NewProtectedRouter(rbacService)) // 注册受保护路由
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
