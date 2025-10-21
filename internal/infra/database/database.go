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
	var err error

	// 构建PostgreSQL连接字符串
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Name,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.SSLMode,
		config.AppConfig.Database.TimeZone,
	)

	// 配置GORM
	dbConfig := &gorm.Config{
		// 可以在这里添加GORM配置选项
	}

	// 连接数据库
	DB, err = gorm.Open(postgres.Open(dsn), dbConfig)
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}

	// 获取通用数据库对象 sql.DB 以进行连接池配置
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("无法获取数据库对象: %v", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)           // 空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

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
			log.Printf("无法获取数据库对象: %v", err)
			return
		}
		err = sqlDB.Close()
		if err != nil {
			log.Printf("关闭数据库连接失败: %v", err)
		} else {
			log.Println("数据库连接已关闭")
		}
	}
}
