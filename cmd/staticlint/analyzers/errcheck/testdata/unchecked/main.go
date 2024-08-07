package unchecked

import "fmt"

func fn() {
	_ = getErr()              // want `assignment with unchecked error`
	_, _ = getIntAndErr()     // want `assignment with unchecked error`
	_, _ = getInt(), getErr() // want `assignment with unchecked error`
	getErr()                  // want `expression returns unchecked error`
	getPtr()
}

func getErr() error {
	return fmt.Errorf("error")
}

func getInt() int {
	return 0
}

func getIntAndErr() (int, error) {
	return 0, fmt.Errorf("error")
}

func getPtr() *error {
	err := fmt.Errorf("error")
	return &err
}
