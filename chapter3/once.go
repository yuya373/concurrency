package chapter3

import (
	"fmt"
	"sync"
)

func Once() {
	var count int

	inc := func() { count++ }

	var once sync.Once
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			once.Do(inc)
		}()
	}

	wg.Wait()
	fmt.Println("Count is", count)
}

func Pitfall() {
	var count int

	inc := func() { count++ }
	dec := func() { count-- }

	var once sync.Once
	once.Do(inc)
	once.Do(dec)

	fmt.Println("Count is", count)
}

func Deadlock() {
	var a, b sync.Once
	var initB func()

	initA := func() { b.Do(initB) }
	initB = func() { a.Do(initA) }
	a.Do(initA)
}
