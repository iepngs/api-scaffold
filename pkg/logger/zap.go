package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewLogger 创建 Zap 日志记录器
func NewLogger() *zap.Logger {
	// 配置日志文件输出
	fileWriter := &lumberjack.Logger{
		Filename:   "logs/2025-07-15.log",
		MaxSize:    100, // MB
		MaxBackups: 3,   // 保留最近3天的日志
		MaxAge:     3,   // 保留3天
		Compress:   true,
	}

	// 配置终端输出
	consoleWriter := zapcore.Lock(os.Stdout)

	// 创建多输出核心
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(fileWriter),
			zapcore.InfoLevel,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			consoleWriter,
			zapcore.InfoLevel,
		),
	)

	// 创建 Zap 日志记录器
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger
}
