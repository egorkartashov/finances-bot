package logger

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.Logger

func InitLogger(devMode bool) {

	var l *zap.Logger
	var err error
	if devMode {
		l, err = zap.NewDevelopment()
	} else {
		cfg := zap.NewProductionConfig()
		cfg.DisableCaller = true
		l, err = cfg.Build()
	}

	if err != nil {
		log.Fatal("Could create logger")
	}

	logger = l
	logger.Info("devmode", zap.Bool("devmode", devMode))
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
