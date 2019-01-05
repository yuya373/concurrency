package chapter4

import (
	"fmt"
	"time"
)

func orDone(done, c <-chan interface{}) <-chan interface{} {
	ch := make(chan interface{})

	go func() {
		defer close(ch)

		for {
			// fmt.Println("in For Loop")
			select {
			case <-done:
				// fmt.Println("Receive done1")
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case ch <- v:
				case <-done:
					// fmt.Println("Receive done2")
				}
			}
		}
	}()

	// in For Loop
	// closing rep chan
	// in For Loop
	// Receive done2
	// in For Loop
	// Receive done1
	// 9

	return ch
}

func OrDone() {
	done1 := make(chan interface{})
	done2 := make(chan interface{})

	rep := repeat(done2, 9)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("closing rep chan")
		close(done1)
		// close(done2)
	}()

	for i := range orDone(done1, rep) {
		fmt.Println(i)
	}
}
