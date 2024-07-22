package checked

import "fmt"

func fn() {
	err := getErr()
	if err != nil {
	}
}

func getErr() error {
	return fmt.Errorf("error")
}
