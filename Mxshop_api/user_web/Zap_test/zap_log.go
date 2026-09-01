package main

import (
	"go.uber.org/zap"
	"time"
)

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./myproject.log",
	}
	return cfg.Build()
}

func runStructuredLogExample() {
	logger, err := NewLogger()
	if err != nil {
		panic(err)
	}

	su := logger.Sugar()
	defer su.Sync()
	url := "https://imoc.com"
	su.Info("failed to fetch URL",
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}
