package gateway

import "fmt"

var DefaultQueueCapacity = 100

type activeConnectionQueue struct {
	queue []Conn
	head  int
	tail  int
	size  int
}

func newActiveConnectionQueue(capacity int) *activeConnectionQueue {
	if capacity == 0 {
		capacity = DefaultQueueCapacity
	}

	return &activeConnectionQueue{
		queue: make([]Conn, capacity),
		head:  0,
		tail:  0,
		size:  0,
	}
}

func (aq *activeConnectionQueue) IsEmpty() bool    { return aq.size == 0 }
func (aq *activeConnectionQueue) GetCapacity() int { return len(aq.queue) }
func (aq *activeConnectionQueue) GetSize() int     { return aq.size }

func (aq *activeConnectionQueue) Resize(newSize int) error {
	if newSize < aq.size {
		return fmt.Errorf("Too many entries. Cannot resize queue in ActiveConnectQueue.Resize()")
	} else if newSize == aq.size {
		return nil
	}

	newQueue := make([]Conn, newSize)
	idx := 0
	for i := aq.head; i != aq.tail; i = (i + 1) % len(aq.queue) {
		newQueue[idx] = aq.queue[i]
		aq.queue[i] = nil
		idx++
	}
	aq.queue = newQueue
	aq.head = 0
	aq.tail = aq.size
	return nil
}

func (aq *activeConnectionQueue) IsFull() bool {
	return aq.size == len(aq.queue)
}

func (aq *activeConnectionQueue) Push(c Conn) error {
	if aq.IsFull() {
		return fmt.Errorf("Queue is full in activeConnectionQueue.Push()")
	}
	aq.queue[aq.tail] = c
	aq.tail = (aq.tail + 1) % len(aq.queue)
	aq.size++
	return nil
}

func (aq *activeConnectionQueue) Pop() Conn {
	if aq.size == 0 {
		return nil
	}

	c := aq.queue[aq.head]
	aq.queue[aq.head] = nil
	aq.head = (aq.head + 1) % len(aq.queue)
	aq.size--
	return c
}
