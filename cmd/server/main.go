package main

import (
	"github.com/gennadyterekhov/metrics-storage/internal/router"
	"net/http"
)

func main() {
	err := http.ListenAndServe(`:8080`, router.GetRouter())
	if err != nil {
		panic(err)
	}
}
