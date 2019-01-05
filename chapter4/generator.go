package chapter4

import (
	"fmt"
	"math/rand"
)

func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	ch := make(chan interface{})

	go func() {
		defer close(ch)

		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case ch <- v:
				}
			}
		}
	}()

	return ch
}

func take(done <-chan interface{}, in <-chan interface{}, n int) <-chan interface{} {
	ch := make(chan interface{})

	go func() {
		defer close(ch)

		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case v := <-in:
				ch <- v
			}
		}
	}()

	return ch
}

func Repeat() {
	done := make(chan interface{})
	defer close(done)

	source := repeat(done, 1)
	for n := range take(done, source, 10) {
		fmt.Printf("%v", n)
	}
}

func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	ch := make(chan interface{})

	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			case ch <- fn():
			}
		}
	}()

	return ch
}

func RepeatFn() {
	done := make(chan interface{})
	defer close(done)

	rand := func() interface{} {
		return rand.Int()
	}

	source := repeatFn(done, rand)

	for n := range take(done, source, 10) {
		fmt.Println(n)
	}
}
