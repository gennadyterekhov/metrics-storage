package main

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/agent"
)

func main() {
	config := getConfig()

	err := agent.RunAgent(config.Addr, config.ReportInterval, config.PollInterval)
	if err != nil {
		fmt.Println(err.Error())
	}
}
