package logger

import (
	"go.uber.org/zap"
)

var (
	logger *zap.SugaredLogger
)

func Init() {
	if logger != nil {
		return
	}

	zapLogger, err := zap.NewProduction(
		zap.IncreaseLevel(zap.InfoLevel),
		zap.AddCallerSkip(1))

	if err != nil {
		panic(err)
	}

	logger = zapLogger.Sugar()
}

func Infof(template string, args ...interface{}) {
	if logger != nil {
		logger.Infof(template, args...)
	}
}

func Errorf(template string, args ...interface{}) {
	if logger != nil {
		logger.Errorf(template, args...)
	}
}
