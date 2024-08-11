// client part of metrics-storage.
// it sends metrics to the server in configured intervals in agent.Config.
package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/gennadyterekhov/metrics-storage/internal/agent"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
)

// use go run -ldflags "-X main.buildVersion=v1.0.1 -X main.buildDate=01.01.2020 -X main.buildCommit=cafebabe" . to set these vars
var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printBuildInfo()

	config := getConfig()
	_, err := fmt.Printf("Agent started with server addr %v\n", config.Addr)
	if err != nil {
		panic(err)
	}
	if config.IsGzip {
		logger.Custom.Infoln("Attention, using gzip")
	}

	rootContext, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	go gracefulShutdown(rootContext)

	agent.RunAgent(rootContext, config)
}

func printBuildInfo() {
	printOrPanic(fmt.Sprintf("Build version: %v", buildVersion))
	printOrPanic(fmt.Sprintf("Build date: %v", buildDate))
	printOrPanic(fmt.Sprintf("Build commit: %v", buildCommit))
}

func printOrPanic(data string) {
	_, err := fmt.Println(data)
	if err != nil {
		panic(err)
	}
}

// gracefulShutdown - this code runs if app gets any of (syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
func gracefulShutdown(ctx context.Context) {
	<-ctx.Done()
	logger.Custom.Infoln("graceful shutdown. waiting a little")
	time.Sleep(time.Second)
}
