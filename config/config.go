package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 保存所有应用配置
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Log    LogConfig    `mapstructure:"log"`
}

// ServerConfig 保存服务器相关配置
type ServerConfig struct {
	Port        string        `mapstructure:"port"`
	Host        string        `mapstructure:"host"`
	Environment string        `mapstructure:"environment"`
	Timeout     time.Duration `mapstructure:"timeout"`
	EnableCORS  bool          `mapstructure:"enable_cors"`
}

// LogConfig 保存日志相关配置
type LogConfig struct {
	Enabled         bool   `mapstructure:"enabled"`
	Level           string `mapstructure:"level"`
	File            string `mapstructure:"file"`
	Format          string `mapstructure:"format"`
	MaxSize         int    `mapstructure:"max_size"`
	MaxBackups      int    `mapstructure:"max_backups"`
	MaxAge          int    `mapstructure:"max_age"`
	Compress        bool   `mapstructure:"compress"`
	EnableColors    bool   `mapstructure:"enable_colors"`
	TimestampFormat string `mapstructure:"timestamp_format"`
}

// AppConfig 是全局配置实例
var AppConfig *Config

// Init 初始化应用配置
func Init(configPath string) {
	v := viper.New()

	// 1. 设置默认值
	v.SetDefault("server.port", "7070")
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.environment", "development")
	v.SetDefault("server.timeout", "30s")
	v.SetDefault("server.enable_cors", true)
	v.SetDefault("log.level", "info")
	v.SetDefault("log.file", "./logs/server.log")
	v.SetDefault("log.format", "text")
	v.SetDefault("log.max_size", 100)
	v.SetDefault("log.max_backups", 3)
	v.SetDefault("log.max_age", 7)
	v.SetDefault("log.compress", true)
	v.SetDefault("log.enable_colors", true)
	v.SetDefault("log.timestamp_format", "2006-01-02 15:04:05")

	// 2. 设置配置文件
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	if configPath != "" {
		v.AddConfigPath(configPath)
	} else {
		v.AddConfigPath(".")
		v.AddConfigPath("./conf")
	}

	// 3. 绑定环境变量
	v.SetEnvPrefix("STARTER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 4. 读取和解析配置
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("未找到配置文件，查找路径: %v", v.ConfigFileUsed())
			log.Printf("当前工作目录和查找的配置路径:")
			for _, path := range []string{".", "./conf"} {
				log.Printf("  - %s", path)
			}
			log.Println("将使用默认值和环境变量")
		} else {
			log.Fatalf("无法读取配置文件: %v", err)
		}
	} else {
		log.Printf("配置文件加载成功: %s", v.ConfigFileUsed())
	}

	if err := v.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("无法解析配置: %v", err)
	}

	log.Println("配置加载成功")
}

// IsProduction 检查是否为生产环境
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// IsDevelopment 检查是否为开发环境
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}
