package chapter5

import (
	"fmt"
	"time"
)

func doWork(
	done <-chan interface{},
	pulseInterfal time.Duration,
) (<-chan interface{}, <-chan time.Time) {
	heartbeat := make(chan interface{})
	results := make(chan time.Time)

	go func() {
		defer close(heartbeat)
		defer close(results)

		pulse := time.NewTicker(pulseInterfal)
		defer pulse.Stop()
		workGen := time.NewTicker(2 * pulseInterfal)
		defer workGen.Stop()

		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}:
			default:
			}
		}

		sendResult := func(r time.Time) {
			for {
				select {
				case <-done:
					return
				case <-pulse.C:
					sendPulse()
				case results <- r:
					return
				}
			}
		}

		for {
			select {
			case <-done:
				return
			case <-pulse.C:
				sendPulse()
			case r := <-workGen.C:
				sendResult(r)
			}
		}
	}()

	return heartbeat, results
}

func HeartBeat() {
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)

	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Printf("results: %v\n", r.Second())
		case <-time.After(timeout):
			return
		}
	}
}

func doWorkTwice(
	done <-chan interface{},
	pulseInterval time.Duration,
) (<-chan interface{}, <-chan time.Time) {
	heartbeat := make(chan interface{})
	results := make(chan time.Time)

	go func() {
		pulse := time.NewTicker(pulseInterval)
		workGen := time.NewTicker(2 * pulseInterval)

		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}:
			default:
			}
		}

		sendResult := func(r time.Time) {
			for {
				select {
				case <-pulse.C:
					sendPulse()
				case results <- r:
					return
				}
			}
		}

		for i := 0; i < 2; i++ {
			select {
			case <-done:
				return
			case <-pulse.C:
				sendPulse()
			case r := <-workGen.C:
				sendResult(r)
			}
		}
	}()

	return heartbeat, results
}

func HeartBeatWithFail() {
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := doWorkTwice(done, timeout/2)

	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Printf("results: %v\n", r)
		case <-time.After(timeout):
			fmt.Println("worker goroutine is not healthy!")
			return
		}
	}
}
