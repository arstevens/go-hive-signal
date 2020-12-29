package manager

import (
	"fmt"
	"sync"
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
		gateway:   gateway,
		negotiate: negotiate,
		tracker:   tracker,
		closed:    false,
		id:        swarmID,
		changes:   0,
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
	offerer, err := sm.gateway.GetEndpoint()
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
	//Connect new endpoint with old endpoint so that state can be copied over
	conn, ok := c.(Conn)
	if !ok {
		return fmt.Errorf("Failed to add endpoint in SwarmManager.AddEndpoint(): parameter of wrong type")
	}
	err := connectForContextRetrieval(conn, sm.negotiate, sm.gateway)
	if err != nil {
		return err
	}

	//Add new endpoint to the gateway structure
	err = sm.gateway.AddEndpoint(conn)
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

func (sm *SwarmManager) TakeEndpoint(addr string) error {
	err := sm.gateway.PushEndpointAddr(addr)
	if err != nil {
		err = fmt.Errorf("Failed to take endpoint in SwarmManager.TakeEndpoint(): %v", err)
	}
	return err
}

func connectForContextRetrieval(conn Conn, negotiate AgentNegotiator, gateway SwarmGateway) error {
	offerer, err := gateway.GetEndpoint()
	if err != nil {
		/*If there was an error getting an endpoint thats because the swarm is
		empty and therefore there is no context to retrieve*/
		return nil
	}
	offererConn, ok := offerer.(Conn)
	if !ok {
		return fmt.Errorf("Failed to negotiate in SwarmManager.AddEndpoint(): Connection of wrong type")
	}
	err = negotiate(offererConn, conn)
	if err != nil {
		return fmt.Errorf("Failed to negotiate in SwarmManager.AddEndpoint(): %v", err)
	}
	return nil
}

func (sm *SwarmManager) RemoveEndpoint(c interface{}) error {
	conn, ok := c.(Conn)
	if !ok {
		return fmt.Errorf("Failed to add endpoint in SwarmManager.RemoveEndpoint(): parameter of wrong type")
	}
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

func (sm *SwarmManager) DropEndpoint(addr string) error {
	err := sm.gateway.DropEndpointAddr(addr)
	if err != nil {
		err = fmt.Errorf("Failed to drop endpoint in SwarmManager.DropEndpoint(): %v", err)
	}
	return err
}

func (sm *SwarmManager) GetEndpoints() []string {
	return sm.gateway.GetEndpointAddrs()
}

//GetID returns the swarm ID associated with this manager
func (sm *SwarmManager) GetID() string {
	return sm.id
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
