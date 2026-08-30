package main

import "go.uber.org/zap"

func main() {
	logger, _ := zap.NewProduction()

	defer logger.Sync()

	url := "https://imooc.com"

	logger.Info("failed to fetch URL",
		zap.String("url", url),
		zap.Int("nums", 3),
	)
}
