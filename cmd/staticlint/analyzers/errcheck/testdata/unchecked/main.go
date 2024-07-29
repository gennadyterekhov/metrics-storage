package unchecked

import "fmt"

func fn() {
	_ = getErr() // want `assignment with unchecked error`
}

func getErr() error {
	return fmt.Errorf("error")
}
