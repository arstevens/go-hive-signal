package gateway

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestActiveConnectionQueue(t *testing.T) {
	fmt.Printf("---------------------\nACTIVE CONNECTION QUEUE TEST\n---------------------------\n")
	queueSize := 10
	queue := NewActiveConnectionQueue(queueSize)
	for i := 0; i < queueSize; i++ {
		queue.Push(&FakeConn{
			addr:   "/address/" + strconv.Itoa(i),
			closed: false,
		})
	}
	fmt.Printf("(Add Test)\n\tPopping off with starting size %d...\n", queue.size)
	for i := 0; i < queueSize; i++ {
		c := queue.Pop()
		if c == nil {
			t.Fatalf("Invalid NIL return on queue pop")
		}
		fmt.Printf("\t%s\n", c.GetAddress())
	}
	fmt.Printf("\tEnding size of %d\n", queue.size)

	for i := 0; i < queueSize+1; i++ {
		err := queue.Push(&FakeConn{
			addr:   "/address/" + strconv.Itoa(i*10),
			closed: false,
		})
		if i == queueSize && err == nil {
			t.Fatal(fmt.Errorf("Capacity overflow should have occurred"))
		}
	}

	fmt.Printf("(Wrap and Overflow Test)\n\tPopping off with starting size %d...\n", queue.size)
	for i := 0; i < queueSize+1; i++ {
		c := queue.Pop()
		if i == queueSize && c != nil {
			t.Fatalf("Non NIL return on empty queue pop")
		}
		if c == nil && i != queueSize {
			t.Fatalf("Invalid NIL return on queue pop")
		}
		if c != nil {
			fmt.Printf("\t%s\n", c.GetAddress())
		}
	}
	fmt.Printf("\tEnding size of %d\n", queue.size)

	for i := 0; i < queueSize/2; i++ {
		queue.Push(&FakeConn{
			addr:   "/address/" + strconv.Itoa(i*12),
			closed: false,
		})
	}

	newCap := queueSize * 2
	fmt.Printf("(Resize Test)\n\tResizing from capacity %d->%d with size %d...\n", queueSize, newCap, queue.size)
	err := queue.Resize(newCap)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("\tResized successfully\n")

	err = queue.Push(&FakeConn{
		addr:   "/address/NEWENRY",
		closed: false})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("\tPopping off with starting size %d...\n", queue.size)
	for queue.size > 0 {
		c := queue.Pop()
		if c == nil {
			t.Fatalf("Invalid NIL return on queue pop")
		}
		fmt.Printf("\t%s\n", c.GetAddress())
	}
	fmt.Printf("\tEnding size of %d\n", queue.size)

}

func TestEndpointPriorityQueue(t *testing.T) {
	fmt.Printf("---------------------\nINACTIVE PRIORITY QUEUE TEST\n---------------------------\n")
	pq := NewEndpointPriorityQueue()

	rand.Seed(time.Now().UnixNano())
	totalItems := 50
	for i := 0; i < totalItems; i++ {
		prio := rand.Intn(100)
		addr := strconv.Itoa(prio)
		pq.Push(addr, prio)
	}

	fmt.Printf("Random Insert Pop Order: \n")
	printPQ(pq)

	for i := 0; i < totalItems; i++ {
		addr := strconv.Itoa(i)
		pq.PushNew(addr)
	}
	fmt.Printf("\nPushNew Insertions Pop Order: \n")
	printPQ(pq)

	removal := strconv.Itoa(totalItems / 2)
	pq.PushNew(removal)
	fmt.Printf("\nRemoving Entry %s: \n", removal)
	pq.Remove(removal)
	printPQ(pq)

}

func printPQ(pq *EndpointPriorityQueue) {
	n := pq.Size()
	if n == 0 {
		fmt.Printf("\tEMPTY\n")
		return
	}
	for i := 0; i < n; i++ {
		if i%10 == 0 {
			fmt.Printf("\n\t")
		}
		fmt.Printf("%s ", pq.Pop())
	}
}

type FakeConn struct {
	addr   string
	closed bool
}

func (fc *FakeConn) Read([]byte) (int, error)  { return 0, nil }
func (fc *FakeConn) Write([]byte) (int, error) { return 0, nil }
func (fc *FakeConn) Close() error              { fc.closed = true; return nil }
func (fc *FakeConn) IsClosed() bool            { return fc.closed }
func (fc *FakeConn) GetAddress() string        { return fc.addr }
