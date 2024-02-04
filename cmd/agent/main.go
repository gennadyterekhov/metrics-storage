package main

import "github.com/gennadyterekhov/metrics-storage/internal/agent"

func shouldContinue() bool {
	return true
}

func main() {
	url := `http://localhost:8080`
	err := agent.Agent(url, shouldContinue)
	if err != nil {
		panic(err)
	}
	//resty
}
