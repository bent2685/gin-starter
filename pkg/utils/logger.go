package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"gin-starter/config"
)

// Logger 是 logrus 实例的包装
type Logger struct {
	*logrus.Logger
}

// Log 是全局日志实例
var Log *Logger

// SetupLogging 初始化日志系统
func SetupLogging(configLog config.LogConfig) error {
	// 创建日志实例
	logger := logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(configLog.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// 设置时间戳格式
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: configLog.TimestampFormat,
	})

	// 设置日志格式
	if configLog.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: configLog.TimestampFormat,
		})
	} else {
		// 创建带颜色的格式化器
		formatter := &logrus.TextFormatter{
			TimestampFormat: configLog.TimestampFormat,
			FullTimestamp:   true,
			ForceColors:     configLog.EnableColors,
			DisableColors:   !configLog.EnableColors,
		}
		logger.SetFormatter(formatter)
	}

	// 创建支持颜色的输出
	colorableOutput := colorable.NewColorableStdout()

	// 如果启用了日志文件输出
	if configLog.Enabled {
		// 确保日志目录存在
		logFile := configLog.File
		if logFile != "" {
			dir := filepath.Dir(logFile)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("创建日志目录失败: %v", err)
			}

			// 创建日志轮转写入器
			lumberjackLogger := &lumberjack.Logger{
				Filename:   logFile,
				MaxSize:    configLog.MaxSize,    // megabytes
				MaxBackups: configLog.MaxBackups, // 最多保留文件数
				MaxAge:     configLog.MaxAge,     // 最大保存天数
				Compress:   configLog.Compress,   // 是否压缩
			}

			// 同时输出到文件和控制台（控制台支持颜色）
			logger.SetOutput(io.MultiWriter(colorableOutput, lumberjackLogger))
		}
	} else {
		// 只输出到控制台（支持颜色）
		logger.SetOutput(colorableOutput)
	}

	// 创建全局日志实例
	Log = &Logger{logger}
	return nil
}

// WithField 添加字段到日志条目
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithFields 添加多个字段到日志条目
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// Debug 记录 debug 级别日志
func (l *Logger) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}

// Info 记录 info 级别日志
func (l *Logger) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

// Warn 记录 warning 级别日志
func (l *Logger) Warn(args ...interface{}) {
	l.Logger.Warn(args...)
}

// Error 记录 error 级别日志
func (l *Logger) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

// Fatal 记录 fatal 级别日志并退出程序
func (l *Logger) Fatal(args ...interface{}) {
	l.Logger.Fatal(args...)
}

// Panic 记录 panic 级别日志并引发 panic
func (l *Logger) Panic(args ...interface{}) {
	l.Logger.Panic(args...)
}

// Debugf 记录格式化的 debug 级别日志
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
}

// Infof 记录格式化的 info 级别日志
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

// Warnf 记录格式化的 warning 级别日志
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logger.Warnf(format, args...)
}

// Errorf 记录格式化的 error 级别日志
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

// Fatalf 记录格式化的 fatal 级别日志并退出程序
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args...)
}

// Panicf 记录格式化的 panic 级别日志并引发 panic
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.Logger.Panicf(format, args...)
}

// Trace 记录 trace 级别日志
func (l *Logger) Trace(args ...interface{}) {
	l.Logger.Trace(args...)
}

// Tracef 记录格式化的 trace 级别日志
func (l *Logger) Tracef(format string, args ...interface{}) {
	l.Logger.Tracef(format, args...)
}

// WithError 添加错误信息到日志条目
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// WithTime 添加时间信息到日志条目
func (l *Logger) WithTime(t time.Time) *logrus.Entry {
	return l.Logger.WithTime(t)
}

// 记录函数调用信息的辅助函数
func (l *Logger) WithFuncInfo(file string, line int, funcName string) *logrus.Entry {
	return l.Logger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
		"func": funcName,
	})
}

// 自定义颜色日志方法
func (l *Logger) Success(args ...interface{}) {
	l.Logger.WithFields(logrus.Fields{
		"level": "SUCCESS",
	}).Info(args...)
}

func (l *Logger) Successf(format string, args ...interface{}) {
	l.Logger.WithFields(logrus.Fields{
		"level": "SUCCESS",
	}).Infof(format, args...)
}

// 结构化日志记录方法
func (l *Logger) LogWithDetails(level string, message string, details map[string]interface{}) {
	fields := logrus.Fields{"level": level}
	for k, v := range details {
		fields[k] = v
	}
	l.Logger.WithFields(fields).Info(message)
}

// HTTP 请求日志记录方法
func (l *Logger) LogHTTPRequest(method, url string, statusCode int, latency float64) {
	l.Logger.WithFields(logrus.Fields{
		"http_method": method,
		"http_url":    url,
		"status_code": statusCode,
		"latency_ms":  latency,
		"level":       "HTTP",
	}).Info("HTTP Request")
}
