package gateway

import (
	"fmt"

	"github.com/arstevens/go-hive-signal/internal/manager"
)

var DefaultQueueCapacity = 100

type activeConnectionQueue struct {
	queue []manager.Conn
	head  int
	tail  int
	size  int
}

func newActiveConnectionQueue(capacity int) *activeConnectionQueue {
	if capacity == 0 {
		capacity = DefaultQueueCapacity
	}

	return &activeConnectionQueue{
		queue: make([]manager.Conn, capacity),
		head:  0,
		tail:  0,
		size:  0,
	}
}

func (aq *activeConnectionQueue) IsEmpty() bool    { return aq.size == 0 }
func (aq *activeConnectionQueue) GetCapacity() int { return len(aq.queue) }
func (aq *activeConnectionQueue) GetSize() int     { return aq.size }
func (aq *activeConnectionQueue) GetAddrs() []string {
	addrs := make([]string, aq.size)
	for i := 0; i < aq.size; i++ {
		c := aq.queue[(aq.head+i)%len(aq.queue)]
		addrs[i] = c.GetAddress()
	}
	return addrs
}

func (aq *activeConnectionQueue) Resize(newSize int) error {
	if newSize < aq.size {
		return fmt.Errorf("Too many entries. Cannot resize queue in ActiveConnectQueue.Resize()")
	} else if newSize == aq.size {
		return nil
	}

	newQueue := make([]manager.Conn, newSize)
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

func (aq *activeConnectionQueue) Push(c manager.Conn) error {
	if aq.IsFull() {
		return fmt.Errorf("Queue is full in activeConnectionQueue.Push()")
	}
	aq.queue[aq.tail] = c
	aq.tail = (aq.tail + 1) % len(aq.queue)
	aq.size++
	return nil
}

func (aq *activeConnectionQueue) Pop() manager.Conn {
	if aq.size == 0 {
		return nil
	}

	c := aq.queue[aq.head]
	aq.queue[aq.head] = nil
	aq.head = (aq.head + 1) % len(aq.queue)
	aq.size--
	return c
}
