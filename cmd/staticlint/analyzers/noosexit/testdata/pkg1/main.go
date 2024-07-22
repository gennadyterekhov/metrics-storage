package main

import (
	"os"
)

func main() {
	os.Exit(0) // want "cannot use os Exit in main function of package main"
}
