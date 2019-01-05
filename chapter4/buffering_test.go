package chapter4

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func performWrite(b *testing.B, writer io.Writer) {
	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for bt := range take(done, repeat(done, byte(0)), b.N) {
		writer.Write([]byte{bt.(byte)})
	}
}

func tmpFileOrFatal() *os.File {
	file, err := ioutil.TempFile("", "tmp")

	if err != nil {
		log.Fatal(err)
	}

	return file
}

func BenchmarkUnbufferedWrite(b *testing.B) {
	performWrite(b, tmpFileOrFatal())
}

func BenchmarkBufferedWrite(b *testing.B) {
	bufferedFile := bufio.NewWriter(tmpFileOrFatal())
	performWrite(b, bufferedFile)
}

// [I] ❯ go test -benchtime=10s -bench=. ./chapter4/buffering_test.go ./chapter4/generator.go
// goos: darwin
// goarch: amd64
// BenchmarkUnbufferedWrite-4       1000000             10637 ns/op
// BenchmarkBufferedWrite-4        10000000              1445 ns/op
// PASS
// ok      command-line-arguments  26.562s
