package chapter4

import (
	"bytes"
	"fmt"
	"sync"
)

func AdhocBind() {
	data := make([]int, 4)
	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}

func LexicalBind() {
	producer := func() <-chan int {
		ch := make(chan int)
		go func() {
			defer close(ch)
			for i := 0; i < 5; i++ {
				fmt.Printf("Sending: %d\n", i)
				ch <- i
				fmt.Printf("Sent: %d\n", i)
			}
		}()
		return ch
	}

	consumer := func(ch <-chan int) {
		for e := range ch {
			fmt.Printf("Received: %d\n", e)
		}
		fmt.Println("Done receiving!")
	}

	ch := producer()
	consumer(ch)

	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	data := []byte{'g', 'o', 'l', 'a', 'n', 'g'}
	wg.Add(1)
	go printData(&wg, data[:2])
	wg.Add(1)
	go printData(&wg, data[2:])

	wg.Wait()
}
