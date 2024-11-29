package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// go signal notification works by sending os.Signal values on a channel. we create a channel to receive
	// these notifications. note that this channel should be buffered
	sigs := make(chan os.Signal, 1)
	// signal.Notify registers the given channel to receive notifications of the specified signals.
	// Notify can take multiple signals starting from the second parameter. in this case,
	// we listen for sigint and sigterm
	// after running the program, use ctrl+c to send the signal
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1) // create buffering channel
	// this goroutine executes a blocking receive for signals. when it gets one it'll print it out and
	// then notify the programt that it can finish
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()
	fmt.Println("awaiting signal")
	// the program will wait here until it gets the expected signal
	<-done
	fmt.Println("exiting")
}
