package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// AppConfig 全局配置实例
var AppConfig *Config

// Config 应用配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	S3       S3Config       `mapstructure:"s3"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	Environment string `mapstructure:"environment"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
	Timezone string `mapstructure:"timezone"`
}

// LogConfig 日志配置
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

// S3Config S3配置
type S3Config struct {
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Region          string `mapstructure:"region"`
	Endpoint        string `mapstructure:"endpoint"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

// Init 初始化配置
func Init(configPath string) {
	// 设置配置文件名和路径
	viper.SetConfigName("config")
	if configPath != "" {
		viper.AddConfigPath(configPath)
	} else {
		viper.AddConfigPath(".")
	}

	// 设置环境变量前缀
	viper.SetEnvPrefix("STARTER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("配置文件未找到，使用默认配置")
		} else {
			log.Fatalf("配置文件读取失败: %v", err)
		}
	} else {
		log.Printf("配置文件加载成功: %s", viper.ConfigFileUsed())
	}

	// 设置默认值
	setDefaults()

	// 绑定环境变量
	bindEnvs()

	// 解析配置
	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("配置解析失败: %v", err)
	}

	log.Println("配置加载成功")
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器配置默认值
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", "7070")
	viper.SetDefault("server.environment", "development")

	// 数据库配置默认值
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.name", "gin_starter")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.timezone", "Asia/Shanghai")

	// 日志配置默认值
	viper.SetDefault("log.enabled", true)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.file", "./logs/server.log")
	viper.SetDefault("log.format", "text")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 3)
	viper.SetDefault("log.max_age", 7)
	viper.SetDefault("log.compress", true)
	viper.SetDefault("log.enable_colors", true)
	viper.SetDefault("log.timestamp_format", "2006-01-02 15:04:05")

	// JWT配置默认值
	viper.SetDefault("jwt.secret", "gin-starter-secret-key")
}

// bindEnvs 绑定环境变量
func bindEnvs() {
	// 服务器配置环境变量绑定
	viper.BindEnv("server.host", "STARTER_SERVER_HOST")
	viper.BindEnv("server.port", "STARTER_SERVER_PORT")
	viper.BindEnv("server.environment", "STARTER_SERVER_ENVIRONMENT")

	// 数据库配置环境变量绑定
	viper.BindEnv("database.host", "STARTER_DATABASE_HOST")
	viper.BindEnv("database.port", "STARTER_DATABASE_PORT")
	viper.BindEnv("database.user", "STARTER_DATABASE_USER")
	viper.BindEnv("database.password", "STARTER_DATABASE_PASSWORD")
	viper.BindEnv("database.name", "STARTER_DATABASE_NAME")
	viper.BindEnv("database.sslmode", "STARTER_DATABASE_SSLMODE")
	viper.BindEnv("database.timezone", "STARTER_DATABASE_TIMEZONE")

	// 日志配置环境变量绑定
	viper.BindEnv("log.enabled", "STARTER_LOG_ENABLED")
	viper.BindEnv("log.level", "STARTER_LOG_LEVEL")
	viper.BindEnv("log.file", "STARTER_LOG_FILE")
	viper.BindEnv("log.format", "STARTER_LOG_FORMAT")
	viper.BindEnv("log.max_size", "STARTER_LOG_MAX_SIZE")
	viper.BindEnv("log.max_backups", "STARTER_LOG_MAX_BACKUPS")
	viper.BindEnv("log.max_age", "STARTER_LOG_MAX_AGE")
	viper.BindEnv("log.compress", "STARTER_LOG_COMPRESS")
	viper.BindEnv("log.enable_colors", "STARTER_LOG_ENABLE_COLORS")
	viper.BindEnv("log.timestamp_format", "STARTER_LOG_TIMESTAMP_FORMAT")

	// JWT配置环境变量绑定
	viper.BindEnv("jwt.secret", "STARTER_JWT_SECRET")

	// S3配置环境变量绑定
	viper.BindEnv("s3.access_key_id", "STARTER_S3_ACCESS_KEY_ID")
	viper.BindEnv("s3.secret_access_key", "STARTER_S3_SECRET_ACCESS_KEY")
	viper.BindEnv("s3.region", "STARTER_S3_REGION")
	viper.BindEnv("s3.endpoint", "STARTER_S3_ENDPOINT")
}

// GetJWTSecret 获取JWT密钥
func GetJWTSecret() string {
	if AppConfig != nil {
		return AppConfig.JWT.Secret
	}
	return "gin-starter-secret-key"
}
