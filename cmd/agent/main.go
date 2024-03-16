package main

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/agent"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
)

func main() {
	config := getConfig()
	config.IsBatch = true
	fmt.Printf("Agent started with server addr %v\n", config.Addr)
	if config.IsGzip {
		logger.ZapSugarLogger.Infoln("Attention, using gzip")
	}

	err := agent.RunAgent(config)
	if err != nil {
		fmt.Println(err.Error())
	}
}
