package main

import "github.com/gennadyterekhov/metrics-storage/internal/agent"

func main() {
	err := agent.Agent()
	if err != nil {
		panic(err)
	}
}
