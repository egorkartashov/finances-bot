package logger

import (
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic("Could create logger")
	}

	logger = l.Sugar()
	defer func(logger *zap.SugaredLogger) {
		_ = logger.Sync()
	}(logger)
}

func Info(msg string) {
	logger.Info(msg)
}

func Infof(msg string, args ...interface{}) {
	logger.Infof(msg, args)
}

func Error(args ...interface{}) {
	logger.Error(args)
}

func Errorf(msg string, args ...interface{}) {
	logger.Errorf(msg, args)
}

func Fatal(msg string, args ...interface{}) {
	logger.Fatal(msg, args)
}
