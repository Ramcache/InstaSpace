package logger

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() (*zap.Logger, error) {
	logFilePath := fmt.Sprintf("logs/%s.log", time.Now().Format("2006-01-02"))

	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		return nil, fmt.Errorf("не удалось создать директорию для логов: %w", err)
	}

	logRotation := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}

	writeSyncer := zapcore.AddSync(logRotation)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), writeSyncer, zapcore.InfoLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
	)

	return zap.New(core, zap.AddCaller()), nil
}
