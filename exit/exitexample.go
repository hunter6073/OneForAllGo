package main

import (
	"fmt"
	"os"
)

func main() {
	// normally defers run at the end of the program, but if you us os.Exit, then
	// the defers will never be called.
	defer fmt.Println("!")

	// use os.Exit to exit the program with a non-zero status
	// the exit will be picked up by go and printed
	os.Exit(3)
}
