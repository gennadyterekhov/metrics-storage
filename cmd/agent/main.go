package main

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/agent"
)

func shouldContinue(iter int) bool {
	return true
}

func main() {
	netAddress, reportInterval, pollInterval := parseFlags()

	err := agent.Agent(netAddress, shouldContinue, reportInterval, pollInterval)
	if err != nil {
		fmt.Println(err.Error())
	}
}
