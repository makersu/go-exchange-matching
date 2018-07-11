package engine

import "fmt"

// TODO: Refactor for logging
func myPrintln(format string, a ...interface{}) {
	var log bool = false
	if log {
		fmt.Println(format, a)
	}
}
