package database

import (
	"gin-starter/internal/domain/models"
	"gin-starter/internal/domain/models/rbac"
	"log"

	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() {
	if DB == nil {
		log.Fatal("数据库未初始化")
	}

	// 添加需要迁移的模型
	// 注意：Casbin 使用自己的表来管理用户-角色关系和角色-权限关系
	err := DB.AutoMigrate(
		&models.User{},
		&rbac.Role{}, // 只保留角色表用于业务逻辑，权限关系由Casbin管理
	)

	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 使用Casbin官方适配器创建表
	// gorm-adapter 会在首次使用时自动创建所需的表
	_, err = gormadapter.NewAdapterByDB(DB)
	if err != nil {
		log.Fatalf("Casbin适配器初始化失败: %v", err)
	}

	log.Println("数据库迁移完成")
}
