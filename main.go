package main

import (
	"fmt"
	"gin-starter/config"
	"gin-starter/internal/application/services/rbac"
	"gin-starter/internal/infra/database"
	"gin-starter/internal/infra/ofs"
	"gin-starter/internal/interfaces/routes"
	"gin-starter/internal/interfaces/validators"
	"gin-starter/internal/middleware"
	"gin-starter/pkg/utils"
	"log"
	"os"
	"slices"

	"github.com/gin-gonic/gin"
)

func main() {
	// 检查是否有迁移参数
	args := os.Args[1:]
	runMigration := slices.Contains(args, "migrate")

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
	defer database.Close()
	database.InitDatabase()
	// 初始化s3
	ofs.InitOfs()
	if runMigration {
		utils.Log.Info("执行数据库迁移...")
		database.AutoMigrate()

		if err := rbac.InitRBAC(); err != nil {
			utils.Log.Fatalf("RBAC服务初始化失败: %v", err)
		}

		rbac.AddPolicy("admin", "/users/*", "*")
		rbac.AddPolicy("user", "/users/:id", "GET")
		rbac.AddPolicy("super_admin", "*", "*")
		rbac.AddRoleForUser("1", "admin")
		rbac.SavePolicy()

		utils.Log.Info("数据库迁移完成")
		return
	}

	if err := rbac.InitRBAC(); err != nil {
		utils.Log.Fatalf("RBAC服务初始化失败: %v", err)
	}

	r := gin.New()
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.RequestIDMiddleware())

	routerManager := routes.NewRouterManager()
	routerManager.RegisterRouter(routes.NewTestRouter())
	routerManager.RegisterRouter(routes.NewUserRouter())
	routerManager.RegisterRouter(routes.NewRBACRouter())
	routerManager.RegisterRouter(routes.NewDepartmentRouter())
	routerManager.RegisterRouter(routes.NewProtectedRouter())
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
