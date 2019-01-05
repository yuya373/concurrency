package chapter4

import (
	"fmt"
	"math/rand"
	"time"
)

func SimpleLeak() {
	doWork := func(ch <-chan string) <-chan interface{} {
		completed := make(chan interface{})

		go func() {
			defer fmt.Println("doWork exites.")
			defer close(completed)

			for s := range ch {
				fmt.Println(s)
			}
		}()

		return completed
	}

	doWork(nil) // leak
	fmt.Println("Done")
}

func FixSimpleLeak() {
	doWork := func(
		killWork <-chan interface{},
		inputCh <-chan string,
	) <-chan interface{} {
		workFinished := make(chan interface{})

		go func() {
			defer fmt.Println("doWork exited.")
			defer close(workFinished)

			for {
				select {
				case s := <-inputCh:
					fmt.Println(s)
				case <-killWork:
					return
				}
			}
		}()

		return workFinished
	}

	killWork := make(chan interface{})
	workFinished := doWork(killWork, nil)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(killWork)
	}()

	<-workFinished
	fmt.Println("Done.")
}

func SendLeak() {
	randCh := func() <-chan int {
		ch := make(chan int)
		go func() {
			defer fmt.Println("randCh closure exited.")
			defer close(ch)
			for {
				ch <- rand.Int()
			}
		}()
		return ch
	}

	ch := randCh()
	fmt.Println("3 random ints:")
	for i := 0; i < 3; i++ {
		fmt.Printf("%d: %d\n", i, <-ch)
	}
}

func FixSendLeak() {
	randCh := func(kill <-chan interface{}) <-chan int {
		ch := make(chan int)
		go func() {
			defer fmt.Println("randCh closure exited.")
			defer close(ch)

			for {
				select {
				case ch <- rand.Int():
				case <-kill:
					return
				}
			}
		}()
		return ch
	}

	killCh := make(chan interface{})
	ch := randCh(killCh)
	fmt.Println("3 random ints:")
	for i := 0; i < 3; i++ {
		fmt.Printf("%d: %d\n", i, <-ch)
	}
	close(killCh)
	time.Sleep(1 * time.Second)
}
