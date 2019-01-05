package chapter4

import (
	"fmt"
)

func bridge(
	done <-chan interface{},
	chanCh <-chan (<-chan interface{}),
) <-chan interface{} {
	outCh := make(chan interface{})

	go func() {
		defer close(outCh)

		for {
			var ch <-chan interface{}
			select {
			case maybeCh, ok := <-chanCh:
				if ok == false {
					return
				}
				ch = maybeCh
			case <-done:
				return
			}

			for v := range orDone(done, ch) {
				select {
				case outCh <- v:
				case <-done:
				}
			}
		}
	}()

	return outCh
}

func Bridge() {
	generateValues := func() <-chan (<-chan interface{}) {
		ch := make(chan (<-chan interface{}))

		go func() {
			defer close(ch)
			for i := 0; i < 10; i++ {
				innerCh := make(chan interface{}, 1)
				innerCh <- i
				close(innerCh)
				ch <- innerCh
			}
		}()

		return ch
	}

	for v := range bridge(nil, generateValues()) {
		fmt.Printf("%v ", v)
	}
}
