package main

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/agent"
)

func main() {
	config := getConfig()
	fmt.Printf("Agent started with server addr %v\n", config.Addr)

	err := agent.RunAgent(config)
	if err != nil {
		fmt.Println(err.Error())
	}
}
