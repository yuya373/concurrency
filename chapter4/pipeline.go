package chapter4

import (
	"fmt"
)

func multiply(values []int, multiplier int) []int {
	results := make([]int, len(values))

	for i, v := range values {
		results[i] = v * multiplier
	}

	return results
}

func add(values []int, additive int) []int {
	results := make([]int, len(values))

	for i, v := range values {
		results[i] = v + additive
	}

	return results
}

func BasicPipeline() {
	ints := []int{1, 2, 3, 4}
	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}

	for _, v := range multiply(add(multiply(ints, 2), 1), 2) {
		fmt.Println(v)
	}
}

func ChannelPipeline() {
	generator := func(
		done <-chan interface{},
		integers ...int,
	) <-chan int {
		ch := make(chan int, len(integers))
		go func() {
			defer close(ch)
			for _, i := range integers {
				select {
				case <-done:
					return
				case ch <- i:
				}
			}
		}()
		return ch
	}

	multiply := func(
		done <-chan interface{},
		intCh <-chan int,
		multiplier int,
	) <-chan int {
		ch := make(chan int)
		go func() {
			defer close(ch)
			for i := range intCh {
				select {
				case <-done:
					return
				case ch <- i * multiplier:
				}
			}
		}()
		return ch
	}

	add := func(
		done <-chan interface{},
		intCh <-chan int,
		additive int,
	) <-chan int {
		ch := make(chan int)
		go func() {
			defer close(ch)
			for i := range intCh {
				select {
				case <-done:
					return
				case ch <- i + additive:
				}
			}
		}()
		return ch
	}

	done := make(chan interface{})
	defer close(done)

	intCh := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intCh, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}
