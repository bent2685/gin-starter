# Gin Starter

一个基于 Go 语言和 Gin 框架的 Web 应用启动模板。

## 功能特性

- 基于 Gin 框架的高性能 Web 服务
- 配置管理 (Viper)
- 结构化日志系统 (Logrus + Lumberjack)
- 数据库集成 (GORM + PostgreSQL)
- 模块化设计

## 数据库集成

本项目集成了 GORM ORM 框架和 PostgreSQL 数据库。

### 配置

数据库配置可以通过 `config.yaml` 文件进行配置:

```yaml
database:
  host: "localhost"              # 数据库主机地址
  port: "5432"                   # 数据库端口
  user: "postgres"               # 数据库用户名
  password: "postgres"           # 数据库密码
  name: "gin_starter"            # 数据库名称
  sslmode: "disable"             # SSL模式
  timezone: "Asia/Shanghai"      # 时区
```

或者通过环境变量配置:

```bash
STARTER_DATABASE_HOST=localhost
STARTER_DATABASE_PORT=5432
STARTER_DATABASE_USER=postgres
STARTER_DATABASE_PASSWORD=postgres
STARTER_DATABASE_NAME=gin_starter
STARTER_DATABASE_SSLMODE=disable
STARTER_DATABASE_TIMEZONE=Asia/Shanghai
```

### 数据库迁移

项目支持自动数据库迁移功能:

```bash
# 执行数据库迁移
go run main.go migrate
```

### 使用方法

在代码中使用数据库:

```go
import "gin-starter/internal/infra/database"

// 获取数据库实例
db := database.GetDB()

// 使用 GORM 进行数据库操作
var user User
db.First(&user, 1)
```

## 日志系统

本项目使用 Logrus 作为日志框架，并集成 Lumberjack 实现日志轮转功能。

### 特性

- 结构化日志记录
- 彩色控制台输出
- 日志级别控制 (trace, debug, info, warn, error, fatal, panic)
- 日志文件轮转 (基于文件大小和时间)
- 控制台和文件双输出
- 支持 JSON 和文本两种格式

### 配置

日志系统可以通过 `config.yaml` 文件进行配置:

```yaml
log:
  enabled: true                  # 是否启用日志文件输出
  level: "info"                  # 日志级别: debug, info, warn, error, fatal, panic, trace
  file: "./logs/server.log"      # 日志文件路径
  format: "text"                 # 日志格式: text, json
  max_size: 100                  # 日志文件最大大小(MB)
  max_backups: 3                 # 保留的旧日志文件最大数量
  max_age: 7                     # 保留旧日志文件的最大天数
  compress: true                 # 是否压缩旧日志文件
  enable_colors: true            # 是否启用控制台颜色输出
  timestamp_format: "2006-01-02 15:04:05"  # 时间戳格式
```

### 使用方法

在代码中使用全局日志实例:

```go
import "gin-starter/pkg/utils"

// 基本日志记录
utils.Log.Info("这是一条信息日志")
utils.Log.Warn("这是一条警告日志")
utils.Log.Error("这是一条错误日志")

// 不同级别的日志
utils.Log.Debug("调试信息")
utils.Log.Trace("跟踪信息")

// 格式化日志记录
utils.Log.Infof("用户 %s 登录成功", "张三")

// 带字段的日志记录
utils.Log.WithFields(logrus.Fields{
    "user_id": 12345,
    "ip": "192.168.1.100",
}).Info("用户登录")

// 带错误信息的日志记录
utils.Log.WithError(err).Error("操作失败")

// 自定义成功日志
utils.Log.Success("操作成功完成")

// HTTP 请求日志记录
utils.Log.LogHTTPRequest("GET", "/api/users", 200, 45.5)
```

### 日志颜色说明

- `PANIC`: 红色背景白色字体
- `FATAL`: 红色背景白色字体
- `ERROR`: 红色字体
- `WARN`: 黄色字体
- `INFO`: 蓝色字体
- `DEBUG`: 绿色字体
- `TRACE`: 紫色字体
- `SUCCESS`: 绿色背景白色字体
- `HTTP`: 青色字体