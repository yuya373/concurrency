package chapter4

import (
	"fmt"
)

func tee(
	done <-chan interface{},
	in <-chan interface{},
) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	go func() {
		defer close(out1)
		defer close(out2)

		for val := range orDone(done, in) {
			var out1, out2 = out1, out2

			for i := 0; i < 2; i++ {
				select {
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()

	return out1, out2
}

func Tee() {
	done := make(chan interface{})
	defer close(done)

	source := repeat(done, 1, 2)
	out1, out2 := tee(done, take(done, source, 4))

	for v := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", v, <-out2)
	}
}
