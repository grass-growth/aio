package aio

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func bechmarkWithSizes(b *testing.B, fn func(*testing.B, []byte)) {
	sizes := []int{128, 256, 512, 1024, 2048, 4096, 8192}
	for _, size := range sizes {
		buf := make([]byte, size)
		rand.Read(buf)

		b.Run(fmt.Sprint(size), func(sb *testing.B) {
			fn(sb, buf)
		})
	}
}

func createOut(b *testing.B) io.Writer {
	b.Helper()

	filename := filepath.Join(b.TempDir(), strings.ReplaceAll(b.Name(), "/", "_"))
	f, err := os.Create(filename)
	if err != nil {
		b.Fatalf("Failed to create file: %v", err)
	}

	b.Cleanup(func() {
		f.Close()
		os.Remove(filename)
	})

	return f
}

func BenchmarkIOPipe(b *testing.B) {
	bechmarkWithSizes(b, func(sb *testing.B, buf []byte) {
		out := createOut(sb)

		pr, pw := io.Pipe()
		donech := make(chan struct{})

		go func() {
			defer close(donech)
			io.Copy(out, pr)
		}()

		sb.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := pw.Write(buf)
			if err != nil {
				sb.Fatal(err)
			}
		}

		pw.Close()
		<-donech
		pr.Close()
	})
}

func BenchmarkNoPipe(b *testing.B) {
	bechmarkWithSizes(b, func(sb *testing.B, buf []byte) {
		out := createOut(sb)

		sb.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := out.Write(buf)
			if err != nil {
				sb.Fatal(err)
			}
		}
	})
}
