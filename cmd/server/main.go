package main

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/router"
	"net/http"
)

func main() {
	fmt.Println("server func main")

	err := http.ListenAndServe(`:8080`, router.GetRouter())
	if err != nil {
		panic(err)
	}
}
