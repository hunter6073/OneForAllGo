package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// in this example our state will b owned by a single goroutine, in order to read or write that state,
// other goroutines will ssend messages to the owning goroutine and receive corresponding replies.
type readOp struct {
	key  int
	resp chan int
}
type writeOp struct {
	key  int
	val  int
	resp chan bool
}

// function that prints string 3 times
func fun3(from string) {
	for i := 0; i < 3; i++ {
		fmt.Println(from, ":", i)
	}
}

// this is the function we'll run in a goroutine. the done channel
// will be used to notify another goroutine that this function's work is done
func worker(done chan bool) {
	fmt.Print("working...")
	time.Sleep(time.Second)
	fmt.Println("done")
	done <- true // send true to channel
}

// ping is a function that only receives
func ping(pings chan<- string, msg string) {
	pings <- msg
}

// pong receives from pings and sends to pongs
func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}

// workers receive work on the jobs channel and send the result on results channel
func worker2(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		time.Sleep(time.Second) // sleep a second per job to simulate an expensive task
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

// this is the function to simulate an expensive task
func worker3(id int) {
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

// container holds a map fo counters,since we want to update it concurrently from multiple goroutines, we add a mutext to synchronize access.
// note that mutexes must not be copied. so if this struct is passed around, it should be done by pointer
type Container struct {
	mu       sync.Mutex
	counters map[string]int
}

// increment a key's value
func (c *Container) inc(name string) {
	c.mu.Lock()         // lock the mutex before accessing counters
	defer c.mu.Unlock() // unlock it at the end of the function
	c.counters[name]++  // update key
}

func main() {
	// goroutines, lightweight threads of execution
	fun3("direct")        // running a function synchronously
	go fun3("goroutine")  // funning a function using goroutine, will execute concurrently with the calling one
	go func(msg string) { // running an anonymous function
		fmt.Println(msg)
	}("going") // this is the input parameter

	// our two calls are running in seperate goroutines now, wait for them to finish,otherwise the function will end before the output
	// sleep function
	time.Sleep(time.Second)
	fmt.Println("fun3 and anonymous goroutine done")

	// by default channels are unbuffered, meaning that they will only
	// accept sends(chan<-) if there is a corresponding receive(<-chan) ready to receive the sent value
	messages := make(chan string) // create a new unbuffering channel, channels are typed by the values they convey
	// send a value into a channel using the channel <- syntax
	go func() { messages <- "ping" }()
	// get the message from the channel
	msg := <-messages // by default sends and receives block until both the sender and receiver are ready
	fmt.Println(msg)

	// buffered channels accept a limited number of values without a cooresponding receiver for those values
	messages2 := make(chan string, 2) // this means the channel allows buffering up to 2 values
	messages2 <- "buffered"
	messages2 <- "channel"
	//messages2 <- "third" // do not attempt to enter another value, buffer is full
	// first in first out
	fmt.Println("getting the first message from messages2: ", <-messages2)  // get s buffered
	fmt.Println("getting the second message from messages2: ", <-messages2) // gets channel

	messages2 <- "buffered again" // buffered channels can receive extra input after buffer is cleared
	fmt.Println(<-messages2)

	// channel synchronization
	done := make(chan bool, 1) // create a buffered channel
	go worker(done)            // start the goroutine, giving it the channel to notify on
	// if we remove <-done, the program would exit before the worker even started
	<-done // blocks until we receive a notification from the worker on the channel

	// create two channel functions pings and pongs
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)

	// using select across two channels
	c1 := make(chan string)
	c2 := make(chan string)

	// each channel will receive a value after some amount of time to
	// simulate blocking rpc operations executing in concurrent goroutines
	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(1 * time.Second)
		c2 <- "two"
	}()

	// use select to await both of these values simultaneously, printing each one as it arrives
	// select will wait until both messages arrive
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}

	// execute a call that returns on channel c1 after 2s, c1 is buffered
	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "result 1"
	}()

	// perform select to catch c1 or a timeout
	select {
	case res := <-c1:
		fmt.Println(res)
	// time.After is a channel as well
	case <-time.After(1 * time.Second): // timeout happens
		fmt.Println("timeout 1")
	}

	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "result 2"
	}()
	// for c2, the time to sleep is 2 seconds, the receive will succceed
	select {
	case res := <-c2:
		fmt.Println(res)
	case <-time.After(3 * time.Second):
		fmt.Println("timeout 2")
	}

	// using default clause in select to implement non-blocking sends,receives and multi-way selects
	signals := make(chan bool) // create an unbuffered channel signals

	// go func() { messages <- "what is love?" }() note that channel input must be done using goroutine
	// here is a non-blocking receive. if a value is available on messages then select will take the
	// <-messages case with that value. if not it will immediately take the default case
	select {
	case msg2 := <-messages:
		fmt.Println("received message", msg2)
	default:
		fmt.Println("no message received")
	}

	// here is a non-blocking send. here msg2 cannot be sent to the messages channel, because the channel
	// has no buffer and there is no receiver. Therefore the default case is selected
	msg2 := "hi"
	select {
	case messages <- msg2:
		fmt.Println("sent message", msg2)
	default:
		fmt.Println("no message sent")
	}

	// we can use multiple cases above the default clause to implement a multi-way non-blocking select.
	// here thwe attempt non-blocking receives on both messages and signals
	select {
	case msg2 := <-messages:
		fmt.Println("received message", msg2)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity")
	}

	// closing channels
	jobs := make(chan int, 5)
	done2 := make(chan bool)

	// worker goroutine, repetatedly receives from jobs with j,more:=<-jobs.
	go func() {
		for {
			// the more value will be false if jobs has been closed and all values in the channel have already been received.
			j, more := <-jobs
			if more {
				fmt.Println("received job", j)
			} else {
				fmt.Println("received all jobs")
				done2 <- true
				return
			}
		}
	}()

	// this sends 3 jobs to the worker over the jobs channel, then closes it
	for j := 1; j <= 3; j++ {
		jobs <- j
		fmt.Println("sent job", j)
	}
	// close a channel
	close(jobs)
	fmt.Println("sent all jobs")
	<-done2 // we await the wroker using the synchronization approach
	//reading from a closed channel succeeds immediately, returning the zero value of the underlying type.
	// the optional second value is true if the value received was delivered by a successful send operation or
	// false if it was a zero value generated because the channel is closed and empty.
	_, ok := <-jobs
	fmt.Println("received more jobs", ok)

	// perfroming range on channels
	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"
	fmt.Println("retrieve an element from queue:", <-queue)
	queue <- "three"
	close(queue)
	// the range iterates over each element as it's received from queue. because we closed the channel above,
	// the iteration terminates after receiving the 2 elements
	for elem := range queue { // not that it is possible to close a non-empty channel but still have the remaining values be received
		fmt.Println("final elements in queue:", elem)
	}

	// timer represents a single event in the future. you tell the timer how long you want to wait,
	// and it provides a channel that will be notified at that time. this timer will wait 2 seconds
	timer1 := time.NewTimer(2 * time.Second)
	<-timer1.C // blocks on the timer's channel C(must be C, C stands for channel) until it sends a value indicating that the timer fired
	fmt.Println("Timer 1 fired")
	// if you just wanted to wait, you could have used time.Sleep,
	// one reason a timer may be useful is that you can cancel the timer before it fires.
	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 fired")
	}()
	stop2 := timer2.Stop()
	if stop2 {
		fmt.Println("Timer 2 stopped")
	}
	// Give the timer2 enough time to fire,
	// if it ever was going to, to show it is in fact stopped.

	// the first timer will fire 2s after we start the program, but the second should be stopped before it has a chance to fire
	time.Sleep(2 * time.Second)

	// tickers use a similar mechanism to timers:a channel that is sent values.
	ticker := time.NewTicker(500 * time.Millisecond)
	// tickers can be stopped like timers, once a tiker is stopped it wont received any more values on its channel
	go func() {
		for {
			select {
			case <-done2:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()
	// stop the ticker after 1600ms
	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	// synchronously wait till done
	done2 <- true
	fmt.Println("Ticker stopped")

	//create buffering channels
	const numJobs = 5
	jobs3 := make(chan int, numJobs)
	results3 := make(chan int, numJobs)
	// start up 3 workers, initially blocked because there aren't any jobs yet
	for w := 1; w <= 3; w++ {
		go worker2(w, jobs3, results3)
	}

	// send 5 jobs, whoever gets the job first start doing it
	for j := 1; j <= numJobs; j++ {
		jobs3 <- j
	}
	close(jobs3) // close the channel to indicate that's all the work we have

	// finally collect all the results of the work
	for a := 1; a <= numJobs; a++ {
		<-results3
	}

	var wg sync.WaitGroup // this waitgroup is used to wait for all the goroutines launched here to finsh
	// if a waitgroup is explicitly passed into functions, it should be done by pointer

	for i := 1; i <= 5; i++ {
		wg.Add(1) // increment the waitgroup counter
		//Wrap the worker call in a closure that makes sure to tell the WaitGroup that this worker is done.
		// This way the worker itself does not have to be aware of the concurrency primitives involved in its execution.
		go func() {
			defer wg.Done()
			worker3(i)
		}()
	}
	//Block until the WaitGroup counter goes back to 0; all the workers notified they’re done.
	wg.Wait()

	// rate limiting
	requests := make(chan int, 5) // create a buffering channel that takes in at most 5 ints
	// buffer 1 to 5 into requests
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests) // close the channel, 1-5 is stored in che channel buffer

	// the limiter channel will receive a value every 200ms, this is the regulator in our rate limiting scheme
	limiter := time.Tick(200 * time.Millisecond)

	// by blocking on a receive from the limiter channel before serving each request,we limit ourselves to 1 request every 200 ms
	for req := range requests {
		<-limiter // limiter is triggered every 200ms, and so is this function
		fmt.Println("request", req, time.Now())
	}

	// creating a new buffering channel, this channel will allow bursts of up to 3 events
	burstyLimiter := make(chan time.Time, 3) // the channel type is time.Time

	// fill up the channel to represent allowed bursting
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	// every 200ms we'll try to add a new value to burstylimiter , up to it's limit of 3
	go func() {
		for t := range time.Tick(200 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()

	// now simulate 5 more incoming requests. the first 3 of these will benefit from the burst capability of burstyLimiter
	burstyRequests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		burstyRequests <- i
	}
	close(burstyRequests) // close channel

	// ranging over bursty requests, which is five ints
	for req := range burstyRequests {
		<-burstyLimiter // first three executes immediately, because limiter already has values buffered, the next two waits at 200ms intervals
		fmt.Println("request", req, time.Now())
	}

	// atomic counters
	var ops atomic.Uint64 // atomic integer type to represent our counter(always positive
	for i := 0; i < 50; i++ {
		wg.Add(1) // add task to waitgroup
		// perform task, this is starting 50 goroutines that each increment the counter 1000 times
		go func() {
			for c := 0; c < 1000; c++ {
				ops.Add(1) // atomically increment the counter
			}
			wg.Done() // finish the task
		}()
	}
	wg.Wait()                       // wait for task to finish
	fmt.Println("ops:", ops.Load()) // read the value of the counter

	// create a container object
	con := Container{
		// the zero value of a mutex is usable as-is, so no initialization is required here
		counters: map[string]int{"a": 0, "b": 0},
	}
	//This function increments a named counter in a loop.
	doIncrement := func(name string, n int) {
		for i := 0; i < n; i++ {
			con.inc(name)
		}
		wg.Done() // alert the waitgroup the task is done
	}
	wg.Add(3) // add task number to wait group
	//Run several goroutines concurrently; note that they all access the same Container, and two of them access the same counter.
	go doIncrement("a", 10000)
	go doIncrement("a", 10000)
	go doIncrement("b", 10000)
	wg.Wait() //Wait for the goroutines to finish
	fmt.Println("container counters: ", con.counters)

	var readOps uint64
	var writeOps uint64
	// create unbuffered channel
	reads := make(chan readOp)
	writes := make(chan writeOp)

	// here is the goroutine that owns the state, which is a map as in the previous example but now private to the stateful goroutine
	// this goroutine repeatedly selects on the reads and writes channels,responding to requests as they arrive.
	// a response is executed by first performing the reuqested opertaion and then sending the value on a response channel resp to indicate success
	go func() {
		var state = make(map[int]int)
		for {
			select {
			case read := <-reads:
				read.resp <- state[read.key]
			case write := <-writes:
				state[write.key] = write.val
				write.resp <- true
			}
		}
	}()

	// this starts 100 goroutines to issue reads to the state owning goroutine via the reads channel.
	//each read requires constructing a readOp,sending it over the reads channel, and then receiving the result over the provided resp channel
	for r := 0; r < 100; r++ {
		go func() {
			for {
				read := readOp{
					key:  rand.Intn(5),
					resp: make(chan int),
				}
				reads <- read
				<-read.resp
				atomic.AddUint64(&readOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	// we start 10 writes as well ,with a similar approach
	for w := 0; w < 10; w++ {
		go func() {
			for {
				write := writeOp{
					key:  rand.Intn(5),
					val:  rand.Intn(100),
					resp: make(chan bool)}
				writes <- write
				<-write.resp
				atomic.AddUint64(&writeOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}
	// let the goroutines work for a second
	time.Sleep(time.Second)

	// finally, capture and return the op counts
	readOpsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint64(&writeOps)
	fmt.Println("writeOps:", writeOpsFinal)
}
