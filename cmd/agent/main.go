package main

import (
	"context"
	"fmt"

	"github.com/gennadyterekhov/metrics-storage/internal/agent"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
)

func main() {
	config := getConfig()
	fmt.Printf("Agent started with server addr %v\n", config.Addr)
	if config.IsGzip {
		logger.ZapSugarLogger.Infoln("Attention, using gzip")
	}

	agent.RunAgent(context.Background(), config)
}
