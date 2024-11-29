package main

import (
	"bufio"
	"fmt"
	"net/http"
	"time"
)

// create a buffering channel to make sure server closes after getting the end handler
var done = make(chan bool, 1)

/************************* server functions ********************************/

// a handler is an object implementing the http.Handler interface.
// functions serving as handlers take a http.responsewriter and a http.request as arguments.
// the response writer is used to fill in the http response.
func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n") // a simple hello as response, writes to w
}

// this handler reads all the http request headers and echo them into the response body
func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			// write all the headers from the request into w
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

// this handler performs a
func context(w http.ResponseWriter, req *http.Request) {
	// a context.Context is created for each request by the net/http machinery,
	// and is available with the Context() method
	ctx := req.Context() // get the context from the request
	fmt.Println("server: hello handler started")

	defer fmt.Println("server: hello handler ended") // output that the handler has ended

	select {
	// wait for a few seconds before sending a reply to the client. this could simulate some work the server is doing.
	case <-time.After(10 * time.Second):
		// if everything is fine and the connection is still stable after 10 seconds, the server will send hello to the client
		fmt.Fprintf(w, "hello\n")
	// while working, keep an eye on the context's Done() channel for a signal that we should cancel the work and return asap
	case <-ctx.Done():
		// the context's Err()method returns an error that explains why the Done() channel was closed
		err := ctx.Err()
		fmt.Println("server error: ", err) // stop the connection to simulate this error
		// write the error to w, note that the client might not catch this error,
		// since it could be the client that prematurely closed the connection
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
	}
}

// this handler will terminate the program
func end(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "goodbye\n") // a simple hello as response, writes to w
	done <- true                // signal done channel to end the program
}

/************************* client functions ********************************/

func ShowHttpClientExample() {
	// issue an http get request to our own server
	resp, err := http.Get("http://localhost:8090/headers")
	if err != nil {
		panic(err) // yse panic to fail on errors that shouldn't occur during normal operation
	}
	// use defer to make sure the response is closed at the end of the program
	defer resp.Body.Close()
	// print response status
	fmt.Println("Response status:", resp.Status)
	// use scanner to read the first 5 lines of the response body
	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}
	// if there is an error with scanner, then panic
	if err := scanner.Err(); err != nil {
		panic(err)
	}

}

// will block the thread
func main() {

	// start a goroutine for the server, so that it does not block the client
	go func() {
		http.HandleFunc("/hello", hello)
		http.HandleFunc("/headers", headers)
		http.HandleFunc("/context", context)
		http.HandleFunc("/end", end)
		http.ListenAndServe(":8090", nil)
	}()
	fmt.Println("server started on port 8090")
	// since the server is running on a goroutine, the client will run and send the get to the server
	ShowHttpClientExample()
	// this blocks the main process and wait, will terminate when receives the end handler
	<-done

}
