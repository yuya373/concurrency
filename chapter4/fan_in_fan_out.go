package chapter4

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func toInt(done <-chan interface{}, input <-chan interface{}) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)

		for v := range input {
			select {
			case <-done:
				return
			case ch <- v.(int):
			}
		}
	}()
	return ch
}

func primeFinder(
	done <-chan interface{},
	input <-chan int,
) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)

		for i := range input {
			i -= 1
			prime := true
			for divisor := i - 1; divisor > 1; divisor-- {
				if i%divisor == 0 {
					prime = false
					break
				}
			}

			if prime {
				select {
				case <-done:
					return
				case ch <- i:
				}
			}
		}
	}()
	return ch
}

func NaivePrimeFinder() {
	rand := func() interface{} { return rand.Intn(50000000) }

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntCh := toInt(done, repeatFn(done, rand))
	source := primeFinder(done, randIntCh)

	fmt.Println("Primes:")

	for prime := range take(done, source, 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v", time.Since(start))
}

func fanIn(
	done <-chan interface{},
	channels ...<-chan interface{},
) <-chan interface{} {
	var wg sync.WaitGroup
	out := make(chan interface{})

	for _, c := range channels {
		wg.Add(1)
		go func(ch <-chan interface{}) {
			defer wg.Done()

			for v := range ch {
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func FanOutFanInPrimeFinder() {
	rand := func() interface{} { return rand.Intn(50000000) }

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntCh := toInt(done, repeatFn(done, rand))

	nCpu := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders.\n", nCpu)

	fmt.Println("Primes:")

	finders := make([]<-chan interface{}, nCpu)
	for i := 0; i < nCpu; i++ {
		finders[i] = primeFinder(done, randIntCh)
	}

	source := fanIn(done, finders...)

	for prime := range take(done, source, 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v", time.Since(start))
}
