// server part of metrics-storage.
// it deals with saving metrics to the db.
// exposes several endpoints to get and save custom metrics.
package main

import (
	"fmt"

	"github.com/gennadyterekhov/metrics-storage/internal/server/app"
)

func main() {
	fmt.Println("Starting")

	appInstance := app.New()
	err := appInstance.StartServer()
	if err != nil {
		panic(err)
	}
}
