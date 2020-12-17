package gateway

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/arstevens/go-hive-signal/internal/manager"
)

var OperationSuccess byte = 1
var MessageEndian = binary.LittleEndian
var NewConnWraper func(net.Conn) Conn = nil

var dialEndpoint = func(addr string) (Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewConnWraper(conn), nil
}

type SwarmGateway struct {
	activeQueue   *activeConnectionQueue
	inactiveQueue *endpointPriorityQueue
	priorityCache map[string]int
}

func New(activeSize int, inactiveSize int) *SwarmGateway {
	return &SwarmGateway{
		activeQueue:   newActiveConnectionQueue(activeSize),
		inactiveQueue: newEndpointPriorityQueue(),
		priorityCache: make(map[string]int),
	}
}

func (sg *SwarmGateway) AddEndpoint(c manager.Conn) error {
	defer c.Close()
	conn := c.(Conn)

	addr := conn.GetAddress()
	err := binary.Write(conn, MessageEndian, OperationSuccess)
	if err != nil {
		return fmt.Errorf("Failed to communicate endpoint addition in SwarmGateway.AddEndpoint(): %v", err)
	}
	sg.inactiveQueue.PushNew(addr)
	return nil
}

func (sg *SwarmGateway) RetireEndpoint(c manager.Conn) error {
	defer c.Close()
	conn := c.(Conn)
	addr := conn.GetAddress()
	sg.inactiveQueue.Remove(addr)
	return nil
}

func (sg *SwarmGateway) GetEndpoint() (manager.Conn, error) {
	if sg.activeQueue.IsEmpty() {
		err := sg.populateActiveQueue()
		if err != nil {
			return nil, fmt.Errorf("Failed in SwarmGateway.GetEndpoint(): %v", err)
		}
	}

	conn := sg.activeQueue.Pop()
	for conn != nil && conn.IsClosed() {
		delete(sg.priorityCache, conn.GetAddress())
		conn = sg.activeQueue.Pop()
	}
	if conn == nil {
		return nil, fmt.Errorf("No active endpoints in SwarmGateway.GetEndpoint()")
	}

	addr := conn.GetAddress()
	priority := sg.priorityCache[addr] + 1
	delete(sg.priorityCache, addr)
	sg.inactiveQueue.Push(addr, priority)

	sg.populateActiveQueue()
	return conn, nil
}

func (sg *SwarmGateway) EvenlySplit() (manager.SwarmGateway, error) {
	newActiveQueue := newActiveConnectionQueue(sg.activeQueue.GetCapacity())
	newInactiveQueue := newEndpointPriorityQueue()
	newPriorityCache := make(map[string]int)

	queueSize := sg.activeQueue.GetSize() / 2
	for i := 0; i < queueSize; i++ {
		conn := sg.activeQueue.Pop()
		newActiveQueue.Push(conn)

		addr := conn.GetAddress()
		newPriorityCache[addr] = sg.priorityCache[addr]
		delete(sg.priorityCache, addr)
	}

	iQueueSize := sg.inactiveQueue.GetSize() / 2
	for i := 0; i < iQueueSize; i++ {
		addr, prio := sg.inactiveQueue.Pop()
		newInactiveQueue.Push(addr, prio)
	}

	return &SwarmGateway{
		activeQueue:   newActiveQueue,
		inactiveQueue: newInactiveQueue,
		priorityCache: newPriorityCache,
	}, nil
}

func (sg *SwarmGateway) GetTotalEndpoints() int {
	return sg.activeQueue.GetSize() + sg.inactiveQueue.GetSize()
}

func (sg *SwarmGateway) Merge(g manager.SwarmGateway) error {
	gateway := g.(*SwarmGateway)

	for !sg.activeQueue.IsFull() && !gateway.activeQueue.IsEmpty() {
		conn := gateway.activeQueue.Pop()
		addr := conn.GetAddress()

		sg.priorityCache[addr] = gateway.priorityCache[addr]
		sg.activeQueue.Push(conn)
	}

	for !gateway.inactiveQueue.IsEmpty() {
		addr, prio := gateway.inactiveQueue.Pop()
		sg.inactiveQueue.Push(addr, prio)
	}
	gateway.Close()
	return nil
}

func (sg *SwarmGateway) Close() error {
  for !sg.inactiveQueue.IsEmpty() {
    sg.inactiveQueue.Pop()
  }
	for !sg.activeQueue.IsEmpty() {
		conn := sg.activeQueue.Pop()
		conn.Close()
	}
	return nil
}

func (sg *SwarmGateway) populateActiveQueue() error {
	if sg.inactiveQueue.IsEmpty() {
		return fmt.Errorf("InactiveQueue is empty")
	}
	for !sg.inactiveQueue.IsEmpty() && !sg.activeQueue.IsFull() {
		addr, prio := sg.inactiveQueue.Pop()
		conn, err := dialEndpoint(addr)
		if err != nil {
			return err
		}
		sg.activeQueue.Push(conn)
		sg.priorityCache[addr] = prio
	}
	return nil
}
