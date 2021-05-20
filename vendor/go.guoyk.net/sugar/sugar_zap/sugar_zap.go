package sugar_zap

import (
	"go.guoyk.net/sugar"
	"go.uber.org/zap"
)

type wrapper struct {
	logger *zap.SugaredLogger
}

func (w *wrapper) Debug(message string, items ...interface{}) {
	w.logger.Debugw(message, items...)
}

func (w *wrapper) Info(message string, items ...interface{}) {
	w.logger.Infow(message, items...)
}

func (w *wrapper) Warn(message string, items ...interface{}) {
	w.logger.Warnw(message, items...)
}

func (w *wrapper) Error(message string, items ...interface{}) {
	w.logger.Errorw(message, items...)
}

func (w *wrapper) Panic(message string, items ...interface{}) {
	w.logger.Panicw(message, items...)
}

func (w *wrapper) Fatal(message string, items ...interface{}) {
	w.logger.Fatalw(message, items...)
}

func Wrap(logger *zap.Logger) sugar.Logger {
	return &wrapper{
		logger: logger.WithOptions(zap.AddCallerSkip(1)).Sugar(),
	}
}
