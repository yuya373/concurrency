package chapter3

import (
	"fmt"
	"sync"
	"time"
)

func Cond() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("Removed from queue")
		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		fmt.Println("Iteration", i, "=================")
		c.L.Lock()

		for len(queue) == 2 {
			fmt.Println("Wait start")
			c.Wait()
			fmt.Println("Wait end")
		}

		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second)
		c.L.Unlock()
	}
}

type Button struct {
	Clicked *sync.Cond
}

func Broadcast() {
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) {
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			wg.Done()

			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		wg.Wait()
	}

	var wg sync.WaitGroup
	wg.Add(1)
	subscribe(button.Clicked, func() {
		defer wg.Done()
		fmt.Println("Maximizing window.")
	})
	wg.Add(1)
	subscribe(button.Clicked, func() {
		defer wg.Done()
		fmt.Println("Display annoying dialog box!")
	})
	wg.Add(1)
	subscribe(button.Clicked, func() {
		defer wg.Done()
		fmt.Println("Mouse clicked.")
	})

	button.Clicked.Broadcast()
	wg.Wait()
}
