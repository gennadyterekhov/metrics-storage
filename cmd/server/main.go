package main

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/handlers"
	"net/http"
)

func main() {
	address := getAddress()
	fmt.Printf("Server started on %v\n", address)
	err := http.ListenAndServe(address, handlers.GetRouter())
	if err != nil {
		panic(err)
	}
}
