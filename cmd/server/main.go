// server part of metrics-storage.
// it deals with saving metrics to the db.
// exposes several endpoints to get and save custom metrics.
package main

import (
	"context"
	"fmt"

	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/app"
)

// use go run -ldflags "-X main.buildVersion=v1.0.1 -X main.buildDate=01.01.2020 -X main.buildCommit=cafebabe" . to set these vars
var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printBuildInfo()

	appInstance := app.New()
	err := appInstance.StartServer(context.Background())
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
