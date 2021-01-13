package gateway

import (
	"fmt"
	"sync"

	"github.com/arstevens/go-hive-signal/internal/manager"
)

type SwarmGateway struct {
	activeQueue *activeConnectionQueue
	aqMutex     *sync.Mutex
}

func New(activeSize int) *SwarmGateway {
	return &SwarmGateway{
		activeQueue: newActiveConnectionQueue(activeSize),
		aqMutex:     &sync.Mutex{},
	}
}

func (sg *SwarmGateway) PushEndpoint(c manager.Conn) error {
	sg.aqMutex.Lock()
	sg.activeQueue.Push(c.(Conn))
	sg.aqMutex.Unlock()
	return nil
}

func (sg *SwarmGateway) GetEndpoint() (manager.Conn, error) {
	sg.aqMutex.Lock()
	conn := sg.activeQueue.Pop()
	for conn != nil && conn.IsClosed() {
		conn = sg.activeQueue.Pop()
	}
	if conn == nil {
		sg.aqMutex.Unlock()
		return nil, fmt.Errorf("No active endpoints in SwarmGateway.GetEndpoint()")
	}
	sg.activeQueue.Push(conn)
	sg.aqMutex.Unlock()
	return conn, nil
}

func (sg *SwarmGateway) GetTotalEndpoints() int {
	sg.aqMutex.Lock()
	defer sg.aqMutex.Unlock()
	return sg.activeQueue.GetSize()
}

func (sg *SwarmGateway) Close() error {
	sg.aqMutex.Lock()
	for !sg.activeQueue.IsEmpty() {
		conn := sg.activeQueue.Pop()
		conn.Close()
	}
	sg.aqMutex.Unlock()
	return nil
}
