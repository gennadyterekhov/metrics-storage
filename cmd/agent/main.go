package main

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/agent"
)

func main() {
	netAddress, reportInterval, pollInterval := parseFlags()

	err := agent.Agent(netAddress, reportInterval, pollInterval)
	if err != nil {
		fmt.Println(err.Error())
	}
}
