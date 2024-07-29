package main

import (
	"context"
	"fmt"

	"github.com/gennadyterekhov/metrics-storage/internal/agent"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
)

func main() {
	config := getConfig()
	_, err := fmt.Printf("Agent started with server addr %v\n", config.Addr)
	if err != nil {
		panic(err)
	}
	if config.IsGzip {
		logger.ZapSugarLogger.Infoln("Attention, using gzip")
	}

	agent.RunAgent(context.Background(), config)
}
