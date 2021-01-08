package gateway

import (
	"fmt"
	"log"
	"sync"

	"github.com/arstevens/go-hive-signal/internal/manager"
)

var DialEndpoint func(addr string) (manager.Conn, error) = nil

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

func (sg *SwarmGateway) GetEndpoint() (manager.Conn, int, error) {
	sg.aqMutex.Lock()
	if sg.activeQueue.IsEmpty() {
		sg.aqMutex.Unlock()
		err := sg.populateActiveQueue()
		if err != nil {
			return nil, -1, fmt.Errorf("Failed in SwarmGateway.GetEndpoint(): %v", err)
		}
	} else {
		sg.aqMutex.Unlock()
	}

	sg.aqMutex.Lock()
	sg.pcMutex.Lock()
	conn, debriefValue := sg.activeQueue.Pop()
	for conn != nil && conn.IsClosed() {
		delete(sg.priorityCache, conn.GetAddress())
		conn, debriefValue = sg.activeQueue.Pop()
	}
	if conn == nil {
		sg.aqMutex.Unlock()
		sg.pcMutex.Unlock()
		return nil, debriefValue, fmt.Errorf("No active endpoints in SwarmGateway.GetEndpoint()")
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
	return conn, debriefValue, nil
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
	addrs := make([]string, iqSize)

	i := 0
	for _, entry := range *sg.inactiveQueue.pq {
		addrs[i] = entry.address
		i++
	}
	activeAddrs := sg.activeQueue.GetAddrs()
	addrs = append(addrs, activeAddrs...)

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
		conn, _ := sg.activeQueue.Pop()
		conn.Close()
	}
	sg.aqMutex.Unlock()
	return nil
}

func (sg *SwarmGateway) populateActiveQueue() error {
	sg.iqMutex.Lock()
	defer sg.iqMutex.Unlock()
	if sg.inactiveQueue.IsEmpty() {
		log.Println("SwarmGateway.populateActiveQueue(): InactiveQueue is empty")
		return nil
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
