package chapter4

import (
	"fmt"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-or(append(channels[3:], orDone)...):
			}
		}
	}()

	return orDone
}

func Or() {
	signal := func(after time.Duration) <-chan interface{} {
		ch := make(chan interface{})

		go func() {
			defer close(ch)
			time.Sleep(after)
		}()

		return ch
	}

	start := time.Now()
	<-or(
		signal(2*time.Hour),
		signal(5*time.Minute),
		signal(1*time.Second),
		signal(1*time.Hour),
		signal(1*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}
