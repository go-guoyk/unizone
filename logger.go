package main

import "go.uber.org/zap"

func NewLogger(verbose bool) *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()
	if verbose {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.Development = true
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		cfg.Development = false
	}
	cfg.Encoding = "console"
	logger, _ := cfg.Build()
	return logger.Sugar()
}
