package chapter3

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)

func UniDirectionalChan() {
	var recv <-chan interface{}
	var send chan<- interface{}

	ch := make(chan interface{})
	recv = ch
	send = ch

	go func() {
		send <- "Hello"
	}()

	fmt.Println(<-recv)
	// <-send // compile error: receive from send-only type
	// recv <- "Hello" // send to receive-only type
}

func ClosedChan() {
	intCh := make(chan int)
	close(intCh)
	integer, ok := <-intCh
	fmt.Printf("(%v): %v", ok, integer)
}

func HandleChannelWithLoop() {
	intChan := make(chan int)

	go func() {
		defer close(intChan)
		for i := 0; i < 5; i++ {
			intChan <- i
		}
	}()

	for i := range intChan {
		fmt.Printf("%v", i)
	}
}

func NotifyChannelWithClose() {
	beg := make(chan interface{})

	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-beg
			fmt.Printf("%v has begun\n", i)
		}(i)
	}

	fmt.Println("Unblocking goroutines...")
	close(beg)
	wg.Wait()
}

func BufferedChan() {
	var buf bytes.Buffer
	defer buf.WriteTo(os.Stdout)

	intCh := make(chan int, 4)

	go func() {
		defer close(intCh)
		defer fmt.Fprintln(&buf, "Producer Done.")

		for i := 0; i < 5; i++ {
			fmt.Fprintf(&buf, "Sending: %d\n", i)
			intCh <- i
		}
	}()

	for i := range intCh {
		fmt.Fprintf(&buf, "Received %v.\n", i)
	}
}

func NilChannel() {
	var ch chan interface{}
	// <-ch // deadlock
	// ch <- struct{}{} // deadlock
	close(ch) // panic
}

func OwnerAndConsumer() {
	chanOwner := func() <-chan int {
		ch := make(chan int, 5)
		go func() {
			defer close(ch)
			for i := 0; i <= 5; i++ {
				ch <- i
			}
		}()
		return ch
	}

	ch := chanOwner()
	for e := range ch {
		fmt.Printf("Received: %v\n", e)
	}
	fmt.Println("Done receiving!")
}

func Select() {
	start := time.Now()
	ch := make(chan interface{})

	go func() {
		defer close(ch)
		time.Sleep(5 * time.Second)
	}()

	fmt.Println("Blocking on read...")
	select {
	case <-ch:
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}
}

func SelectMultiple() {
	c1 := make(chan interface{})
	c2 := make(chan interface{})
	close(c1)
	close(c2)

	var c1Count, c2Count int

	for i := 1000; i > 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}

	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}

func SelectTimeout() {
	var c <-chan int // nil channel
	select {
	case <-c: // read from nil channel blocking goroutine
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out.")
	}
}

func SelectDefault() {
	start := time.Now()
	var c1, c2 <-chan int
	select {
	case <-c1:
	case <-c2:
	default:
		fmt.Printf("In default clause after %v\n\n", time.Since(start))
	}
}

func SelectBreak() {
	done := make(chan interface{})
	go func() {
		defer close(done)
		time.Sleep(5 * time.Second)
	}()

	workCounter := 0
loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}

		workCounter++
		time.Sleep(1 * time.Second)
	}

	fmt.Printf(
		"Achived %v cycles of work before signalled to stop.\n",
		workCounter,
	)
}
