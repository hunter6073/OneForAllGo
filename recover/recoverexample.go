package main

import "fmt"

// this function panics
func mayPanic() {
	panic("a problem")
}

func main() {
	//recover must be called within a deferred function. when the enclosing function panics,
	// the defer will activate and a recover call within it will catch the panic
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	mayPanic()
	// this line will not run, maypanic will panic and defer will run, thus this line is never reached
	fmt.Println("After mayPanic()")
}
