package database

import (
	"gin-starter/internal/domain/models"
	"log"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() {
	if DB == nil {
		log.Fatal("数据库未初始化")
	}

	// 添加需要迁移的模型
	err := DB.AutoMigrate(
		&models.User{},
	)

	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库迁移完成")
}
