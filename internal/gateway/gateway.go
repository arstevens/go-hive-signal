package gateway

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	"github.com/arstevens/go-hive-signal/internal/manager"
)

var OperationSuccess byte = 1
var MessageEndian = binary.LittleEndian
var NewConnWraper func(net.Conn) manager.Conn = nil

var DialEndpoint = func(addr string) (manager.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewConnWraper(conn), nil
}

type SwarmGateway struct {
	activeQueue   *activeConnectionQueue
	aqMutex       *sync.Mutex
	inactiveQueue *endpointPriorityQueue
	iqMutex       *sync.Mutex
	priorityCache map[string]int
	pcMutex       *sync.Mutex
}

func New(activeSize int, inactiveSize int) *SwarmGateway {
	return &SwarmGateway{
		activeQueue:   newActiveConnectionQueue(activeSize),
		aqMutex:       &sync.Mutex{},
		inactiveQueue: newEndpointPriorityQueue(),
		iqMutex:       &sync.Mutex{},
		priorityCache: make(map[string]int),
		pcMutex:       &sync.Mutex{},
	}
}

func (sg *SwarmGateway) PushEndpoint(addr string) error {
	sg.iqMutex.Lock()
	sg.inactiveQueue.PushNew(addr)
	sg.iqMutex.Unlock()
	return nil
}

func (sg *SwarmGateway) RemoveEndpoint(addr string) error {
	sg.iqMutex.Lock()
	sg.inactiveQueue.Remove(addr)
	sg.iqMutex.Unlock()
	return nil
}

func (sg *SwarmGateway) GetEndpoint() (manager.Conn, error) {
	sg.aqMutex.Lock()
	if sg.activeQueue.IsEmpty() {
		sg.aqMutex.Unlock()
		err := sg.populateActiveQueue()
		if err != nil {
			return nil, fmt.Errorf("Failed in SwarmGateway.GetEndpoint(): %v", err)
		}
	} else {
		sg.aqMutex.Unlock()
	}

	sg.aqMutex.Lock()
	sg.pcMutex.Lock()
	conn := sg.activeQueue.Pop()
	for conn != nil && conn.IsClosed() {
		delete(sg.priorityCache, conn.GetAddress())
		conn = sg.activeQueue.Pop()
	}
	if conn == nil {
		sg.aqMutex.Unlock()
		sg.pcMutex.Unlock()
		return nil, fmt.Errorf("No active endpoints in SwarmGateway.GetEndpoint()")
	}

	addr := conn.GetAddress()
	priority := sg.priorityCache[addr] + 1
	delete(sg.priorityCache, addr)
	sg.pcMutex.Unlock()
	sg.aqMutex.Unlock()

	sg.iqMutex.Lock()
	sg.inactiveQueue.Push(addr, priority)
	sg.iqMutex.Unlock()

	sg.populateActiveQueue()
	return conn, nil
}

func (sg *SwarmGateway) GetTotalEndpoints() int {
	sg.iqMutex.Lock()
	sg.aqMutex.Lock()
	defer sg.iqMutex.Unlock()
	defer sg.aqMutex.Unlock()
	return sg.activeQueue.GetSize() + sg.inactiveQueue.GetSize()
}

func (sg *SwarmGateway) GetEndpointAddrs() []string {
	sg.iqMutex.Lock()
	sg.aqMutex.Lock()

	iqSize := sg.inactiveQueue.GetSize()
	aqSize := sg.activeQueue.GetSize()
	addrs := make([]string, iqSize+aqSize)

	i := 0
	for _, entry := range *sg.inactiveQueue.pq {
		addrs[i] = entry.address
		i++
	}
	for _, conn := range sg.activeQueue.queue {
		addrs[i] = conn.GetAddress()
		i++
	}

	sg.iqMutex.Unlock()
	sg.aqMutex.Unlock()
	return addrs
}

func (sg *SwarmGateway) Close() error {
	sg.iqMutex.Lock()
	for !sg.inactiveQueue.IsEmpty() {
		sg.inactiveQueue.Pop()
	}
	sg.iqMutex.Unlock()

	sg.aqMutex.Lock()
	for !sg.activeQueue.IsEmpty() {
		conn := sg.activeQueue.Pop()
		conn.Close()
	}
	sg.aqMutex.Unlock()
	return nil
}

func (sg *SwarmGateway) populateActiveQueue() error {
	sg.iqMutex.Lock()
	defer sg.iqMutex.Unlock()
	if sg.inactiveQueue.IsEmpty() {
		return fmt.Errorf("InactiveQueue is empty")
	}

	sg.aqMutex.Lock()
	for !sg.inactiveQueue.IsEmpty() && !sg.activeQueue.IsFull() {
		sg.aqMutex.Unlock()
		addr, prio := sg.inactiveQueue.Pop()
		conn, err := DialEndpoint(addr)
		if err != nil {
			return err
		}
		sg.aqMutex.Lock()
		sg.activeQueue.Push(conn)

		sg.pcMutex.Lock()
		sg.priorityCache[addr] = prio
		sg.pcMutex.Unlock()
	}
	sg.aqMutex.Unlock()
	return nil
}
