package chapter3

import (
	"fmt"
	"sync"
)

func Pool() {
	p := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new instance")
			return struct{}{}
		},
	}

	p.Get()
	i := p.Get()
	p.Put(i)
	p.Get()
}

func MemPool() {
	var n int

	p := &sync.Pool{
		New: func() interface{} {
			n += 1
			mem := make([]byte, 1024)
			return &mem
		},
	}

	for i := 0; i < 2; i++ {
		p.Put(p.New())
	}

	const nWorkers = 1024 * 1024
	var wg sync.WaitGroup
	for i := nWorkers; i > 0; i-- {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mem := p.Get().(*[]byte)
			defer p.Put(mem)
		}()
	}

	wg.Wait()
	fmt.Println(n, "calculator was created.")
}
