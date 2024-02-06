package main

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/handlers"
	"net/http"
)

func main() {
	address := parseFlags()
	fmt.Printf("Server started on %v\n", address)
	err := http.ListenAndServe(address, handlers.GetRouter())
	if err != nil {
		panic(err)
	}
}
