package gateway

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/arstevens/go-hive-signal/internal/manager"
)

func TestSwarmGateway(t *testing.T) {
	fmt.Printf("---------------------------\n    SWARM GATEWAY TEST\n---------------------------\n")
	DialEndpoint = func(addr string) (manager.Conn, error) {
		return &FakeConn{addr: addr, closed: false}, nil
	}
	DebriefProcedure = func(conn io.Reader) interface{} {
		var debriefValue int32
		err := binary.Read(conn, binary.BigEndian, &debriefValue)
		if err != nil {
			log.Printf("Failed to debrief in gateway.debriefConnection(): %v", err)
			return -1
		}
		return int(debriefValue)
	}

	activeSize := 10
	inactiveSize := 20
	gateway := New(activeSize, inactiveSize)
	totalAdds := activeSize + inactiveSize

	fmt.Printf("Populating gateway with %d connections...\n", totalAdds)
	for i := 0; i < totalAdds; i++ {
		conn := &FakeConn{addr: "/address/" + strconv.Itoa(i), closed: false}
		err := gateway.PushEndpoint(conn.addr)
		if err != nil {
			t.Fatal(err)
		}
	}
	fmt.Printf("\tTotal Connections: active=%d inactive=%d\n", gateway.activeQueue.GetSize(), gateway.inactiveQueue.GetSize())

	removals := make([]*FakeConn, activeSize)
	fmt.Printf("Getting %d endpoints...\n", activeSize)
	for i := 0; i < activeSize; i++ {
		conn, pref, err := gateway.GetEndpoint()
		if err != nil {
			t.Fatal(err)
		}
		c := conn.(*FakeConn)
		removals[i] = c
		fmt.Printf("\tFetched endpoint at %s with Pref_Load %d\n", c.GetAddress(), pref)
	}
	fmt.Printf("\tTotal Connections: active=%d inactive=%d\n", gateway.activeQueue.GetSize(), gateway.inactiveQueue.GetSize())

	fmt.Printf("Retiring %d endpoints...\n", len(removals))
	for i := 0; i < len(removals); i++ {
		err := gateway.RemoveEndpoint(removals[i].addr)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("\tRemoved endpoint at %s\n", removals[i].GetAddress())
	}
	fmt.Printf("\tTotal Connections: active=%d inactive=%d\n", gateway.activeQueue.GetSize(), gateway.inactiveQueue.GetSize())

	fmt.Printf("Outputing all endpoint addrs...\n")
	addrs := gateway.GetEndpointAddrs()
	for i := 0; i < len(addrs); i++ {
		fmt.Printf("\t%s\n", addrs[i])
	}
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
		c, pref := queue.Pop()
		if c == nil {
			t.Fatalf("Invalid NIL return on queue pop")
		}
		fmt.Printf("\t%s Pref: %d\n", c.GetAddress(), pref)
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
		c, pref := queue.Pop()
		if i == queueSize && c != nil {
			t.Fatalf("Non NIL return on empty queue pop")
		}
		if c == nil && i != queueSize {
			t.Fatalf("Invalid NIL return on queue pop")
		}
		if c != nil {
			fmt.Printf("\t%s Pref: %d\n", c.GetAddress(), pref)
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
		c, pref := queue.Pop()
		if c == nil {
			t.Fatalf("Invalid NIL return on queue pop")
		}
		fmt.Printf("\t%s Pref: %d\n", c.GetAddress(), pref)
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
