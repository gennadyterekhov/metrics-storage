// client part of metrics-storage.
// it sends metrics to the server in configured intervals in agent.Config.
package main

import (
	"fmt"

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

	instance := agent.New()
	err := instance.Start()
	if err != nil {
		logger.Custom.Infoln(err.Error())
	}
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
