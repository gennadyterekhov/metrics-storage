package main

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/handlers"
	"net/http"
)

func main() {
	fmt.Println("server func main")

	err := http.ListenAndServe(`:8080`, handlers.GetRouter())
	if err != nil {
		panic(err)
	}
}
