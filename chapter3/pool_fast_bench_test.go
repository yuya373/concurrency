package chapter3

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

func fConnectToService() interface{} {
	time.Sleep(1 * time.Second)
	return struct{}{}
}

func warmServiceConnCache() *sync.Pool {
	p := &sync.Pool{
		New: fConnectToService,
	}

	for i := 0; i < 10; i++ {
		p.Put(p.New())
	}
	return p
}

func fStartNetworkDaemon() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		pool := warmServiceConnCache()
		server, err := net.Listen("tcp", "localhost:9999")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer server.Close()
		wg.Done()

		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue
			}
			cache := pool.Get()
			fmt.Fprintln(conn, "")
			pool.Put(cache)
			conn.Close()
		}
	}()
	return &wg
}

func init() {
	daemonStarted := fStartNetworkDaemon()
	daemonStarted.Wait()
}

func BenchmarkFastNetworkRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:9999")
		if err != nil {
			b.Fatalf("cannot dial host: %v", err)
		}
		if _, err := ioutil.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %v", err)
		}
		conn.Close()
	}
}

// [I] ❯ go test -benchtime=10s -bench=. ./chapter3/pool_fast_bench_test.go
// goos: darwin
// goarch: amd64
// BenchmarkFastNetworkRequest-4               2000           5194733 ns/op
// PASS
// ok      command-line-arguments  29.870s
