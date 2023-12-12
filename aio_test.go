package aio

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/djherbis/buffer"
	"github.com/djherbis/nio/v3"
)

func benchmarkWithSizes(b *testing.B, fn func(*testing.B, []byte)) {
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
	benchmarkWithSizes(b, func(sb *testing.B, buf []byte) {
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

func BenchmarkIOPipeWithTimeout(b *testing.B) {
	benchmarkWithSizes(b, func(sb *testing.B, buf []byte) {
		out := createOut(sb)

		pr, pw := io.Pipe()
		donech := make(chan struct{})

		go func() {
			defer close(donech)
			var buf [4096]byte

			for {
				n, err := pr.Read(buf[:])
				if n > 0 {
					out.Write(buf[:n])
				}

				if err == io.EOF {
					break
				}

				if err != nil {
					b.Logf("Failed to read: %v", err)
					return
				}

				time.Sleep(1 * time.Microsecond)
			}
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
	benchmarkWithSizes(b, func(sb *testing.B, buf []byte) {
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

func BenchmarkNioPipe(b *testing.B) {
	benchmarkWithSizes(b, func(sb *testing.B, buf []byte) {
		out := createOut(sb)

		pbuf := buffer.New(32 * 1024) // 32KB In memory Buffer
		pr, pw := nio.Pipe(pbuf)
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

func BenchmarkNioPipeWithTimeout(b *testing.B) {
	benchmarkWithSizes(b, func(sb *testing.B, buf []byte) {
		out := createOut(sb)

		pbuf := buffer.New(32 * 1024) // 32KB In memory Buffer
		pr, pw := nio.Pipe(pbuf)
		donech := make(chan struct{})

		go func() {
			defer close(donech)
			var buf [4096]byte

			for {
				n, err := pr.Read(buf[:])
				if n > 0 {
					out.Write(buf[:n])
				}

				if err == io.EOF {
					break
				}

				if err != nil {
					b.Logf("Failed to read: %v", err)
					return
				}

				time.Sleep(1 * time.Microsecond)
			}
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
