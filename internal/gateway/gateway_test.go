package gateway

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestSwarmGateway(t *testing.T) {
	fmt.Printf("---------------------------\n    SWARM GATEWAY TEST\n---------------------------\n")
	dialEndpoint = func(addr string) (Conn, error) {
		return &FakeConn{addr: addr, closed: false}, nil
	}

	activeSize := 10
	inactiveSize := 20
	gateway := New(activeSize, inactiveSize)
	totalAdds := activeSize + inactiveSize

	fmt.Printf("Populating gateway with %d connections...\n", totalAdds)
	for i := 0; i < totalAdds; i++ {
		conn := &FakeConn{addr: "/address/" + strconv.Itoa(i), closed: false}
		err := gateway.AddEndpoint(conn)
		if err != nil {
			t.Fatal(err)
		}
	}
	fmt.Printf("\tTotal Connections: active=%d inactive=%d\n", gateway.activeQueue.GetSize(), gateway.inactiveQueue.GetSize())

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
	fmt.Printf("\tTotal Connections: active=%d inactive=%d\n", gateway.activeQueue.GetSize(), gateway.inactiveQueue.GetSize())

	fmt.Printf("Retiring %d endpoints...\n", len(removals))
	for i := 0; i < len(removals); i++ {
		err := gateway.RetireEndpoint(removals[i])
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("\tRemoved endpoint at %s\n", removals[i].GetAddress())
	}
	fmt.Printf("\tTotal Connections: active=%d inactive=%d\n", gateway.activeQueue.GetSize(), gateway.inactiveQueue.GetSize())

	fmt.Printf("Evenly splitting swarm...\n")
	g2, err := gateway.EvenlySplit()
	if err != nil {
		t.Fatal(err)
	}
	gateway2 := g2.(*SwarmGateway)
	fmt.Printf("\t(Original)Total Connections: active=%d inactive=%d\n", gateway.activeQueue.GetSize(), gateway.inactiveQueue.GetSize())
	fmt.Printf("\t(New)Total Connections: active=%d inactive=%d\n", gateway2.activeQueue.GetSize(), gateway2.inactiveQueue.GetSize())

	fmt.Printf("Merging Original swarm into New swarm...\n")
	err = gateway2.Merge(gateway)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("\t(Original)Total Connections: active=%d inactive=%d\n", gateway.activeQueue.GetSize(), gateway.inactiveQueue.GetSize())
	fmt.Printf("\t(New)Total Connections: active=%d inactive=%d\n", gateway2.activeQueue.GetSize(), gateway2.inactiveQueue.GetSize())

	fmt.Printf("Closing New swarm...\n")
	gateway2.Close()
	fmt.Printf("\t(New)Total Connections: active=%d inactive=%d\n", gateway2.activeQueue.GetSize(), gateway2.inactiveQueue.GetSize())
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
	fmt.Printf("\n---------------------------\nINACTIVE PRIORITY QUEUE TEST\n---------------------------\n")
	pq := newEndpointPriorityQueue()

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

func printPQ(pq *endpointPriorityQueue) {
	n := pq.GetSize()
	if n == 0 {
		fmt.Printf("\tEMPTY\n")
		return
	}
	for i := 0; i < n; i++ {
		if i%10 == 0 {
			fmt.Printf("\n\t")
		}
		s, _ := pq.Pop()
		fmt.Printf("%s ", s)
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
