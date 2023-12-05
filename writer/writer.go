package writer

import (
	"fmt"
	"io"
	"sync"
)

type Node struct {
	data []byte
	next *Node
}

func (n *Node) Size() int {
	return len(n.data)
}

func NewNode(p []byte) *Node {
	return &Node{
		data: p,
		next: nil,
	}
}

type Buffer struct {
	head     *Node
	tail     *Node
	capacity int
	size     int
}

func (b *Buffer) Write(p []byte) error {
	lp := len(p)
	if b.capacity < lp {
		return fmt.Errorf("buffer capacity=%d id lower than message size=%d", b.capacity, lp)
	}

	b.shrink(lp)
	head := NewNode(p)
	if b.head == nil {
		b.head = head
		b.tail = head
		return nil
	}

	b.head.next = head
	b.head = head
	return nil
}

func (b *Buffer) Read() []byte {
	if b.tail == nil {
		return nil
	}
	res := b.tail.data
	b.tail = b.tail.next
	if b.tail == nil {
		b.head = nil
	}
	return res
}

func (b *Buffer) shrink(s int) {
	for b.size+s > b.capacity {
		b.size -= b.tail.Size()
		b.tail = b.tail.next
	}
}

type Writer struct {
	inner  io.Writer
	mutex  sync.Mutex
	buffer *Buffer
	chanel chan int
}

func NewWriter(w io.Writer, cap int) *Writer {
	res := &Writer{
		inner:  w,
		mutex:  sync.Mutex{},
		buffer: &Buffer{capacity: cap},
		chanel: make(chan int, cap),
	}

	go res.loop()
	return res
}

func (w *Writer) Write(b []byte) (int, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	err := w.buffer.Write(b)
	if err != nil {
		return 0, err
	}
	w.chanel <- len(b)

	return len(b), nil
}

func (w *Writer) loop() {
	for range w.chanel {
		w.buffer.Read()
		_, err := w.inner.Write(w.buffer.Read())
		if err != nil {
		}
	}
}
