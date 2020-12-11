package manager

import (
	"fmt"
	"sync"

	"github.com/arstevens/go-hive-signal/internal/transmuter"
)

var ChangeTriggerLimit int = 20

/*SwarmManager is an object that can be used to connect new requesters
to a peer-to-peer swarm*/
type SwarmManager struct {
	gateway     SwarmGateway
	gatewayLock *sync.Mutex
	negotiate   AgentNegotiator
	tracker     SwarmSizeTracker
	closed      bool
	id          string
	changes     int
}

//New creates a new SwarmManager
func New(swarmID string, gateway SwarmGateway, negotiate AgentNegotiator, tracker SwarmSizeTracker) *SwarmManager {
	return &SwarmManager{
		gateway:     gateway,
		gatewayLock: &sync.Mutex{},
		negotiate:   negotiate,
		tracker:     tracker,
		closed:      false,
		id:          swarmID,
		changes:     0,
	}
}

//AttemptToPair attempts to pair 'conn' with someone from the swarm
func (sm *SwarmManager) AttemptToPair(conn interface{}) error {
	if sm.closed {
		return fmt.Errorf("Failed to attempt pair in SwarmManager.AttemptToPair(). Object closed")
	}

	acceptorConn, ok := conn.(Conn)
	if !ok {
		return fmt.Errorf("Failed to pair in SwarmManager.AttemptToPair(). 'conn' does not conform to Conn interface")
	}
	sm.gatewayLock.Lock()
	offerer, err := sm.gateway.GetEndpoint()
	sm.gatewayLock.Unlock()
	if err != nil {
		return fmt.Errorf("Failed to pair in SwarmManager.AttemptToPair(): %v", err)
	}
	offererConn, ok := offerer.(Conn)
	if !ok {
		return fmt.Errorf("Failed to pair in SwarmManager.AttemptToPair(). offerer 'conn' does not conform to Conn interface")
	}

	err = sm.negotiate(offererConn, acceptorConn)
	if err != nil {
		return fmt.Errorf("Failed to pair in SwarmManager.AttemptToPair(): %v", err)
	}
	return nil
}

func (sm *SwarmManager) AddEndpoint(c interface{}) error {
	conn, ok := c.(Conn)
	if !ok {
		return fmt.Errorf("Failed to add endpoint in SwarmManager.AddEndpoint(): parameter of wrong type")
	}
	sm.gatewayLock.Lock()
	defer sm.gatewayLock.Unlock()
	err := sm.gateway.AddEndpoint(conn)
	if err != nil {
		return fmt.Errorf("Failed to add endpoint in SwarmManager.AddEndpoint(): %v", err)
	}
	sm.changes++
	if sm.changes > ChangeTriggerLimit {
		sm.tracker.SetSize(sm.id, sm.gateway.GetTotalEndpoints())
		sm.changes = 0
	}
	return nil
}

func (sm *SwarmManager) RemoveEndpoint(c interface{}) error {
	conn, ok := c.(Conn)
	if !ok {
		return fmt.Errorf("Failed to add endpoint in SwarmManager.RemoveEndpoint(): parameter of wrong type")
	}
	sm.gatewayLock.Lock()
	defer sm.gatewayLock.Unlock()
	err := sm.gateway.RetireEndpoint(conn)
	if err != nil {
		return fmt.Errorf("Failed to add endpoint in SwarmManager.RemoveEndpoint(): %v", err)
	}
	sm.changes++
	if sm.changes > ChangeTriggerLimit {
		sm.tracker.SetSize(sm.id, sm.gateway.GetTotalEndpoints())
		sm.changes = 0
	}
	return nil
}

func (sm *SwarmManager) Bisect() (transmuter.SwarmManager, error) {
	sm.gatewayLock.Lock()
	defer sm.gatewayLock.Unlock()
	newGateway, err := sm.gateway.EvenlySplit()
	if err != nil {
		return nil, fmt.Errorf("Failed to bisect in SwarmManager.Bisect(): %v", err)
	}
	sm.tracker.SetSize(sm.id, sm.gateway.GetTotalEndpoints())
	newManager := New("", newGateway, sm.negotiate, sm.tracker)
	return newManager, nil
}

func (sm *SwarmManager) Stitch(m transmuter.SwarmManager) error {
	manager, ok := m.(*SwarmManager)
	if !ok {
		return fmt.Errorf("Failed to stitch in SwarmManager.Stitch(): Wrong parameter type")
	}
	sm.gatewayLock.Lock()
	defer sm.gatewayLock.Unlock()
	err := sm.gateway.Merge(manager.gateway)
	if err != nil {
		return fmt.Errorf("Failed to stitch in SwarmManager.Stitch(): %v", err)
	}
	sm.tracker.SetSize(sm.id, sm.gateway.GetTotalEndpoints())
	sm.tracker.SetSize(manager.id, 0)
	manager.Close()
	return nil
}

//GetID returns the swarm ID associated with this manager
func (sm *SwarmManager) GetID() string {
	return sm.id
}

func (sm *SwarmManager) SetID(id string) {
	sm.gatewayLock.Lock()
	sm.tracker.SetSize(id, sm.gateway.GetTotalEndpoints())
	sm.gatewayLock.Unlock()
	sm.id = id
}

//Close closes the SwarmManager for use
func (sm *SwarmManager) Close() error {
	if sm.closed {
		return fmt.Errorf("Failed to close in SwarmManager.Close() Already closed")
	}
	sm.gateway.Close()
	sm.closed = true
	return nil
}
