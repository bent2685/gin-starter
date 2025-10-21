package database

import (
	"fmt"
	"gin-starter/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() {
	if config.AppConfig == nil {
		log.Fatal("配置未初始化")
	}

	// 构建PostgreSQL连接字符串
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Name,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.SSLMode,
		config.AppConfig.Database.Timezone, // 修复拼写错误
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 获取通用数据库对象 sql.DB 以设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("获取数据库对象失败: %v", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)           // 空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接可复用的最大时间

	log.Println("数据库连接成功")
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// Close 关闭数据库连接
func Close() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("获取数据库对象失败: %v", err)
			return
		}
		sqlDB.Close()
		log.Println("数据库连接已关闭")
	}
}
