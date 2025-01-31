package logger

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() (*zap.Logger, *zap.SugaredLogger, error) {
	logFilePath := fmt.Sprintf("logs/%s.log", time.Now().Format("2006-01-02"))

	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		return nil, nil, fmt.Errorf("не удалось создать директорию для логов: %w", err)
	}

	logRotation := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}

	logLevel := zapcore.InfoLevel
	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		if parsedLevel, err := zapcore.ParseLevel(lvl); err == nil {
			logLevel = parsedLevel
		}
	}

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
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(logRotation), logLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), logLevel),
	)

	logger := zap.New(core, zap.AddCaller())
	sugar := logger.Sugar()

	return logger, sugar, nil
}
