package gateway

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func TestSwarmGateway(t *testing.T) {
	fmt.Printf("---------------------------\n    SWARM GATEWAY TEST\n---------------------------\n")

	activeSize := 11
	gateway := New(activeSize-1)

	fmt.Printf("Populating gateway with %d connections...\n", activeSize)
	for i := 0; i < activeSize; i++ {
		conn := &FakeConn{addr: "/address/" + strconv.Itoa(i), closed: false}
		err := gateway.PushEndpoint(conn)
		if err != nil {
			t.Fatal(err)
		}
	}
	fmt.Printf("\tTotal Connections: active=%d\n", gateway.activeQueue.GetSize())

	removals := make([]*FakeConn, activeSize)
	fmt.Printf("Getting %d endpoints...\n", activeSize)
	for i := 0; i < activeSize; i++ {
		conn, err := gateway.GetEndpoint()
		if err != nil {
			t.Fatal(err)
		}
		c := conn.(*FakeConn)
		removals[i] = c
		fmt.Printf("\tFetched endpoint at %s\n", c.GetAddress())
	}
	fmt.Printf("\tTotal Connections: active=%d\n", gateway.activeQueue.GetSize())
}

func TestActiveConnectionQueue(t *testing.T) {
	fmt.Printf("\n---------------------------\nACTIVE CONNECTION QUEUE TEST\n---------------------------\n")
	queueSize := 10
	queue := newActiveConnectionQueue(queueSize)
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
		fmt.Printf("\tPopped: %s\n", c.GetAddress())
	}
	fmt.Printf("\tEnding size of %d\n", queue.size)

	for i := 0; i < queueSize; i++ {
		err := queue.Push(&FakeConn{
			addr:   "/address/" + strconv.Itoa(i*10),
			closed: false,
		})
		if err != nil {
			t.Fatalf("Failed: %v", err)
		}
	}

	fmt.Printf("(Wrap and Overflow Test)\n\tPopping off with starting size %d...\n", queue.size)
	for i := 0; i < queueSize; i++ {
		c := queue.Pop()
		if c == nil {
			t.Fatalf("Invalid NIL return on queue pop")
		}
		fmt.Printf("\tPopped: %s\n", c.GetAddress())
	}
	fmt.Printf("\tEnding size of %d\n", queue.size)

	for i := 0; i < queueSize/2; i++ {
		queue.Push(&FakeConn{
			addr:   "/address/" + strconv.Itoa(i*12),
			closed: false,
		})
	}

	err := queue.Push(&FakeConn{
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
		fmt.Printf("\tPopped: %s\n", c.GetAddress())
	}
	fmt.Printf("\tEnding size of %d\n", queue.size)

}

type FakeConn struct {
	addr   string
	closed bool
}

func (fc *FakeConn) Read(b []byte) (int, error) {
	for i := 0; i < len(b); i++ {
		b[i] = byte(rand.Intn(256))
	}
	return len(b), nil
}
func (fc *FakeConn) Write([]byte) (int, error) { return 0, nil }
func (fc *FakeConn) Close() error              { fc.closed = true; return nil }
func (fc *FakeConn) IsClosed() bool            { return fc.closed }
func (fc *FakeConn) GetAddress() string        { return fc.addr }
