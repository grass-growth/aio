package aio

import (
	"io"
	"os"
	"sync"
	"testing"
)

var testData = []byte("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz")

func BenchmarkWrite(b *testing.B) {
	// Create a pipe
	pr, pw := io.Pipe()

	// Create a buffer with text data

	var wg sync.WaitGroup

	// Read from the read end concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := io.Copy(os.Stdout, pr)
		if err != nil {
			panic(err)
		}
	}()

	// Run the Write function b.N times
	for i := 0; i < b.N; i++ {
		_, err := pw.Write(testData)
		if err != nil {
			b.Fatal(err)
		}
	}

	// Close the write end of the pipe
	pw.Close()

	wg.Wait()
	// Close the read end of the pipe
	pr.Close()
}
